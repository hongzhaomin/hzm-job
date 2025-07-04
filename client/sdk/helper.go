package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hongzhaomin/hzm-job/client/model"
	"log/slog"
)

func Post[Param any, Res any](ctx context.Context, client *RemotingUtil, url, accessToken string, body *Param) *Res {
	jsonStr, err := client.PostJSON(ctx, url, accessToken, body)
	if err != nil {
		slog.Error(fmt.Sprintf("请求http异常，url: %s", url), err.Error())
		return nil
	}

	var res model.Result[Res]
	_ = json.Unmarshal(jsonStr, &res)
	if res.Success {
		return res.Data
	}
	return nil
}
