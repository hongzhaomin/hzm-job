package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hongzhaomin/hzm-job/core/sdk"
	"io"
	"log/slog"
	"net/http"
	"path"
)

const (
	baseUrl   = "/api"
	heartBeat = "heart-beat"
	jobHandle = "job-handle"
)

var _ JobClientApi = (*HttpJobClient)(nil)

type HttpJobClient struct{}

func (my *HttpJobClient) HeatBeat() {
	fmt.Println("heat beat job")
}

func (my *HttpJobClient) JobHandle(req *JobHandleReq) {
	marshal, _ := json.Marshal(req)
	fmt.Println("job handle:", string(marshal))
}

func (my *HttpJobClient) Start() {
	http.HandleFunc(path.Join(baseUrl, heartBeat), func(w http.ResponseWriter, r *http.Request) {
		commonHandlerFun(w, r, func(param any) (any, error) {
			my.HeatBeat()
			return true, nil
		})
	})

	http.HandleFunc(path.Join(baseUrl, jobHandle), func(w http.ResponseWriter, r *http.Request) {
		commonHandlerFun[JobHandleReq](w, r, func(param JobHandleReq) (any, error) {
			my.JobHandle(&param)
			return true, nil
		})
	})

	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		slog.Error("ListenAndServe: ", err)
	}
}

func commonHandlerFun[P any](res http.ResponseWriter, req *http.Request, fn func(param P) (any, error)) {
	defer func() {
		if err := recover(); err != nil {
			slog.Error("http server panic: ", err)
			if err2, ok := err.(error); ok {
				http.Error(res, err2.Error(), http.StatusInternalServerError)
			} else {
				http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}
	}()

	if req.Method != http.MethodPost {
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var param P
	if err := json.NewDecoder(req.Body).Decode(&param); err != nil && !errors.Is(err, io.EOF) {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := fn(param)
	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "application/json")
	if err != nil {
		_ = json.NewEncoder(res).Encode(sdk.Fail[any](err.Error()))
	} else {
		_ = json.NewEncoder(res).Encode(sdk.Ok2[any](data))
	}

}
