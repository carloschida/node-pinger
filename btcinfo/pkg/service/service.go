package service

import (
	"context"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/tkanos/gonfig"
)

type Service interface {
	SyncStatus(ctx context.Context) (uint, error)
	BlockTxs(ctx context.Context, hashString string) (*btcjson.GetBlockVerboseResult, error)
}

type infoService struct {}

func (infoService) SyncStatus(ctx context.Context) (uint, error) {
	lb, err := getLatestBlock()
	if err != nil {
		return 0, err
	}
	return lb, nil
}

func (infoService) BlockTxs(ctx context.Context, hashString string) (*btcjson.GetBlockVerboseResult, error) {
	txs, err := getTxsFromBlockHash(hashString)
	if err != nil {
		return nil, err
	}
	return txs, nil
}

func NewService() Service {
	return infoService{}
}

/*
 * Business Logic
 */

// TODO Move config to a suitable place and make singleton

type rpcConf struct {
	Host string `env:"RPC_HOST"`
	User string `env:"RPC_USER"`
	Pass string `env:"RPC_PASS"`
}

func getLatestBlock() (uint, error) {
	// Retrieve local conf from file
	conf := rpcConf{}
	err := gonfig.GetConf("./rpcConf.json", &conf)
	if err != nil {
		panic(err)
	}

	// Connect to local bitcoin core RPC server using HTTP POST mode.
	connCfg := &rpcclient.ConnConfig{
		Host:         conf.Host,
		User:         conf.User,
		Pass:         conf.Pass,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return 0, err
	}
	defer client.Shutdown()

	// Get the current block count
	blockCount, err := client.GetBlockCount()
	if err != nil {
		return 0, err
	}

	return uint(blockCount), nil
}

func getTxsFromBlockHash(blockHash string) (*btcjson.GetBlockVerboseResult, error) {
	// Retrieve local conf from file
	conf := rpcConf{}
	err := gonfig.GetConf("./rpcConf.json", &conf)
	if err != nil {
		panic(err)
	}

	// Connect to local bitcoin core RPC server using HTTP POST mode.
	connCfg := &rpcclient.ConnConfig{
		Host:         conf.Host,
		User:         conf.User,
		Pass:         conf.Pass,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return &btcjson.GetBlockVerboseResult{}, err
	}
	defer client.Shutdown()

	hash, err := chainhash.NewHashFromStr(blockHash)
	if err != nil {
		return &btcjson.GetBlockVerboseResult{}, err
	}

	/*
	 * Using `GetBlockVerbose` is compatible with
	 */
	blockVerbose, err := client.GetBlockVerbose(hash)
	if err != nil {
		return &btcjson.GetBlockVerboseResult{}, err
	}

	return blockVerbose, nil
}
