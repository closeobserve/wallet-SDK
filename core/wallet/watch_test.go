package wallet

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChainType(t *testing.T) {
	address := "😁"

	chains := ChainTypeOfWatchAddress(address)
	t.Log(chains.JsonString())
	t.Log(chains.Count())
}

func TestChainTypeOfPrivateKey(t *testing.T) {
	tests := []struct {
		name   string
		prikey string
		want   string
	}{
		{
			name:   "emoji",
			prikey: "😁",
			want:   "null",
		},
		{
			name:   "length 64",
			prikey: "0xfc0e2f9586b6ba8e4380737250824b64e7abc1d5e26d4357097809ad27e5e096",
			want:   `["bitcoin","ethereum","polka","signet","dogecoin","cosmos","terra","aptos","sui","starcoin"]`,
		},
		{
			name:   "length 128",
			prikey: "0xfc0e2f9586b6ba8e431d5e26d43537250824b64e7abc1a8e424b64e7abc97809ad27e5e096fc0e2f9586b6380737d5e26d4357080772508b097809ad27e5e096",
			want:   `["solana"]`,
		},
		{
			name:   "length 64, without 0x",
			prikey: "fc0e2f9586b6ba8e4380737250824b64e7abc1d5e26d4357097809ad27e5e096",
			want:   `["bitcoin","ethereum","polka","signet","dogecoin","cosmos","terra","aptos","sui","starcoin"]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ChainTypeOfPrivateKey(tt.prikey)
			t.Log(got.JsonString())
			require.Equal(t, got.JsonString(), tt.want)
		})
	}
}

func TestChainTypeOfPrivateKey_btc(t *testing.T) {
	prikey := "cTkZaPpb1pDdor36V5VY4uu5LE6tgzrjRADvrEXimEqWqvwRbfXY"
	typeArr := ChainTypeOfPrivateKey(prikey)
	t.Log(typeArr)
}
