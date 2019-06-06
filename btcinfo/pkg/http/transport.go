package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/btcsuite/btcd/btcjson"
)

/*
 * Mappers
 */

type SyncStatusRequest struct {}

type SyncStatusResponse struct {
	LatestBlock uint   `json:"latestBlock,omitempty"`
	Err         string `json:"err,omitempty"`
}

type BlockTxsRequest struct {
	HashString string
}

type BlockTxsResponse struct {
	BlockInfo btcjson.GetBlockVerboseResult `json:"txs,omitempty"`
	Err       string                        `json:"err,omitempty"`
}

/*
 * Decoders
 */

func DecodeSyncStatusRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req SyncStatusRequest
	return req, nil
}

func DecodeBlockTxsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req BlockTxsRequest
	hashString := r.URL.Query().Get("blockHash")
	if hashString == "" {
		return nil, errors.New("no parameter `blockHash` in the query")
	}
	req.HashString = hashString
	return req, nil
}

/*
 * Encoders
 */
func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}