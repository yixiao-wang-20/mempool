package main

import (
	_ "github.com/lib/pq"
	"txstore/db"
)

func main() {
	store, err := db.NewLocalhostTransactionDatabase(
		"postgres",
		"2312jsd2jklq99287nsxz",
		"localhost",
		"tx_log",
		true)
	if err != nil {
		panic(err)
	}
	defer store.Close()
}
