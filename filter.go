package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	gw "github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
)

var (
	getActionsPath                   = "/v1/getActions"
	estimateGasForActionPath         = "/v1/estimateGasForAction"
	getBlockMetasPath                = "/v1/getBlockMetas"
	estimateActionGasConsumptionPath = "/v1/estimateActionGasConsumption"
	getLogsPath                      = "/v1/getLogs"

	readContractPath                          = "/v1/readContract"
	getLogsByBlockPath                        = "/v1/getLogs/byBlock"
	getLogsByRangePath                        = "/v1/getLogs/byRange"
	streamLogsPath                            = "/v1/streamLogs"
	getBlockMetasByIndexPath                  = "/v1/getBlockMetas/byIndex"
	getBlockMetasByHashPath                   = "/v1/getBlockMetas/byHash"
	sendActionTransferPath                    = "/v1/sendAction/transfer"
	sendActionExecutionPath                   = "/v1/sendAction/execution"
	getActionsByIndexPath                     = "/v1/getActions/byIndex"
	getActionsByHashPath                      = "/v1/getActions/byHash"
	getActionsByAddrPath                      = "/v1/getActions/byAddr"
	getActionsUnconfirmedByAddrPath           = "/v1/getActions/unconfirmedByAddr"
	getActionsByBlkPath                       = "/v1/getActions/byBlk"
	estimateGasForActionTransferPath          = "/v1/estimateGasForAction/transfer"
	estimateGasForActionExecutionPath         = "/v1/estimateGasForAction/execution"
	estimateActionGasConsumptionTransferPath  = "/v1/estimateActionGasConsumption/transfer"
	estimateActionGasConsumptionExecutionPath = "/v1/estimateActionGasConsumption/execution"
)

