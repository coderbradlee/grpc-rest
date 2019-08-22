package main

import (
	"bytes"
	"context" // Use "golang.org/x/net/context" for Golang version <= 1.6
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"strings"

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
		if r.Method == "GET" {
			switch r.URL.Path {
			case "/v1/readContract":
				readContract(r)
			case "/v1/getLogs/byBlock":
				getLogsByBlock(r)
			}
		}

		h.ServeHTTP(w, r)
	})
}
func getLogsByBlock(r *http.Request) {
	kv := r.URL.Query()
	r.Method = "POST"
	data := kv.Get("topics")
	var decodeBytes []byte
	var err error
	if !strings.EqualFold(data, "") {
		decodeBytes, err = base64.StdEncoding.DecodeString(data)
		if err != nil {
			fmt.Println("a", err)
			return
		}
	}
	blockHashString := kv.Get("blockHash")
	fmt.Println(blockHashString)
	blockHashBytes, err := base64.StdEncoding.DecodeString(blockHashString)
	if err != nil {
		fmt.Println("b", err)
		return
	}
	topic := []*gw.Topics{
		&gw.Topics{
			Topic: [][]byte{decodeBytes},
		},
	}
	req := &gw.GetLogsRequest{
		Filter: &gw.LogsFilter{
			Address: []string{kv.Get("address")},
			Topics:  topic,
		},
		Lookup: &gw.GetLogsRequest_ByBlock{
			ByBlock: &gw.GetLogsByBlock{
				BlockHash: blockHashBytes,
			},
		},
		//Lookup: &iotexapi.GetLogsRequest_ByRange{
		//	ByRange: &iotexapi.GetLogsByRange{
		//		FromBlock: test.fromBlock,
		//		Count:     test.count,
		//	},
		//},
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		fmt.Println("c", err)
		return
	}
	fmt.Println(string(reqBytes))
	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBytes))
}
func readContract(r *http.Request) {
	kv := r.URL.Query()
	r.Method = "POST"
	data := kv.Get("data")
	decodeBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return
	}
	req := gw.ReadContractRequest{
		Execution: &iotextypes.Execution{
			Amount:   kv.Get("amount"),
			Contract: kv.Get("contract"),
			Data:     decodeBytes,
		},
		CallerAddress: kv.Get("callerAddress"),
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return
	}
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
