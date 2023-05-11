package db

import (
	_ "github.com/lib/pq"
)

func TestDatabasedBackedStorage_AddBlock(t *testing.T) {
	db, err := NewLocalhostTransactionDatabase(User, Password, Server, DbName, true)
	require.NoError(t, err)
	defer db.Close()


	peerIp := net.TCPAddr{
		IP:   net.ParseIP("128.0.0.1"),
		Port: 1000,
		Zone: "",
	}
	err = db.AddBlock("testing", time.Now(), "testing node", "testing node name", &peerIp, 0,
		common.Hash{})
	require.NoError(t, err)
}

func TestDatabasedBackedStorage_AddTransaction(t *testing.T) {
	db, err := NewLocalhostTransactionDatabase(User, Password, Server, DbName, true)
	require.NoError(t, err)
	defer db.Close()

	peerIp := net.TCPAddr{
		IP:   net.ParseIP("128.0.0.1"),
		Port: 1000,
		Zone: "",
	}
	for i := 0; i < 100; i++ {
		err = db.AddTransaction("fmt.Sprintf("%s%d", "test tx ",i)" ,
		time.Now(),
		"testing node",
		"testing node name",
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
	// err = db.AddTransaction("testing",
	// 	time.Now(),
	// 	"testing node",
	// 	"testing node name",
	// 	&peerIp,
	// 	common.Address{},	
	// 	&common.Address{},
	// 	100,
	// 	big.NewInt(0),
	// 	234234,
	// 	big.NewInt(0),
	// 	[]byte("payload"),
	// 	big.NewInt(0),
	// 	big.NewInt(100),
	// 	big.NewInt(0),
	// 	common.Hash{},
	// 	)
	// require.NoError(t, err)
}
