package service

import (
	"context"
	"testing"
)

func TestSmoke(t *testing.T) {
	srv, ctx := setup()

	// Just test the existence of the service and the context
	ok := srv != nil && ctx != nil
	if !ok {
		t.Errorf("Missing service or context")
	}
}

func TestSyncStatus(t *testing.T) {
	srv, ctx := setup()

	lb, err := srv.SyncStatus(ctx)
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	/* If it reaches a proper node in the right way,
	 * it is highly likely that it has at least one
	 * block
	 */
	ok := lb > 0
	if !ok {
		t.Errorf("Expected to have at least 1 block")
	}
}

func TestBlockTxs(t *testing.T) {
	srv, ctx := setup()

	// https://blockstream.info/block/00000000000000000008fb52139fb6d289136883c1e3222630af3ebd81495fec
	hashString  := "00000000000000000008fb52139fb6d289136883c1e3222630af3ebd81495fec"
	blockHeight := 579375

	info, err := srv.BlockTxs(ctx, hashString)
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	ok := info.Height == int64(blockHeight)
	if !ok {
		t.Errorf("Obtained info did not match expected block height")
	}

}

func setup() (srv Service, ctx context.Context) {
	return NewService(), context.Background()
}
