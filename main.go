package main

import (
	"bytes"
	"context" // Use "golang.org/x/net/context" for Golang version <= 1.6
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	gw "github.com/iotexproject/iotex-proto/golang/iotexapi" // Update
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
)

var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "api.testnet.iotex.one:80", "gRPC server endpoint")
)

func filter(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//if origin := r.Header.Get("Origin"); origin != "" {
		//	w.Header().Set("Access-Control-Allow-Origin", origin)
		//	if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
		//		return
		//	}
		//}

		switch r.URL.Path {
		case "/v1/readContract":
			changeQueryToBody(r)
		}

		h.ServeHTTP(w, r)
	})
}

func changeQueryToBody(r *http.Request) {
	kv := r.URL.Query()
	for k, v := range kv {
		fmt.Println(k, ":", v)
	}
	fmt.Println("empty?")
	r.Method = "POST"
	//{"execution":{"amount":"0","contract":"io1hhu3gwt5uankzl3zlp2cz8w0sl9uj336rq0334","data":"Bv3eAw=="}, "callerAddress": "io1vdtfpzkwpyngzvx7u2mauepnzja7kd5rryp0sg"}
	req := gw.ReadContractRequest{
		Execution: &iotextypes.Execution{
			Amount:   kv.Get("amount")[0],
			Contract: kv.Get("data")[0],
			Data:     kv.Get("data")[0],
		},
		CallerAddress: kv.Get("callerAddress")[0],
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return
	}
	fmt.Println("req:", string(reqBytes))
	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBytes))
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
