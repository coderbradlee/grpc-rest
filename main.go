package main

import (
	"bytes"
	"context" // Use "golang.org/x/net/context" for Golang version <= 1.6
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"strconv"
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
			case "/v1/getLogs/byRange":
				getLogsByRange(r)
			case "/v1/streamLogs":
				streamlogs(r)
			case "/v1/getBlockMetas/byIndex":
				getBlockMetas(r, true)
			case "/v1/getBlockMetas/byHash":
				getBlockMetas(r, false)
			case "/v1/sendAction/transfer":
				sendTransfer(r)
			}
		}

		h.ServeHTTP(w, r)
	})
}
func sendTransfer(r *http.Request) {
	//curl -X GET localhost:8081/v1/sendAction?version=1\&nonce=2\&gasLimit=10000\&gasPrice=10\&amount=100\&recipient=io1sxm6zl56um2c3ntq5fwqjar4za5ka560x53muy\&senderPubKey=BOk7WxyPumkmNlKkg61VMY5O7VtRIjFMt/2wd9jHKVCXzsku5QsRCNx0lalyDlkh5W0wSON6vmpnFtfJuRPp8uY=\&signature=9mrqFBggiRocZhkRVUswxs83NaEFNdEYYczI8049vlovHEP4YMQz+3Isznc3CrYaJxAbc2PTIz7y2meerJ8bHAA=
	//{"action": {"core": {"version": 1, "nonce": 2, "gasLimit": 10000, "gasPrice": "10", "transfer": {"amount": "100", "recipient": "io1sxm6zl56um2c3ntq5fwqjar4za5ka560x53muy"}}, "senderPubKey": "BOk7WxyPumkmNlKkg61VMY5O7VtRIjFMt/2wd9jHKVCXzsku5QsRCNx0lalyDlkh5W0wSON6vmpnFtfJuRPp8uY=", "signature": "9mrqFBggiRocZhkRVUswxs83NaEFNdEYYczI8049vlovHEP4YMQz+3Isznc3CrYaJxAbc2PTIz7y2meerJ8bHAA="}}
	kv := r.URL.Query()
	r.Method = "POST"
	version, err := strconv.ParseUint(kv.Get("version"), 10, 32)
	if err != nil {
		fmt.Println("version:", err)
		return
	}
	nonce, err := strconv.ParseUint(kv.Get("nonce"), 10, 64)
	if err != nil {
		fmt.Println("nonce:", err)
		return
	}
	gasLimit, err := strconv.ParseUint(kv.Get("gasLimit"), 10, 64)
	if err != nil {
		fmt.Println("gaslimit:", err)
		return
	}
	type actionCore struct {
		Version  uint32               `json:"version,omitempty"`
		Nonce    uint64               `json:"nonce,omitempty"`
		GasLimit uint64               `json:"gasLimit,omitempty"`
		GasPrice string               `json:"gasPrice,omitempty"`
		Transfer *iotextypes.Transfer `json:"transfer,omitempty"`
	}
	type sendActionStruct struct {
		Core         *actionCore `json:"core,omitempty"`
		SenderPubKey []byte      `json:"senderPubKey,omitempty"`
		Signature    []byte      `json:"signature,omitempty"`
	}

	senderPubKey, err := base64.StdEncoding.DecodeString(kv.Get("senderPubKey"))
	if err != nil {
		fmt.Println("b", err)
		return
	}
	signature, err := base64.StdEncoding.DecodeString(kv.Get("signature"))
	if err != nil {
		fmt.Println("b", err)
		return
	}
	payload, err := base64.StdEncoding.DecodeString(kv.Get("payload"))
	if err != nil {
		fmt.Println("b", err)
		return
	}
	req := &sendActionStruct{
		Core: &core{
			Version:  uint32(version),
			Nonce:    nonce,
			GasLimit: gasLimit,
			GasPrice: kv.Get("gasPrice"),
			Transfer: &iotextypes.Transfer{
				Amount:    kv.Get("amount"),
				Recipient: kv.Get("recipient"),
				Payload:   payload,
			},
		},
		SenderPubKey: senderPubKey,
		Signature:    signature,
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		fmt.Println("c", err)
		return
	}
	fmt.Println(string(reqBytes))
	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBytes))
	r.URL.Path = "/v1/sendAction"
}

