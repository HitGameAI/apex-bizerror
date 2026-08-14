# Go BizError

> `go-bizerror` 是 `go-rpc-gateway` 服务的通用业务错误包，业务侧只需要传入错误码字符串，包内会负责错误码 → gRPC 状态码 / 网关 AppError 错误码的映射、i18n 消息解析（按请求语言）、`{reason}` 模板注入，以及 Model-PB 泛型转换失败的统一归并

[![稳定版本](https://img.shields.io/badge/版本-稳定-green.svg)](https://github.com/kamalyes/go-bizerror)
[![许可协议](https://img.shields.io/github/license/kamalyes/go-bizerror)](LICENSE)
[![最新版本](https://img.shields.io/github/v/release/kamalyes/go-bizerror)](https://github.com/kamalyes/go-bizerror/releases)
[![提交记录](https://img.shields.io/github/last-commit/kamalyes/go-bizerror)](https://github.com/kamalyes/go-bizerror/commits/master)
[![Go 版本](https://img.shields.io/github/go-mod/go-version/kamalyes/go-bizerror)](https://golang.org/)

## 错误码约定

错误码为 `error.xxx` 形式的字符串常量，通用码已内置映射：

| 分类 | 错误码 | gRPC | 网关 |
| --- | --- | --- | --- |
| 通用 | `error.invalid_argument` | InvalidArgument | BadRequest |
| 通用 | `error.validation_failed` | InvalidArgument | BadRequest |
| 通用 | `error.invalid_parameter` | InvalidArgument | BadRequest |
| 通用 | `error.unauthenticated` | Unauthenticated | Unauthorized |
| 通用 | `error.permission_denied` | PermissionDenied | Forbidden |
| 通用 | `error.not_found` | NotFound | NotFound |
| 通用 | `error.conflict` | AlreadyExists | Conflict |
| 通用 | `error.cannot_be_self` | FailedPrecondition | BadRequest |
| 内部 | `error.internal` | Internal | InternalServerError |
| 内部 | `error.database_error` | Internal | InternalServerError |
| 内部 | `error.database_operation_failed` | Internal | InternalServerError |
| 内部 | `error.conversion_failed` | Internal | InternalServerError |
| 内部 | `error.unsupported_operation` | Unimplemented | InternalServerError |
| 内部 | `error.external_service_error` | Internal | InternalServerError |

映射解析顺序：显式注册 > `*_not_found` 后缀规则（自动映射为 NotFound）> Internal 兜底

## 内置 i18n 文案

全部通用错误码的 zh / en 文案已内置在包内，新服务**零 locales 配置**即有可读错误消息。消息解析优先级：

1. 服务自身 `locales/<lang>.json`（同名 key 覆盖内置文案，仅维护服务专属码即可）
2. 包内置通用文案（按请求语言，未收录语言回退 en）
3. 错误码本身（始终可读的兜底）

`ToGRPCError*` 与 `NewBizAppError*` 两条路径均已接入该回退链

## 安装

```bash
go get github.com/kamalyes/go-bizerror
```

本地联调时可以在业务项目 `go.mod` 里使用本地 replace：

```go
require github.com/kamalyes/go-bizerror v0.0.0

replace github.com/kamalyes/go-bizerror => ../go-bizerror
```

## 快速开始

服务启动期注册专属错误码映射（如 `bootstrap` 初始化时），业务代码里只需返回错误码：

```go
package bootstrap

import (
 "google.golang.org/grpc/codes"

 bizerror "github.com/kamalyes/go-bizerror"
 gwerrors "github.com/kamalyes/go-rpc-gateway/errors"
)

func registerBizErrors() {
 bizerror.RegisterMany(map[string]bizerror.CodeMapping{
  "error.dashboard_time_range_invalid": {GRPC: codes.InvalidArgument, Gateway: gwerrors.ErrCodeBadRequest},
  "error.roi_consumption_duplicate":    {GRPC: codes.AlreadyExists, Gateway: gwerrors.ErrCodeConflict},
 })
}
```

业务方法中直接转换返回：

```go
func (s *ServiceImpl) GetFoo(ctx context.Context, req *pb.GetFooRequest) (*pb.GetFooResponse, error) {
 stored, err := s.repo.Find(ctx, req.GetId())
 if err != nil {
  return nil, bizerror.ToGRPCErrorWithCause(ctx, "error.database_operation_failed", err)
 }
 if stored == nil {
  return nil, bizerror.ToGRPCError(ctx, "error.foo_not_found") // _not_found 后缀自动映射 NotFound
 }
 return &pb.GetFooResponse{}, nil
}
```

## 常用写法

Model-PB 泛型转换（转换失败统一返回 `error.conversion_failed`，需先在 go-pbmo 注册转换器）：

```go
item, err := bizerror.ToPB[models.FooModel, pb.FooInfo](ctx, stored)
if err != nil {
 return nil, err
}

items, err := bizerror.ToPBs[models.FooModel, pb.FooInfo](ctx, rows)
```

i18n 模板变量（`{reason}` 由 WithCause 自动注入，其他占位符手工传入）：

```go
return bizerror.ToGRPCErrorWithTemplate(ctx, "error.quota_exceeded", map[string]any{
 "limit": 100,
})
```

非 gRPC 场景转换为网关 `AppError`：

```go
appErr := bizerror.NewBizAppError(ctx, "error.conflict")
```

跨服务调用侧提取可读错误消息（兼容 `fmt.Errorf("%w")` 包裹场景）：

```go
msg := bizerror.ExtractRpcErrorMsg(err)
```

对应关系：`ToGRPCError*` 系列用于 gRPC 服务实现内返回；`NewBizAppError*` 系列用于网关/HTTP 处理器内返回；两者共享同一份映射注册表
