package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/hongzhaomin/hzm-job/client/internal"
	"github.com/hongzhaomin/hzm-job/client/internal/anno"
	"github.com/hongzhaomin/hzm-job/client/internal/consts"
	"github.com/hongzhaomin/hzm-job/client/internal/global"
	"github.com/hongzhaomin/hzm-job/client/internal/prop"
	"github.com/hongzhaomin/hzm-job/core/config"
	"github.com/hongzhaomin/hzm-job/core/ezconfig"
	"github.com/hongzhaomin/hzm-job/core/sdk"
	"io"
	"net/http"
	"path"
)

var _ JobClientApi = (*HttpJobClient)(nil)

type HttpJobClient struct {
	server *http.Server
}

func (my *HttpJobClient) HeatBeat() {
	global.SingletonPool().Log.Debug("收到调度中心心跳检测")
}

func (my *HttpJobClient) JobHandle(req *JobHandleReq) {
	if req == nil {
		panic(errors.New("schedule job params is nil"))
	}

	jobName := req.JobName
	if jobName == nil {
		panic(errors.New("job name is nil"))
	}

	// 任务异步化，执行完进行服务端回调
	ctx, cancel := context.WithCancel(context.Background())
	// 将任务日志与取消上下文绑定存储
	global.SingletonPool().JobCancelCtx.Put(req.LogId, cancel)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				jobParams, _ := json.Marshal(req)
				global.SingletonPool().Log.Error("[hzm-job]任务回调异常",
					"jobParams", string(jobParams),
					"err", err)
			}
		}()
		defer global.SingletonPool().JobCancelCtx.CancelAndRemove(req.LogId)
		jobHandler := internal.DefaultJobRegister().GetJob(*jobName)
		jobExeErr := anno.CallDoHandle(jobHandler, ctx, req.JobParams)
		if jobExeErr != nil {
			global.SingletonPool().Log.Error("hzm-job: 任务执行失败",
				"jobName", *jobName, "jobParams", req.JobParams, "err", jobExeErr)

			// 失败回调
			err := internal.Callback2Admin(*req.LogId, 500, jobExeErr.Error())
			if err != nil {
				global.SingletonPool().Log.Error("hzm-job: 执行器任务回调失败",
					"logId", *req.LogId,
					"jobName", *jobName,
					"jobExeResult", "失败",
					"jobExeMsg", jobExeErr,
					"err", err)
			}
		} else {
			// 成功回调
			err := internal.Callback2Admin(*req.LogId, 200, "执行成功")
			if err != nil {
				global.SingletonPool().Log.Error("hzm-job: 执行器任务回调失败",
					"logId", *req.LogId,
					"jobName", *jobName,
					"jobExeResult", "成功",
					"err", err)
			}
		}
	}()
}

func (my *HttpJobClient) CancelJob(req *JobCancelReq) {
	if req == nil {
		panic(errors.New("cancel job params is nil"))
	}
	if req.LogId == nil {
		panic(errors.New("logId is nil"))
	}
	global.SingletonPool().JobCancelCtx.CancelAndRemove(req.LogId)
}

// Start 启动 http 服务
func (my *HttpJobClient) Start() {
	http.HandleFunc(path.Join(consts.BaseUrl, consts.HeartBeatUrl), func(w http.ResponseWriter, r *http.Request) {
		commonHandlerFun(w, r, func(param any) (any, error) {
			my.HeatBeat()
			return true, nil
		})
	})

	http.HandleFunc(path.Join(consts.BaseUrl, consts.JobHandleUrl), func(w http.ResponseWriter, r *http.Request) {
		commonHandlerFun[JobHandleReq](w, r, func(param JobHandleReq) (any, error) {
			my.JobHandle(&param)
			return true, nil
		})
	})

	http.HandleFunc(path.Join(consts.BaseUrl, consts.JobCancelUrl), func(w http.ResponseWriter, r *http.Request) {
		commonHandlerFun[JobCancelReq](w, r, func(param JobCancelReq) (any, error) {
			my.CancelJob(&param)
			return true, nil
		})
	})

	my.server = &http.Server{Addr: ":" + ezconfig.Get[*prop.HzmJobConfigBean]().Port, Handler: nil}
	// 注册服务下线函数
	my.server.RegisterOnShutdown(func() {
		err := internal.Offline2Admin()
		if err != nil {
			global.SingletonPool().Log.Error("hzm-job: 执行器自动下线失败", "err", err)
		}
	})
	err := my.server.ListenAndServe()
	//err := http.ListenAndServe(":8888", nil)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		global.SingletonPool().Log.Error("ListenAndServe", "err", err)
	}
}

// Stop 关闭 http 服务
func (my *HttpJobClient) Stop() {
	global.SingletonPool().Log.Info("embed http serve is closing")
	err := my.server.Shutdown(context.Background())
	if err != nil {
		global.SingletonPool().Log.Error("embed http serve shutdown error", "err", err)
	}
}

func commonHandlerFun[P any](res http.ResponseWriter, req *http.Request, fn func(param P) (any, error)) {
	defer func() {
		if err := recover(); err != nil {
			global.SingletonPool().Log.Error("http server error", "err", err)
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

	// 鉴权：如果未配置token，则认为不需要鉴权
	commonConfig := ezconfig.Get[*config.CommonConfigBean]()
	if commonConfig.AccessToken != "" {
		accessToken := req.Header.Get(sdk.TokenHeaderKey)
		if commonConfig.AccessToken != accessToken {
			http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
	}

	var param P
	if err := json.NewDecoder(req.Body).Decode(&param); err != nil && !errors.Is(err, io.EOF) {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	paramJson, _ := json.Marshal(param)
	global.SingletonPool().Log.Info("http serve request",
		"url", req.URL.String(),
		"param", string(paramJson))

	data, err := fn(param)
	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "application/json")
	if err != nil {
		_ = json.NewEncoder(res).Encode(sdk.Fail(err.Error()))
	} else {
		_ = json.NewEncoder(res).Encode(sdk.Ok2[any](data))
	}

}
