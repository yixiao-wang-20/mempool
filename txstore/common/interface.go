package common

import (
	"io"
	"math/big"
	"net"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type TransactionStorage interface {
	io.Closer
	AddTransaction(myIp string,
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
		txHash common.Hash) error

	LogTransactionMessage(txHash common.Hash,
		myId string,
		messageType byte,
		timeSeen time.Time) error

	LogBlockMessage(blockHash common.Hash,
		myId string,
		messageType byte,
		timeSeen time.Time) error

	AddBlock(myIp string,
		timeSeen time.Time,
		peerInfo string,
		peerName string,
		peerAddr net.Addr,
		blockNum uint64,
		hash common.Hash) error

	GetStoreName() string
}

type WatchingList interface {
	// to may be nil
	ShouldLog(from common.Address, to *common.Address) bool
}
