package common

import (
	"crypto/sha256"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"math/big"
	"net"
	"testing"
	"time"
)

func TestAddBlock(t *testing.T, store TransactionStorage) {
	peerIp := net.TCPAddr{
		IP:   net.ParseIP("128.0.0.1"),
		Port: 1000,
		Zone: "",
	}
	err := store.AddBlock("testing", time.Now(), "testing node", "testing node name", &peerIp, 0,
		sha256.Sum256([]byte("something")))
	require.NoError(t, err)
}

func TestAddTransaction(t *testing.T, store TransactionStorage) {
	peerIp := net.TCPAddr{
		IP:   net.ParseIP("128.0.0.1"),
		Port: 1000,
		Zone: "",
	}

	for i := 0; i < 10; i++ {
		err := store.AddTransaction("testing",
			time.Now(),
			"testing node",
			"testing node name",
			&peerIp,
			byte(2), // txType
			common.Address{},
			&common.Address{},
			100,
			234234,        // gasLimit
			big.NewInt(0), // gasPrice
			big.NewInt(0), // gasTipCap
			big.NewInt(0), // gasFeeCap
			big.NewInt(0), // txAmount
			[]byte("payload"),
			big.NewInt(0),
			big.NewInt(100),
			big.NewInt(0),
			common.Hash{},
		)
		require.NoError(t, err)
	}
}
