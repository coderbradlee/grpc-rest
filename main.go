package main

import (
	"context"
	"flag"
	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	gw "github.com/iotexproject/iotex-proto/golang/iotexapi"
	"google.golang.org/grpc"
)

var (
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "api.testnet.iotex.one:80", "gRPC server endpoint")
)

func filter(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			switch r.URL.Path {
			case readContractPath:
				readContract(r)
			case getLogsByBlockPath:
				getLogsByBlock(r)
			case getLogsByRangePath:
				getLogsByRange(r)
			case streamLogsPath:
				streamlogs(r)
			case getBlockMetasByIndexPath:
				getBlockMetas(r, true)
			case getBlockMetasByHashPath:
				getBlockMetas(r, false)
			case sendActionTransferPath:
				sendTransfer(r, false)
			case sendActionExecutionPath:
				sendExecution(r, false)
			case getActionsByIndexPath:
				getActionsByIndex(r)
			case getActionsByHashPath:
				getActionsByHash(r)
			case getActionsByAddrPath:
				getActionsByAddr(r)
			case getActionsUnconfirmedByAddrPath:
				getActionsUnconfirmedByAddr(r)
			case getActionsByBlkPath:
				getActionsByBlk(r)
			case estimateGasForActionTransferPath:
				sendTransfer(r, true)
			case estimateGasForActionExecutionPath:
				sendExecution(r, true)
			case estimateActionGasConsumptionTransferPath:
				estimateTransferGasConsumption(r)
			case estimateActionGasConsumptionExecutionPath:
				estimateExecutionGasConsumption(r)
			}
		}

		h.ServeHTTP(w, r)
	})
}

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gw.RegisterAPIServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	//return http.ListenAndServe(":8081", mux)
	s := &http.Server{
		Addr:    ":8081",
		Handler: filter(mux),
	}
	return s.ListenAndServe()

}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
