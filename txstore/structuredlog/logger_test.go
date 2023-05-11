package structuredlog

import (
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	_ "github.com/lib/pq"
	"math/big"
	"testing"
	"txstore/db"
	"github.com/ethereum/go-ethereum/common"
	"fmt"
	"net"
	"time"
	"os"
	txstorecommon "txstore/common"
)

const testLogDir = "testlogs"

func TestLoggerBasedStorage_AddBlock(t *testing.T) {
	backend, err := NewLoggerBasedStorage(testLogDir)
	require.NoError(t, err)
	defer backend.Close()

	txstorecommon.TestAddBlock(t, backend)
}

func TestLoggerBasedStorage_AddTransaction(t *testing.T) {
	backend, err := NewLoggerBasedStorage(testLogDir)
	require.NoError(t, err)
	defer backend.Close()

	txstorecommon.TestAddTransaction(t, backend)
}

func TestDatabasedBackedStorage_AddTransaction(t *testing.T) {
	os.Setenv("LOG_LEVEL", "default")
	fmt.Println("Database Test")
	db, err := db.NewLocalhostTransactionDatabase(
		"postgres",
		"2312jsd2jklq99287nsxz",
		"95.216.0.93",
		"tx_log",
		true)	
	fmt.Println(db)

	require.NoError(t, err)
	defer db.Close()

	peerIp := net.TCPAddr{
		IP:   net.ParseIP("128.0.0.1"),
		Port: 1000,
		Zone: "",
	}
	for i := 0; i < 100; i++ {
		err = db.AddTransaction(fmt.Sprintf("%s%d", "test tx ",i) ,
		time.Now(),
		fmt.Sprintf("%s%d", "testing node ",i),
		fmt.Sprintf("%s%d", "testing node name  ",i),
		&peerIp,
		common.Address{},	
		&common.Address{},
		100,
		big.NewInt(0),
		234234,
		big.NewInt(0),
		[]byte("payload"),
		big.NewInt(0),
		big.NewInt(100),
		big.NewInt(0),
		common.Hash{},
		)
	require.NoError(t, err)
	}
	
}