func getBlockMetas(r *http.Request, byIndex bool) {
	kv := r.URL.Query()
	r.Method = "POST"
	var reqBytes []byte
	var err error
	if byIndex {
		type byIndexStruct struct {
			ByIndex *gw.GetBlockMetasByIndexRequest `json:"byIndex,omitempty"`
		}
		start, err := strconv.ParseUint(kv.Get("start"), 10, 64)
		if err != nil {
			return
		}
		count, err := strconv.ParseUint(kv.Get("count"), 10, 64)
		if err != nil {
			return
		}
		req := &byIndexStruct{
			ByIndex: &gw.GetBlockMetasByIndexRequest{
				Start: start,
				Count: count,
			},
		}
		reqBytes, err = json.Marshal(req)
		if err != nil {
			fmt.Println("c", err)
			return
		}
	} else {
		type byHashStruct struct {
			ByHash *gw.GetBlockMetaByHashRequest `json:"byHash,omitempty"`
		}
		req := &byHashStruct{
			ByHash: &gw.GetBlockMetaByHashRequest{
				BlkHash: kv.Get("blkHash"),
			},
		}
		reqBytes, err = json.Marshal(req)
		if err != nil {
			fmt.Println("c", err)
			return
		}
	}

	fmt.Println(string(reqBytes))
	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBytes))
	r.URL.Path = "/v1/getBlockMetas"
}
func streamlogs(r *http.Request) {
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

	var topic []*gw.Topics
	if len(decodeBytes) != 0 {
		topic = []*gw.Topics{
			&gw.Topics{
				Topic: [][]byte{decodeBytes},
			},
		}
	}
	type reqStruct struct {
		Filter *gw.LogsFilter `json:"filter,omitempty"`
	}
	req := &reqStruct{
		Filter: &gw.LogsFilter{
			Address: []string{kv.Get("address")},
			Topics:  topic,
		},
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		fmt.Println("c", err)
		return
	}
	fmt.Println(string(reqBytes))
	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBytes))
}
func getLogsByRange(r *http.Request) {
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

	var topic []*gw.Topics
	if len(decodeBytes) != 0 {
		topic = []*gw.Topics{
			&gw.Topics{
				Topic: [][]byte{decodeBytes},
			},
		}
	}
	from := kv.Get("fromBlock")
	fromUint64, err := strconv.ParseUint(from, 10, 64)
	if err != nil {
		return
	}
	count := kv.Get("count")
	countUint64, err := strconv.ParseUint(count, 10, 64)
	if err != nil {
		return
	}
	type reqStruct struct {
		Filter  *gw.LogsFilter     `json:"filter,omitempty"`
		ByRange *gw.GetLogsByRange `json:"byRange,omitempty"`
	}
	req := &reqStruct{
		Filter: &gw.LogsFilter{
			Address: []string{kv.Get("address")},
			Topics:  topic,
		},
		ByRange: &gw.GetLogsByRange{
			FromBlock: fromUint64,
			Count:     countUint64,
		},
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		fmt.Println("c", err)
		return
	}
	fmt.Println(string(reqBytes))
	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBytes))
	r.URL.Path = "/v1/getLogs"
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
	replaced := strings.ReplaceAll(blockHashString, " ", "+")
	fmt.Println(replaced)
	blockHashBytes, err := base64.StdEncoding.DecodeString(replaced)
	if err != nil {
		fmt.Println("b", err)
		return
	}
	var topic []*gw.Topics
	if len(decodeBytes) != 0 {
		topic = []*gw.Topics{
			&gw.Topics{
				Topic: [][]byte{decodeBytes},
			},
		}
	}

	type reqStruct struct {
		Filter  *gw.LogsFilter     `json:"filter,omitempty"`
		ByBlock *gw.GetLogsByBlock `json:"byBlock,omitempty"`
	}
	req := &reqStruct{
		Filter: &gw.LogsFilter{
			Address: []string{kv.Get("address")},
			Topics:  topic,
		},
		ByBlock: &gw.GetLogsByBlock{
			BlockHash: blockHashBytes,
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
	r.URL.Path = "/v1/getLogs"
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
