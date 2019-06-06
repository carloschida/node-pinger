package service

import (
	"context"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"node-pinger/btcinfo/pkg/endpoint"
	transport "node-pinger/btcinfo/pkg/http"
)

func NewHTTPServer(ctv context.Context, endpoints endpoint.Endpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(commonMiddleware)

	r.Methods("GET").Path("/").Handler(httptransport.NewServer(
		endpoints.SyncStatusEndpoint,
		transport.DecodeSyncStatusRequest,
		transport.EncodeResponse,
	))

	r.Methods("GET").Path("/blocks").Handler(httptransport.NewServer(
		endpoints.BlockTxsEndpoint,
		transport.DecodeBlockTxsRequest,
		transport.EncodeResponse,
	))

	return r
}

// TODO move this middleware to where middleware belong (although it's only 1...)
func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
