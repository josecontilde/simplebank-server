package sqlc

import (
	"context"
	"testing"

	"github.com/josecontilde/simplebank/utils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  decimal.New(utils.RandomMoney(), 2),
		Currency: utils.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.True(t, arg.Balance.Equal(account.Balance))
	require.Equal(t, arg.Currency, account.Currency)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestListAccounts(t *testing.T) {
	createRandomAccount(t)

	accounts, err := testQueries.ListAccounts(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

func TestGetAccount(t *testing.T) {
	accountCreated := createRandomAccount(t)

	accountCurrent, err := testQueries.GetAccount(context.Background(), accountCreated.ID)
	require.NoError(t, err)
	require.NotEmpty(t, accountCurrent)

	require.Equal(t, accountCreated.ID, accountCurrent.ID)
	require.Equal(t, accountCreated.Owner, accountCurrent.Owner)
	require.True(t, accountCreated.Balance.Equal(accountCurrent.Balance))
	require.Equal(t, accountCreated.Currency, accountCurrent.Currency)
}

func TestUpdateAccount(t *testing.T) {

	account := createRandomAccount(t)
	require.NotEmpty(t, account)

	updateArg := UpdateAccountParams{
		ID:       account.ID,
		Owner:    utils.RandomOwner(),
		Balance:  decimal.New(utils.RandomMoney(), 2),
		Currency: utils.RandomCurrency(),
	}

	updateAccount, err := testQueries.UpdateAccount(context.Background(), updateArg)
	require.NoError(t, err)
	require.NotEmpty(t, updateAccount)

	require.Equal(t, account.ID, updateAccount.ID)
	require.NotEqual(t, account.Owner, updateAccount.Owner)
	require.NotEqual(t, account.Balance, updateAccount.Balance)
	require.NotEqual(t, account.Currency, updateAccount.Currency)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)
	require.NotEmpty(t, account)

	accountDeleted, err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, accountDeleted)

	require.Equal(t, account.ID, accountDeleted.ID)
	require.Equal(t, account.Owner, accountDeleted.Owner)
	require.True(t, account.Balance.Equal(accountDeleted.Balance))
	require.Equal(t, account.Currency, accountDeleted.Currency)

	account3, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.Empty(t, account3)
}
