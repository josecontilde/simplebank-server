package sqlc

import (
	"context"
	"fmt"
	"testing"

	"github.com/josecontilde/simplebank/utils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestTransferTxDeadlock(t *testing.T) {

	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	n := 20

	amount := utils.RandomDecimal(0.0, 40.0)
	fmt.Println(">> amount:", amount.String())

	errs := make(chan error)

	//run n concurrent transfers
	for i := 0; i < n; i++ {
		go func(i int) {

			fromAccount, toAccount := account1.ID, account2.ID

			fmt.Println(">> transfer n:", i)
			fmt.Println(">> from:", fromAccount, "to:", toAccount)
			if i%2 == 1 {
				fmt.Println(">> transfer n:", i)
				fmt.Println(">> from:", fromAccount, "to:", toAccount)
				fromAccount, toAccount = account2.ID, account1.ID
			}

			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccount,
				ToAccountID:   toAccount,
				Amount:        amount,
			})

			errs <- err
		}(i)
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check the final balances

	totalAmount := amount.Mul(decimal.NewFromInt(int64(n))).Round(2)

	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.True(t, updatedAccount1.Balance.GreaterThanOrEqual(decimal.Zero))
	require.True(t, updatedAccount2.Balance.GreaterThanOrEqual(decimal.Zero))

	require.True(t, account1.Balance.Equal(updatedAccount1.Balance))
	require.True(t, account2.Balance.Equal(updatedAccount2.Balance))

	fmt.Println(">> transfer total:", totalAmount.String())
}
func TestTransferTx(t *testing.T) {

	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)
	n := 10
	amount := utils.RandomDecimal(0.0, 40.0)
	fmt.Println(">> amount:", amount.String())

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background()
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}
	// check results
	existed := make(map[int64]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		if err != nil {
			fmt.Println(">> err:", err)
		}
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.True(t, amount.Equal(transfer.Amount))

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountsID)
		require.Equal(t, amount.Neg(), fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountsID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)

		//check accounts balance
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance.Sub(fromAccount.Balance)
		diff2 := toAccount.Balance.Sub(account2.Balance)

		require.Equal(t, diff1, diff2)
		require.True(t, diff1.GreaterThan(decimal.Zero))
		require.True(t, diff1.Mod(amount).Equal(decimal.Zero))

		k := diff1.Div(amount).IntPart()
		require.True(t, k >= 1 && k <= int64(n))
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	totalAmount := amount.Mul(decimal.NewFromInt(int64(n))).Round(2)

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.True(t,
		account1.Balance.Sub(totalAmount).Equal(updatedAccount1.Balance),
		"Expected: %s, got: %s",
		account1.Balance.Sub(totalAmount).String(),
		updatedAccount1.Balance.String(),
	) // amount * n ...

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.True(t,
		account2.Balance.Add(totalAmount).Equal(updatedAccount2.Balance),
		"Expected: %s, got: %s",
		account2.Balance.Add(totalAmount).String(),
		updatedAccount2.Balance.String(),
	) // amount * n ...
	require.True(t, updatedAccount1.Balance.GreaterThanOrEqual(decimal.Zero))
	require.True(t, updatedAccount2.Balance.GreaterThanOrEqual(decimal.Zero))
	fmt.Println(">> transfer total:", totalAmount.String())
	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
}
