package executor

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hongzhaomin/hzm-job/client/api"
	"github.com/hongzhaomin/hzm-job/client/model"
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

var _ api.JobClientApi = (*HttpJobClient)(nil)

type HttpJobClient struct{}

func (my *HttpJobClient) HeatBeat() {
	fmt.Println("heat beat job")
}

func (my *HttpJobClient) JobHandle(req *api.JobHandleReq) {
	marshal, _ := json.Marshal(req)
	fmt.Println("job handle:", string(marshal))
}

func (my *HttpJobClient) Start() {
	http.HandleFunc(path.Join(baseUrl, heartBeat), func(w http.ResponseWriter, r *http.Request) {
		commonHandlerFun(w, r, func(param any) (any, error) {
			my.HeatBeat()
			return "ok", nil
		})
	})

	http.HandleFunc(path.Join(baseUrl, jobHandle), func(w http.ResponseWriter, r *http.Request) {
		commonHandlerFun[api.JobHandleReq](w, r, func(param api.JobHandleReq) (any, error) {
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
		json.NewEncoder(res).Encode(model.Fail[any](err.Error()))
	} else {
		json.NewEncoder(res).Encode(model.Ok2[any](data))
	}

}
