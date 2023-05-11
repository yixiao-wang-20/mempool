package db

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	// this credential only works after you connect to the server
	User     = "eth_mea"
	Password = "123456"
	Server   = "localhost"
	DbName   = "tx_log"
)

type Txmes struct {
	Hash common.Hash
	MessageType byte
	TimeSeen time.Time
}

type Txinfo struct {
	Tx *types.Transaction
	P *eth.Peer
	Signer types.Signer
}

type Blmes struct {
	Hash common.Hash
	MessageType byte
	TimeSeen time.Time
}

type Blinfo struct {
	Bl *types.Block
}

type DatabasedBackedStorage struct {
	db                     *sql.DB
	txMessageInsertStmt    *sql.Stmt
	txInsertStmt           *sql.Stmt
	blockMessageInsertStmt *sql.Stmt
	blockInsertStmt        *sql.Stmt
}

const (
	createTransactionMessageTable = `
	CREATE TABLE IF NOT EXISTS "tx_observation" (
		"hash" bpchar(66) NOT NULL,
		"type" int2 NOT NULL,
		"monitor_id" varchar(20) NOT NULL,
		"time_seen" timestamptz NOT NULL
	);`

	createTransactionTable = `
	CREATE TABLE IF NOT EXISTS "tx" (
		"hash" text NOT NULL,
		"from" text NOT NULL,
		"to" text,
		"gas" int8 NOT NULL,
		"gasPrice" text NOT NULL,
		"maxFeePerGas" text,
		"maxPriorityFeePerGas" text,
		"value" text NOT NULL,
		"nonce" int8 NOT NULL,
		"input" text,
		"type" int8 NOT NULL,
		"r" text NOT NULL,
		"s" text NOT NULL,
		"v" text NOT NULL
	);`

	createBlockMessageTable = `
	CREATE TABLE IF NOT EXISTS "block_observation" (
		"hash" bpchar(66) NOT NULL,
		"type" int2 NOT NULL,
		"monitor_id" varchar(20) NOT NULL,
		"time_seen" timestamptz NOT NULL
	);`

	createBlockTable = `
	CREATE TABLE IF NOT EXISTS "block" (
		"hash" text NOT NULL,
		"blockNumber" int8,
		"size" int8,
		"baseFeePerGas" text NOT NULL,
		"parentHash" text NOT NULL,
		"sha3Uncles" text NOT NULL,
		"miner" text NOT NULL,
		"stateRoot" text NOT NULL,
		"transactionsRoot" text NOT NULL,
		"receiptsRoot" text NOT NULL,
		"logsBloom" bytea,
		"difficulty" text NOT NULL,
		"gasLimit" int8 NOT NULL,
		"gasUsed" int8 NOT NULL,
		"timestamp" int8 NOT NULL,
		"extraData" bytea,
		"mixHash" text NOT NULL,
		"nonce" int8,
		"transactions" text[],
		"uncles" int8[]
	);`
)

var	insertTransactionMessage = `
	INSERT INTO "tx_observation" (
		"hash",
		"type",
		"monitor_id",
		"time_seen"
	) VALUES ($1, $2, $3, $4)`

