package endpoint

import (
	"context"
	"errors"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/go-kit/kit/endpoint"

	transport "node-pinger/btcinfo/pkg/http"
	"node-pinger/btcinfo/pkg/service"
)

type Endpoints struct {
	SyncStatusEndpoint endpoint.Endpoint
	BlockTxsEndpoint   endpoint.Endpoint
}

/*
 * Endpoint Mappings
 */

func (e Endpoints) SyncStatus(ctx context.Context) (uint, error) {
	req := transport.SyncStatusRequest{}
	resp, err := e.SyncStatusEndpoint(ctx, req)
	if err != nil {
		return 0, err
	}
	syncStatusResp := resp.(transport.SyncStatusResponse)
	if syncStatusResp.Err != "" {
		return 0, errors.New(syncStatusResp.Err)
	}
	return syncStatusResp.LatestBlock, nil
}

func (e Endpoints) BlockTxs(ctx context.Context, hashString string) (*btcjson.GetBlockVerboseResult, error) {
	req := transport.BlockTxsRequest{HashString: hashString}
	resp, err := e.BlockTxsEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	blockTxsResp := resp.(transport.BlockTxsResponse)
	if blockTxsResp.Err != "" {
		return nil, errors.New(blockTxsResp.Err)
	}
	return &blockTxsResp.BlockInfo, nil
}

/*
 * Response Makers
 */

// `MakeSyncStatusEndpoint` return the response from `Service.SyncStatus`
func MakeSyncStatusEndpoint(srv service.Service) endpoint.Endpoint {
	return func (ctx context.Context, request interface{}) (interface{}, error) {
		block, err := srv.SyncStatus(ctx)
		if err != nil {
			return transport.SyncStatusResponse{Err: err.Error()}, nil
		}
		return transport.SyncStatusResponse{LatestBlock: block}, nil
	}
}

// `MakeBlockTxsEndpoint(src Service) returns the response from `Service.BlockTxs`
func MakeBlockTxsEndpoint(srv service.Service) endpoint.Endpoint {
	return func (ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.BlockTxsRequest)
		txs, err := srv.BlockTxs(ctx, req.HashString)
		if err != nil {
			return transport.BlockTxsResponse{Err: err.Error()}, nil
		}
		return transport.BlockTxsResponse{BlockInfo: *txs}, nil
	}
}
