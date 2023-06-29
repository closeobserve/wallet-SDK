package starcoin

import (
	"testing"

	"github.com/coming-chat/wallet-SDK/core/testcase"
	"github.com/stretchr/testify/require"
)

func TestAccount(t *testing.T) {
	mnemonic := testcase.M1
	account, err := NewAccountWithMnemonic(mnemonic)
	require.Nil(t, err)

	// private key hex to account
	priHex, err := account.PrivateKeyHex()
	require.Nil(t, err)
	acc2, err := AccountWithPrivateKey(priHex)
	require.Nil(t, err)
	require.Equal(t, account.address, acc2.address)

	// public key hex to address
	pubHex := account.PublicKeyHex()
	addr2, err := EncodePublicKeyToAddress(pubHex)
	require.Nil(t, err)
	require.Equal(t, addr2, account.address)
}

func M1Account(t *testing.T) *Account {
	account, err := NewAccountWithMnemonic(testcase.M1)
	require.Nil(t, err)
	return account
}

func M2Account(t *testing.T) *Account {
	account, err := NewAccountWithMnemonic(testcase.M2)
	require.Nil(t, err)
	return account
}

func TestAccountWithPrivatekey(t *testing.T) {
	mnemonic := testcase.M1
	accountFromMnemonic, err := NewAccountWithMnemonic(mnemonic)
	require.Nil(t, err)
	privateKey, err := accountFromMnemonic.PrivateKeyHex()
	require.Nil(t, err)

	accountFromPrikey, err := AccountWithPrivateKey(privateKey)
	require.Nil(t, err)

	require.Equal(t, accountFromMnemonic.Address(), accountFromPrikey.Address())
}
