package txstore

import (
	"fmt"
	"errors"
	"txstore/db"

	_ "github.com/lib/pq"
)

type MempoolMonitor struct {
	myId      string
	txStorage db.DatabasedBackedStorage
}

func NewMempoolMonitor() (*MempoolMonitor, error) {
	myId := "measurement"

	if myId == "" {
		fmt.Println("ERROR CREATING MEMPOOL MONITOR")
		return &MempoolMonitor{}, errors.New("cannot get MONITOR_ID or MEMPOOL_LOG_DIR or MEMPOOL_LOG_BACKEND")
	}

	var store_point *db.DatabasedBackedStorage
	var err error

	store_point, err = db.NewLocalhostTransactionDatabase(true)
//    store :=*store_point

	if err != nil {
		fmt.Println("ERROR CREATING MEMPOOL MONITOR: %v", err)
		return &MempoolMonitor{}, err
	}
	fmt.Println("MEMPOOL MONITOR HAS BEEN CREATED")
	return &MempoolMonitor{
		myId:      myId,
		txStorage: *store_point,
	}, nil
}

func (m MempoolMonitor) LogTransactionMessage(txmeslist []db.Txmes) error {
	return m.txStorage.LogTransactionMessage(txmeslist, m.myId,)
}

func (m MempoolMonitor) LogBlockMessage(blmeslist []db.Blmes) error {
	return m.txStorage.LogBlockMessage(blmeslist, m.myId,)
}

func (m MempoolMonitor) LogTransaction(txinfolist []db.Txinfo) error {
	return m.txStorage.AddTransaction(txinfolist,)
}

func (m MempoolMonitor) LogNewBlock(blinfolist []db.Blinfo) error {
	return m.txStorage.AddBlock(blinfolist,)
}

func (m MempoolMonitor) Close() error {
	return m.txStorage.Close()
}