var	insertTransaction = `
	INSERT INTO "tx" (
		"hash",
		"from",
		"to",
		"gas",
		"gasPrice",
		"maxFeePerGas",
		"maxPriorityFeePerGas",
		"value",
		"nonce",
		"input",
		"type",
		"r",
		"s",
		"v"
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

var	insertBlockMessage = `
	INSERT INTO "block_observation" (
		"hash",
		"type",
		"monitor_id",
		"time_seen"
	) VALUES ($1, $2, $3, $4)`

var	insertBlock = `
	INSERT INTO "block" (
		"hash",
		"blockNumber",
		"size",
		"baseFeePerGas",
		"parentHash",
		"sha3Uncles",
		"miner",
		"stateRoot",
		"transactionsRoot",
		"receiptsRoot",
		"logsBloom",
		"difficulty",
		"gasLimit",
		"gasUsed",
		"timestamp",
		"extraData",
		"mixHash",
		"nonce",
		"transactions",
		"uncles"
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`
 

func (db DatabasedBackedStorage) CreateTables() error {
	_, err := db.db.Exec(createTransactionMessageTable)
	if err != nil {
		return err
	}

	_, err = db.db.Exec(createTransactionTable)
	if err != nil {
		return err
	}

	_, err = db.db.Exec(createBlockMessageTable)
	if err != nil {
		return err
	}

	_, err = db.db.Exec(createBlockTable)
	if err != nil {
		return err
	}
  
  return nil
}

func (db DatabasedBackedStorage) LogTransactionMessage(txmeslist []Txmes, myId string) (err error) {

	valueArgs := make([]interface{}, 0)

	for i := 0; i < 200; i++ {
		valueArgs = append(valueArgs, txmeslist[i].Hash, myId, txmeslist[i].MessageType, txmeslist[i].TimeSeen)
	}

	_, err = db.txMessageInsertStmt.Exec(valueArgs...)

	if err != nil {
		return errors.Wrap(err, "cannot add block to storage")
	}

	return nil
}

func (db DatabasedBackedStorage) LogBlockMessage(blmeslist []Blmes, myId string) (err error) {

	valueArgs := make([]interface{}, 0)

	for i := 0; i < 200; i++ {
		valueArgs = append(valueArgs, blmeslist[i].Hash, myId, blmeslist[i].MessageType, blmeslist[i].TimeSeen)
	}
			
	_, err = db.blockMessageInsertStmt.Exec(valueArgs...)

	if err != nil {
		return errors.Wrap(err, "cannot add block to storage")
	}

	return nil
}

func (db DatabasedBackedStorage) AddTransaction(txinfolist []Txinfo) (err error) {

	valueArgs := make([]interface{}, 0)

	for i := 0; i < 100; i++ {

		txSender, err := types.Sender(txinfolist[i].Signer, txinfolist[i].Tx)
		if err != nil {
			return err
		}

		r, s, v := txinfolist[i].Tx.RawSignatureValues()
		
		var txToStr = ""
		if txinfolist[i].Tx.To() != nil {
			txToStr = txinfolist[i].Tx.To().String()
		}
			
		valueArgs = append(
			valueArgs,
			txinfolist[i].Tx.Hash().String(),
			txSender.String(),
			txToStr,
			txinfolist[i].Tx.Gas(),
			txinfolist[i].Tx.GasPrice().String(),
			txinfolist[i].Tx.GasFeeCap().String(),
			txinfolist[i].Tx.GasTipCap().String(),
			txinfolist[i].Tx.Value().String(),
			txinfolist[i].Tx.Nonce(),
			hex.EncodeToString(txinfolist[i].Tx.Data()),
			txinfolist[i].Tx.Type(),
			r.String(),
			s.String(),
			v.String(),
		)
	}
	
	_, err = db.txInsertStmt.Exec(valueArgs...)

	if err != nil {
		return errors.Wrap(err, "cannot add tx to storage")
	}

	return nil
}

func (db DatabasedBackedStorage) AddBlock(blinfolist []Blinfo) (err error) {

	valueArgs := make([]interface{}, 0)

	for i := 0; i < 30; i++ {
		transactionsHash := []string{}
		for _, tx := range blinfolist[i].Bl.Transactions() {
			transactionsHash=append(transactionsHash, tx.Hash().String())
		}
		unclesNum := []uint64{}
		for _, uncle := range blinfolist[i].Bl.Uncles() {
			unclesNum=append(unclesNum, uncle.Number.Uint64())
		}
	
		valueArgs = append(
			valueArgs, 
			blinfolist[i].Bl.Hash().String(),
			blinfolist[i].Bl.NumberU64(),
			blinfolist[i].Bl.Size(),
			blinfolist[i].Bl.BaseFee().String(),
			blinfolist[i].Bl.ParentHash().String(),
			blinfolist[i].Bl.UncleHash().String(),
			blinfolist[i].Bl.Coinbase().String(),
			blinfolist[i].Bl.Root().String(),
			blinfolist[i].Bl.TxHash().String(),
			blinfolist[i].Bl.ReceiptHash().String(),
			blinfolist[i].Bl.Bloom(),
			blinfolist[i].Bl.Difficulty().String(),
			blinfolist[i].Bl.GasLimit(),
			blinfolist[i].Bl.GasUsed(),
			blinfolist[i].Bl.Time(),
			blinfolist[i].Bl.Extra(),
			blinfolist[i].Bl.MixDigest().String(),
			blinfolist[i].Bl.Nonce(),
			transactionsHash,
			unclesNum,
		)
	}
			
	_, err = db.blockInsertStmt.Exec(valueArgs...)

	if err != nil {
		return errors.Wrap(err, "cannot add block to storage")
	}

	return nil
}

func (_ DatabasedBackedStorage) GetStoreName() string {
	return "postgres"
}

func (db DatabasedBackedStorage) Close() error {
	db.txMessageInsertStmt.Close()
	db.txInsertStmt.Close()
	db.blockMessageInsertStmt.Close()
	db.blockInsertStmt.Close()
	db.db.Close()

	return nil
}

func NewLocalhostTransactionDatabase(createTableIfNotExist bool) (*DatabasedBackedStorage, error) {

	var err error
    var dberr *DatabasedBackedStorage
	
	db, err := sql.Open("postgres", fmt.Sprintf(
		"user=%s password='%s' host='%s' dbname=%s sslmode=disable",
		User,
		Password,
		Server,
		DbName))

	if err != nil {
		return dberr, err
	}

	if err = db.Ping(); err != nil {
		fmt.Println("%v", err)
		return dberr, err
	}

	dbStore := DatabasedBackedStorage{
		db: db,
	}

	if createTableIfNotExist {
		err = dbStore.CreateTables()
		if err != nil {
			return dberr, err
		}
	}

	for i := 1; i < 200; i++ {
		insertTransactionMessage += fmt.Sprintf(", ($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4)
	}

	for i := 1; i < 100; i++ {
		insertTransaction += fmt.Sprintf(", ($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", i*14+1, i*14+2, i*14+3, i*14+4, i*14+5, i*14+6, i*14+7, i*14+8, i*14+9, i*14+10, i*14+11, i*14+12, i*14+13, i*14+14)
	}

	for i := 1; i < 200; i++ {
		insertBlockMessage += fmt.Sprintf(", ($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4)
	}

	for i := 1; i < 30; i++ {
		insertBlock += fmt.Sprintf(", ($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", i*20+1, i*20+2, i*20+3, i*20+4, i*20+5, i*20+6, i*20+7, i*20+8, i*20+9, i*20+10, i*20+11, i*20+12, i*20+13, i*20+14, i*20+15, i*20+16, i*20+17, i*20+18, i*20+19, i*20+20)
	}

	txMessageInsertStmt, err := db.Prepare(insertTransactionMessage)
	if err != nil {
		fmt.Println("%v", err)
		return dberr, err
	}

	txInsertStmt, err := db.Prepare(insertTransaction)
	if err != nil {
		fmt.Println("%+v", err)
		return dberr, err
	}

	blockMessageInsertStmt, err := db.Prepare(insertBlockMessage)
	if err != nil {
		fmt.Println("%v", err)
		return dberr, err
	}

	blockInsertStmt, err := db.Prepare(insertBlock)
	if err != nil {
		fmt.Println("%v", err)
		return dberr, err
	}

	dbStore.txMessageInsertStmt = txMessageInsertStmt
	dbStore.txInsertStmt = txInsertStmt
	dbStore.blockMessageInsertStmt = blockMessageInsertStmt
	dbStore.blockInsertStmt = blockInsertStmt

	return &dbStore, nil
}
