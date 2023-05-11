package hybrid

import (
	"errors"
	"math/big"
	"net"
	"path"
	"time"
	txstorecommon "txstore/common"
	"txstore/structuredlog"

	"github.com/ethereum/go-ethereum/common"
	"github.com/syndtr/goleveldb/leveldb"
)

// HybridStorage logs transaction metadata in a structured log, and store transaction content in a key value store
type HybridStorage struct {
	db *leveldb.DB
	structuredlog.LoggerBasedStorage
}

func NewHybridStorage(logDir string) (txstorecommon.TransactionStorage, error) {
	db, err := leveldb.OpenFile(path.Join(logDir, "txdata.db"), nil)
	if err != nil {
		return nil, err
	}

	log, err := structuredlog.NewLoggerBasedStorage(logDir)
	if err != nil {
		return nil, err
	}

	logBasedStore, ok := log.(structuredlog.LoggerBasedStorage)
	if !ok {
		return nil, errors.New("logBasedStore is not log based???")
	}

	return &HybridStorage{
		db,
		logBasedStore,
	}, nil
}

func (h *HybridStorage) LogTransactionMessage(txHash common.Hash,
	myId string,
	messageType byte,
	timeSeen time.Time) error {
	return nil
}

func (h *HybridStorage) LogBlockMessage(blockHash common.Hash,
	myId string,
	messageType byte,
	timeSeen time.Time) error {
	return nil
}

func (h *HybridStorage) AddTransaction(
	myId string,
	timeSeen time.Time,
	peerInfo string,
	peerName string,
	peerAddr net.Addr,
	txType byte,
	txFrom common.Address,
	txTo *common.Address,
	txNonce uint64,
	gasLimit uint64,
	gasPrice *big.Int,
	gasTipCap *big.Int,
	gasFeeCap *big.Int,
	txAmount *big.Int,
	txPayload []byte,
	sigV *big.Int,
	sigR *big.Int,
	sigS *big.Int,
	txHash common.Hash) error {
	err := h.LoggerBasedStorage.AddTransaction(
		myId,
		timeSeen,
		peerInfo,
		peerName,
		peerAddr,
		txType,
		txFrom,
		txTo,
		txNonce,
		gasLimit,
		gasPrice,
		gasTipCap,
		gasFeeCap,
		txAmount,
		nil,
		sigS, sigR, sigV,
		txHash)
	if err != nil {
		return err
	}

	has, err := h.db.Has(txHash.Bytes(), nil)
	if err != nil {
		return err
	}

	if !has {
		// add to DB if not already
		return h.db.Put(txHash.Bytes(), txPayload, nil)
	}
	return nil
}

func (h *HybridStorage) Close() error {
	return h.db.Close()
}

func (h *HybridStorage) GetStoreName() string {
	return "hybrid"
}
