# hecc-blot-log-sls

阿里云日志服务（SLS）日志后端：实现 core 的 `ILog` 契约，与本地日志二选一接入。

## 安装

```bash
go get github.com/hecc-blot/log-sls
```

## 说明

本模块依赖 [core](https://github.com/hecc-blot/core) 的日志契约（`core/contract/log`），实现 `ILog` 接口将日志写入阿里云 SLS；SLS 配置类型（`SlsConfig`）定义在本模块的 `config` 包。

## 接口契约

实现的是 core 定义的 `ILog`：

```go
// github.com/hecc-blot/core/contract/log
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
    logContract "github.com/hecc-blot/core/contract/log"
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
    endpoint: "cn-hangzhou.log.aliyuncs.com"
    access_key: "your-access-key"
    secret_key: "your-secret"
    secret_token: ""           # 阿里云 STS Token（可选）
    project: "your-project"
    log_store: "your-logstore"
```

| 配置项 | 类型 | 说明 |
|--------|------|------|
| `endpoint` | string | SLS 端点地址 |
| `access_key` | string | 阿里云 AccessKey |
| `secret_key` | string | 阿里云 SecretKey |
| `secret_token` | string | 阿里云 STS Token（可选） |
| `project` | string | SLS Project 名称 |
| `log_store` | string | SLS LogStore 名称 |

## 与本地日志的关系

本地日志（core `service/log`，Zap + lumberjack 文件滚动）与 SLS 日志**二选一**，业务方按需显式指定构造哪一个，不做 `enable` 自动切换：

```go
logSvc, err := log.NewLogger(&config.Log.Local)    // 本地日志
// 或
logSvc, err := logsls.NewLogger(&config.Log.Sls)   // SLS 日志
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
| [core](https://github.com/hecc-blot/core) | `ILog` 契约、本地日志 |
| [trace](https://github.com/hecc-blot/trace) | TraceId 与日志关联 |
