package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hongzhaomin/hzm-job/core/internal/tools"
	"log/slog"
)

func Post[Res any](client *tools.RemotingUtil, url, accessToken string, body BaseParam[Res]) *Res {
	ctx := context.Background()
	jsonStr, err := client.PostJSON(ctx, url, accessToken, body)
	if err != nil {
		slog.Error(fmt.Sprintf("请求http异常，url: %s", url), err.Error())
		return nil
	}

	var res Result[Res]
	_ = json.Unmarshal(jsonStr, &res)
	if res.Success {
		return res.Data
	}
	return nil
}
