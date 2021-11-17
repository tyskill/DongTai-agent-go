package hookServer

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-agent/global"
	"go-agent/model/request"
	"go-agent/utils"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

// 模拟了对http包 client结构体的get的hook  此处不规范 仅为展示效果

func MyServerReq(server *http.ServeMux, w http.ResponseWriter, r *http.Request) {
	id := utils.CatGoroutineID()
	b, err := json.Marshal(r.Header)
	if err != nil {
		fmt.Println(err)
		return
	}
	baseStr := base64.StdEncoding.EncodeToString(b)
	s, _ := io.ReadAll(r.Body)
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	onlyKey, err := strconv.Atoi(strconv.Itoa(global.AgentId) + id + strconv.Itoa(int(time.Now().Unix())))
	if err != nil {
		return
	}
	global.HookGroup[id] = &request.UploadReq{
		Type:     36,
		InvokeId: onlyKey,
		Detail: request.Detail{
			AgentId: global.AgentId,
			Function: request.Function{
				Method:        r.Method,
				Url:           scheme + "://" + r.Host + r.RequestURI,
				Uri:           r.URL.Path,
				Protocol:      r.Proto,
				ClientIp:      r.RemoteAddr,
				Language:      "GO",
				ReplayRequest: false,
				ReqHeader:     baseStr,
				ReqBody:       string(s),
				QueryString:   r.URL.RawQuery,
				Pool:          []request.Pool{},
			},
		},
	}
	//fmt.Println("请求到了\n")
	//fmt.Println("method", r.Method)
	//fmt.Println("scheme", scheme)
	//fmt.Println("secure", r.TLS)
	//fmt.Println("url", scheme+"://"+r.Host+r.RequestURI)
	//fmt.Println("uri", r.URL.Path)
	//fmt.Println("protocol", r.Proto)
	//fmt.Println("clientId", r.RemoteAddr)
	//fmt.Println("language", "GO")
	//fmt.Println("replayRequest", false)
	//fmt.Println("header", baseStr)
	//fmt.Println("reqBody", string(s))
	//fmt.Println("queryString", r.URL.RawQuery)
	//utils.CatContext()
	MyServerTempReq(server, w, r)
	resH, err := json.Marshal(w.Header())
	if err != nil {
		fmt.Println(err)
		return
	}
	resBody := string(reflect.ValueOf(w).Elem().FieldByName("w").Elem().FieldByName("buf").Bytes())
	resHeader := base64.StdEncoding.EncodeToString(resH)
	global.HookGroup[id].Detail.ResHeader = resHeader
	global.HookGroup[id].Detail.ResBody = resBody
	return
}

func MyServerTempReq(server *http.ServeMux, w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 100; i++ {

	}
	return
}
