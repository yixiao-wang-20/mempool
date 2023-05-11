package structuredlog

import (
	"math/big"
	"net"
	"os"
	"path"
	"time"
	txstorecommon "txstore/common"

	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	txLogFileName    = "tx.log"
	blockLogFileName = "blocks.log"
)

type LoggerBasedStorage struct {
	txLog    *zap.Logger
	blockLog *zap.Logger
}

func (back LoggerBasedStorage) LogTransactionMessage(txHash common.Hash,
	myId string,
	messageType byte,
	timeSeen time.Time) error {
	return nil
}

func (back LoggerBasedStorage) LogBlockMessage(blockHash common.Hash,
	myId string,
	messageType byte,
	timeSeen time.Time) error {
	return nil
}

func (back LoggerBasedStorage) AddTransaction(
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
	ip, port, err := txstorecommon.NetToPortInt(peerAddr)
	if err != nil {
		return err
	}

	txToStr := ""
	if txTo != nil {
		txToStr = txTo.String()
	}

	back.txLog.Info(txHash.String(),
		zap.String("monitor_id", myId),
		zap.Time("timeSeen", timeSeen),
		zap.String("peerInfo", peerInfo),
		zap.String("peerName", peerName),
		zap.String("peerIP", ip),
		zap.Int("peerPort", port),
		zap.Int8("txType", int8(txType)),
		zap.String("txFrom", txFrom.String()),
		zap.String("txTo", txToStr),
		zap.Uint64("nonce", txNonce),
		zap.Uint64("gasLimit", gasLimit),
		zap.String("gasPrice", gasPrice.String()),
		zap.String("gasTipCap", gasTipCap.String()),
		zap.String("gasFeeCap", gasFeeCap.String()),
		zap.String("value", txAmount.String()),
		zap.String("payload", common.Bytes2Hex(txPayload)),
		zap.String("sigV", sigV.String()),
		zap.String("sigR", sigR.String()),
		zap.String("sigS", sigS.String()))

	return nil
}

func (back LoggerBasedStorage) AddBlock(myId string, timeSeen time.Time, peerInfo string, peerName string, peerAddr net.Addr, blockNum uint64, hash common.Hash) error {
	ip, port, err := txstorecommon.NetToPortInt(peerAddr)
	if err != nil {
		return err
	}

	back.blockLog.Info(hash.String(),
		zap.String("monitor_id", myId),
		zap.Time("timeSeen", timeSeen),
		zap.String("peerInfo", peerInfo),
		zap.String("peerName", peerName),
		zap.String("peerIP", ip),
		zap.Int("peerPort", port),
		zap.Uint64("blockNum", blockNum))

	return nil
}

func (_ LoggerBasedStorage) GetStoreName() string {
	return "zaplog"
}

func (back LoggerBasedStorage) Close() error {
	back.txLog.Sync()
	return back.blockLog.Sync()
}

func NewLoggerBasedStorage(logDir string) (txstorecommon.TransactionStorage, error) {
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.Mkdir(logDir, os.ModePerm); err != nil {
			return nil, err
		}
	}

	txLogConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "tx_hash",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{path.Join(logDir, txLogFileName)},
		ErrorOutputPaths: []string{"stderr"},
	}
	txLogConfig.OutputPaths = []string{path.Join(logDir, txLogFileName)}

	blockLogConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "block_hash",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{path.Join(logDir, blockLogFileName)},
		ErrorOutputPaths: []string{"stderr"},
	}

	txLogger, err := txLogConfig.Build()
	if err != nil {
		return nil, err
	}

	blockLogger, err := blockLogConfig.Build()
	if err != nil {
		return nil, err
	}

	return LoggerBasedStorage{
		txLog:    txLogger,
		blockLog: blockLogger,
	}, nil
}
