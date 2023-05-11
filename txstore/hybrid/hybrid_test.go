package hybrid

import (
	"github.com/stretchr/testify/require"
	"testing"
	txstorecommon "txstore/common"
)

const testLogDir = "testlogs"

func Test_AddTransactions(t *testing.T) {
	backend, err := NewHybridStorage(testLogDir)
	require.NoError(t, err)
	defer backend.Close()

	txstorecommon.TestAddTransaction(t, backend)
}

func Test_AddBlocks(t *testing.T) {
	backend, err := NewHybridStorage(testLogDir)
	require.NoError(t, err)
	defer backend.Close()

	txstorecommon.TestAddBlock(t, backend)
}
