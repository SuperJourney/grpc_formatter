package grpcformatter

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/SuperJourney/cache_wrapper"
	"github.com/golang/protobuf/proto" //nolint
	"google.golang.org/genproto/googleapis/rpc/status"
	gstatus "google.golang.org/grpc/status"
)

// nolint
// nolint
type M interface {
	cache_wrapper.RequestFormatter
}

type GrpcFormatter struct {
	uniqKey string
}

func NewGrpcFormatter(uniqueKey string) *GrpcFormatter {
	return &GrpcFormatter{
		uniqKey: uniqueKey,
	}
}
func (f *GrpcFormatter) GetUniqKey(reqs ...interface{}) ([]byte, error) {
	// gprc request 结构固定  reqs[0] context.Context ,reqs[1]:  protoV2.Message
	_, ok := reqs[0].(context.Context)
	if !ok {
		return nil, fmt.Errorf("cache, ctx type err")
	}

	req, ok := reqs[1].(proto.Message)
	if !ok {
		return nil, fmt.Errorf("cache, req type err")
	}

	pr, _ := proto.Marshal(req)
	m := md5.New() //nolint
	m.Write([]byte(pr))
	sig := hex.EncodeToString(m.Sum(nil))
	return []byte(fmt.Sprintf("%s:common-cache:%v", f.uniqKey, sig)), nil
}

type RespCache struct {
	Message    []byte         `json:"message,omitempty"`
	GrpcStatus *status.Status `json:"grpc_status,omitempty"`
	ErrMsg     string         `json:"err,omitempty"`
}

func (f *GrpcFormatter) MarshalWrapper(respes ...interface{}) ([]byte, error) {
	// gprc response结构固定 resp[0]:  protoV2.Message , resp[1]: error
	var resp proto.Message = nil
	if respes[0] != nil {
		var ok bool
		resp, ok = respes[0].(proto.Message)
		if !ok {
			return nil, fmt.Errorf("cache, resp type err")
		}
	}

	var resperr error
	if respes[1] == nil {
		resperr = nil
	} else {
		resperr = respes[1].(error)
	}

	var respCacheMsg []byte

	if resp == nil || !proto.MessageV2(resp).ProtoReflect().IsValid() {
		respCacheMsg = nil
	} else {
		var err error
		respCacheMsg, err = proto.Marshal(resp)
		if err != nil {
			return nil, fmt.Errorf("cache, proto Marshal err:%v", err)
		}
	}

	respCache := &RespCache{
		Message: respCacheMsg,
	}

	if resperr != nil {
		s, ok := gstatus.FromError(resperr)
		if !ok {
			respCache.ErrMsg = resperr.Error()
		} else {
			respCache.GrpcStatus = s.Proto()
		}
	}

	respCacheByte, err := json.Marshal(respCache)
	if err != nil {
		return nil, fmt.Errorf("cache, json Marshal  err:%v", err)
	}

	return respCacheByte, nil
}

func (f *GrpcFormatter) UnMarshalWrapper(respStr []byte, respes ...interface{}) ([]interface{}, error) {
	var respCache RespCache
	err := json.Unmarshal(respStr, &respCache)
	if err != nil {
		return nil, fmt.Errorf("cache, json Unmarshal  err:%v", err)
	}

	resp := respes[0]

	var ret []interface{}
	var newResp interface{}
	if respCache.Message != nil {
		// 转换成proto信息
		err = proto.Unmarshal(respCache.Message, resp.(proto.Message))
		if err != nil {
			return nil, fmt.Errorf("cache, proto Unmarshal  err:%v", err)
		}
		newResp = resp
	}
	var respErr error = nil
	if respCache.GrpcStatus != nil {
		respErr = gstatus.ErrorProto(respCache.GrpcStatus)
	} else {
		if respCache.ErrMsg != "" {
			respErr = fmt.Errorf(respCache.ErrMsg)
		}
	}

	ret = append(ret, newResp, respErr)
	return ret, nil
}
