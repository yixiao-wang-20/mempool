package common

import "github.com/ethereum/go-ethereum/common"

type DummyWatcher struct {
}

func (d DummyWatcher) ShouldLog(from common.Address, to *common.Address) bool {
	return true
}