func estimateTransferGasConsumption(r *http.Request) {
	kv := r.URL.Query()
	r.Method = "POST"
	dataString := kv.Get("payload")
	dataString = strings.ReplaceAll(dataString, " ", "+")
	data, err := base64.StdEncoding.DecodeString(dataString)
	if err != nil {
		return
	}
	type estimateRequest struct {
		Transfer      *iotextypes.Transfer `json:"transfer,omitempty"`
		CallerAddress string               `json:"callerAddress,omitempty"`
	}

	req := &estimateRequest{
		Transfer: &iotextypes.Transfer{
			Amount:    kv.Get("amount"),
			Recipient: kv.Get("recipient"),
			Payload:   data,
		},
		CallerAddress: kv.Get("callerAddress"),
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBytes))
	r.URL.Path = estimateActionGasConsumptionPath
}
func estimateExecutionGasConsumption(r *http.Request) {
	kv := r.URL.Query()
	r.Method = "POST"
	dataString := kv.Get("data")
	dataString = strings.ReplaceAll(dataString, " ", "+")
	data, err := base64.StdEncoding.DecodeString(dataString)
	if err != nil {
		return
	}
	type estimateRequest struct {
		Execution     *iotextypes.Execution `json:"execution,omitempty"`
		CallerAddress string                `json:"callerAddress,omitempty"`
	}

	req := &estimateRequest{
		Execution: &iotextypes.Execution{
			Amount:   kv.Get("amount"),
			Contract: kv.Get("contract"),
			Data:     data,
		},
		CallerAddress: kv.Get("callerAddress"),
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBytes))
	r.URL.Path = estimateActionGasConsumptionPath
}
func getActionsByBlk(r *http.Request) {
	kv := r.URL.Query()
	r.Method = "POST"
	var reqBytes []byte
	var err error
	type byAddrStruct struct {
		ByBlk *gw.GetActionsByBlockRequest `json:"byBlk,omitempty"`
	}
	start, err := strconv.ParseUint(kv.Get("start"), 10, 64)
	if err != nil {
		return
	}
	count, err := strconv.ParseUint(kv.Get("count"), 10, 64)
	if err != nil {
		return
	}
	req := &byAddrStruct{
		ByBlk: &gw.GetActionsByBlockRequest{
			BlkHash: kv.Get("blkHash"),
			Start:   start,
			Count:   count,
		},
	}
	reqBytes, err = json.Marshal(req)
	if err != nil {
		return
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBytes))
	r.URL.Path = getActionsPath
}
func getActionsUnconfirmedByAddr(r *http.Request) {
	kv := r.URL.Query()
	r.Method = "POST"
	var reqBytes []byte
	var err error
	type byAddrStruct struct {
		UnconfirmedByAddr *gw.GetUnconfirmedActionsByAddressRequest `json:"unconfirmedByAddr,omitempty"`
	}
	start, err := strconv.ParseUint(kv.Get("start"), 10, 64)
	if err != nil {
		return
	}
	count, err := strconv.ParseUint(kv.Get("count"), 10, 64)
	if err != nil {
		return
	}
	req := &byAddrStruct{
		UnconfirmedByAddr: &gw.GetUnconfirmedActionsByAddressRequest{
			Address: kv.Get("address"),
			Start:   start,
			Count:   count,
		},
	}
	reqBytes, err = json.Marshal(req)
	if err != nil {
		return
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBytes))
	r.URL.Path = getActionsPath
}
func getActionsByAddr(r *http.Request) {
	kv := r.URL.Query()
	r.Method = "POST"
	var reqBytes []byte
	var err error
	type byAddrStruct struct {
		ByAddr *gw.GetActionsByAddressRequest `json:"byAddr,omitempty"`
	}
	start, err := strconv.ParseUint(kv.Get("start"), 10, 64)
	if err != nil {
		return
	}
	count, err := strconv.ParseUint(kv.Get("count"), 10, 64)
	if err != nil {
		return
	}
	req := &byAddrStruct{
		ByAddr: &gw.GetActionsByAddressRequest{
			Address: kv.Get("address"),
			Start:   start,
			Count:   count,
		},
	}
	reqBytes, err = json.Marshal(req)
	if err != nil {
		return
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBytes))
	r.URL.Path = getActionsPath
}
func getActionsByHash(r *http.Request) {
	kv := r.URL.Query()
	r.Method = "POST"
	var reqBytes []byte
	var err error
	type byHashStruct struct {
		ByHash *gw.GetActionByHashRequest `json:"byHash,omitempty"`
	}
	var chekpending bool
	if strings.EqualFold(kv.Get("checkPending"), "true") {
		chekpending = true
	}
	req := &byHashStruct{
		ByHash: &gw.GetActionByHashRequest{
			ActionHash:   kv.Get("actionHash"),
			CheckPending: chekpending,
		},
	}
	reqBytes, err = json.Marshal(req)
	if err != nil {
		return
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBytes))
	r.URL.Path = getActionsPath
}
func getActionsByIndex(r *http.Request) {
	kv := r.URL.Query()
	r.Method = "POST"
	var reqBytes []byte
	var err error
	type byIndexStruct struct {
		ByIndex *gw.GetActionsByIndexRequest `json:"byIndex,omitempty"`
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
		ByIndex: &gw.GetActionsByIndexRequest{
			Start: start,
			Count: count,
		},
	}
	reqBytes, err = json.Marshal(req)
	if err != nil {
		return
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBytes))
	r.URL.Path = getActionsPath
}
func sendExecution(r *http.Request, estimate bool) {
	kv := r.URL.Query()
	r.Method = "POST"
	version, err := strconv.ParseUint(kv.Get("version"), 10, 32)
	if err != nil {
		return
	}
	nonce, err := strconv.ParseUint(kv.Get("nonce"), 10, 64)
	if err != nil {
		return
	}
	gasLimit, err := strconv.ParseUint(kv.Get("gasLimit"), 10, 64)
	if err != nil {
		return
	}
	type actionCore struct {
		Version   uint32                `json:"version,omitempty"`
		Nonce     uint64                `json:"nonce,omitempty"`
		GasLimit  uint64                `json:"gasLimit,omitempty"`
		GasPrice  string                `json:"gasPrice,omitempty"`
		Execution *iotextypes.Execution `json:"execution,omitempty"`
	}
	type sendActionStruct struct {
		Core         *actionCore `json:"core,omitempty"`
		SenderPubKey []byte      `json:"senderPubKey,omitempty"`
		Signature    []byte      `json:"signature,omitempty"`
	}

	senderPubKeyString := kv.Get("senderPubKey")
	senderPubKeyString = strings.ReplaceAll(senderPubKeyString, " ", "+")
	senderPubKey, err := base64.StdEncoding.DecodeString(senderPubKeyString)
	if err != nil {
		return
	}
	signatureString := kv.Get("signature")
	signatureString = strings.ReplaceAll(signatureString, " ", "+")
	signature, err := base64.StdEncoding.DecodeString(signatureString)
	if err != nil {
		return
	}
	dataString := kv.Get("data")
	dataString = strings.ReplaceAll(dataString, " ", "+")
	data, err := base64.StdEncoding.DecodeString(dataString)
	if err != nil {
		return
	}
	type SendActionRequest struct {
		Action *sendActionStruct `json:"action,omitempty"`
	}
	req := &SendActionRequest{
		Action: &sendActionStruct{
			Core: &actionCore{
				Version:  uint32(version),
				Nonce:    nonce,
				GasLimit: gasLimit,
				GasPrice: kv.Get("gasPrice"),
				Execution: &iotextypes.Execution{
					Amount:   kv.Get("amount"),
					Contract: kv.Get("contract"),
					Data:     data,
				},
			},
			SenderPubKey: senderPubKey,
			Signature:    signature,
		},
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBytes))
	r.URL.Path = "/v1/sendAction"
	if estimate {
		r.URL.Path = estimateGasForActionPath
	}
}
func sendTransfer(r *http.Request, estimate bool) {
	kv := r.URL.Query()
	r.Method = "POST"
	version, err := strconv.ParseUint(kv.Get("version"), 10, 32)
	if err != nil {
		return
	}
	nonce, err := strconv.ParseUint(kv.Get("nonce"), 10, 64)
	if err != nil {
		return
	}
	gasLimit, err := strconv.ParseUint(kv.Get("gasLimit"), 10, 64)
	if err != nil {
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

	senderPubKeyString := kv.Get("senderPubKey")
	senderPubKeyString = strings.ReplaceAll(senderPubKeyString, " ", "+")
	senderPubKey, err := base64.StdEncoding.DecodeString(senderPubKeyString)
	if err != nil {
		return
	}
	signatureString := kv.Get("signature")
	signatureString = strings.ReplaceAll(signatureString, " ", "+")
	signature, err := base64.StdEncoding.DecodeString(signatureString)
	if err != nil {
		return
	}
	payloadString := kv.Get("payload")
	payloadString = strings.ReplaceAll(payloadString, " ", "+")
	payload, err := base64.StdEncoding.DecodeString(payloadString)
	if err != nil {
		return
	}
	type SendActionRequest struct {
		Action *sendActionStruct `json:"action,omitempty"`
	}
	req := &SendActionRequest{
		Action: &sendActionStruct{
			Core: &actionCore{
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
		},
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBytes))
	r.URL.Path = "/v1/sendAction"
	if estimate {
		r.URL.Path = estimateGasForActionPath
	}
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
			return
		}
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBytes))
	r.URL.Path = getBlockMetasPath
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
		return
	}
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
		return
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBytes))
	r.URL.Path = getLogsPath
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
			return
		}
	}
	blockHashString := kv.Get("blockHash")
	replaced := strings.ReplaceAll(blockHashString, " ", "+")

	blockHashBytes, err := base64.StdEncoding.DecodeString(replaced)
	if err != nil {
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
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBytes))
	r.URL.Path = getLogsPath
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
