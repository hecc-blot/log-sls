# hecc-blot-log-sls

阿里云日志服务（SLS）日志后端：实现 framework 的 `ILog` 契约，与本地日志二选一接入。

## 安装

```bash
go get github.com/hecc-blot/log-sls
```

## 说明

本模块依赖 [framework](https://github.com/hecc-blot/framework) 的日志契约（`framework/contract/log`）与配置类型（`framework/config/log`），实现 `ILog` 接口将日志写入阿里云 SLS。

## 接口契约

实现的是 framework 定义的 `ILog`：

```go
// github.com/hecc-blot/framework/contract/log
type ILog interface {
    Debug(ctx context.Context, msg string, fields ...interface{})
    Info(ctx context.Context, msg string, fields ...interface{})
    Warn(ctx context.Context, msg string, fields ...interface{})
    Error(ctx context.Context, msg string, fields ...interface{})
}
```

## 初始化

```go
import (
    logContract "github.com/hecc-blot/framework/contract/log"
    logsls "github.com/hecc-blot/log-sls/service"
)

logSvc, err := logsls.NewLogger(&config.Log.Sls)
if err != nil {
    panic(err)
}

container.Set(new(logContract.ILog), logSvc)
```

## 配置

```yaml
log:
  sls:
    enable: true
    endpoint: "cn-hangzhou.log.aliyuncs.com"
    access_key: "your-access-key"
    secret_key: "your-secret"
    secret_token: ""           # 阿里云 STS Token（可选）
    project: "your-project"
    log_store: "your-logstore"
```

| 配置项 | 类型 | 说明 |
|--------|------|------|
| `enable` | bool | 是否启用 |
| `endpoint` | string | SLS 端点地址 |
| `access_key` | string | 阿里云 AccessKey |
| `secret_key` | string | 阿里云 SecretKey |
| `secret_token` | string | 阿里云 STS Token（可选） |
| `project` | string | SLS Project 名称 |
| `log_store` | string | SLS LogStore 名称 |

## 与本地日志的关系

本地日志（framework `service/log`，Zap + lumberjack 文件滚动）与 SLS 日志**二选一**，通过 `enable` 配置启用。组装层按需选择：

```go
var logSvc logContract.ILog
if config.Log.Sls.Enable {
    logSvc = must(logsls.NewLogger(&config.Log.Sls))   // 优先 SLS
} else {
    logSvc = must(log.NewLogger(&config.Log.Local))    // 本地日志
}
```

## 使用

与本地日志一致，通过 IOC 注入 `ILog` 后调用：

```go
type AddApi struct {
    LogSvc    logContract.ILog  `inject:""`
    // ...
}

func (a AddApi) Call(ctx *gin.Context) (interface{}, iCoreError.IError) {
    a.LogSvc.Info(ctx, "add account request", "name", a.AccountName)
    // ...
}
```

## 相关模块

| 模块 | 说明 |
|------|------|
| [framework](https://github.com/hecc-blot/framework) | `ILog` 契约、`SlsConfig` 配置类型、本地日志 |
| [trace](https://github.com/hecc-blot/trace) | TraceId 与日志关联 |
