package main

import (
	"context" // Use "golang.org/x/net/context" for Golang version <= 1.6
	"flag"
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	gw "github.com/iotexproject/iotex-proto/golang/iotexapi" // Update
)

var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "api.testnet.iotex.one:80", "gRPC server endpoint")
)

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//if origin := r.Header.Get("Origin"); origin != "" {
		//	w.Header().Set("Access-Control-Allow-Origin", origin)
		//	if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
		//		return
		//	}
		//}
		fmt.Println(r.URL.Path)
		kv := r.URL.Query()
		for k, v := range kv {
			fmt.Println(k, ":", v)
		}
		fmt.Println("empty?")
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
		Handler: allowCORS(mux),
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
