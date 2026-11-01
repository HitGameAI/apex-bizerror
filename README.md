# Go BizError

> `go-bizerror` 是通用业务错误包：业务侧只传错误码字符串，包内负责错误码 → gRPC 状态码 / 网关错误码的映射、i18n 消息解析（按请求语言）、`{reason}` 模板注入

**三个模块，按需取用**

| 模块 | 职责 | 额外依赖 |
| --- | --- | --- |
| `go-bizerror` | 错误码映射、i18n 文案、gRPC 错误转换 | grpc + go-i18n（go-i18n 本身零依赖） |
| `go-bizerror/gwx` | 转换为 `go-rpc-gateway` 的 `AppError` | + go-rpc-gateway |
| `go-bizerror/pbx` | Model-PB 泛型转换，失败归并为业务错误 | + go-pbmo |

主模块原先直接引用网关的 `middleware` 与 `errors` 两个包，而 `middleware` 一路连到 `cpool`，再连到 ClickHouse 与 AWS S3 SDK——407 行代码拖进 **686 个外部包**
拆成三个模块后，只做 gRPC 错误转换的服务只需 **41 个**（grpc 自身闭包，产生 `status.Error` 必需）

[![稳定版本](https://img.shields.io/badge/版本-稳定-green.svg)](https://github.com/HitGameAI/apex-bizerror)
[![许可协议](https://img.shields.io/github/license/HitGameAI/apex-bizerror)](LICENSE)
[![最新版本](https://img.shields.io/github/v/release/HitGameAI/apex-bizerror)](https://github.com/HitGameAI/apex-bizerror/releases)
[![提交记录](https://img.shields.io/github/last-commit/HitGameAI/apex-bizerror)](https://github.com/HitGameAI/apex-bizerror/commits/master)
[![Go 版本](https://img.shields.io/github/go-mod/go-version/HitGameAI/apex-bizerror)](https://golang.org/)

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

全部通用错误码的 zh / en 文案已内置在包内，新服务**零 locales 配置**即有可读错误消息
消息解析优先级：

1. 服务自身 `locales/<lang>.json`（同名 key 覆盖内置文案，仅维护服务专属码即可）
2. 包内置通用文案（按请求语言，未收录语言回退 en）
3. 错误码本身（始终可读的兜底）

`ToGRPCError*` 与 `gwx.NewBizAppError*` 两条路径均已接入该回退链，内置文案里的占位符也会被 `templateData` 填充

## 安装

```bash
go get github.com/HitGameAI/apex-bizerror
```

本地联调时可以在业务项目 `go.mod` 里使用本地 replace：

```go
require github.com/HitGameAI/apex-bizerror v0.0.0

replace github.com/HitGameAI/apex-bizerror => ../go-bizerror
```

## 快速开始

服务启动期注册专属错误码映射（如 `bootstrap` 初始化时），业务代码里只需返回错误码：

```go
package bootstrap

import (
 "google.golang.org/grpc/codes"

 bizerror "github.com/HitGameAI/apex-bizerror"
 gwerrors "github.com/HitGameAI/apex-rpc-gateway/errors"
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

Model-PB 泛型转换，在 `pbx` 子模块（转换失败统一返回 `error.conversion_failed`，需先在 go-pbmo 注册转换器）：

```go
import "github.com/HitGameAI/apex-bizerror/pbx"

// 读路径：Model -> PB
item, err := pbx.ToPB[models.FooModel, pb.FooInfo](ctx, stored)
items, err := pbx.ToPBs[models.FooModel, pb.FooInfo](ctx, rows)

// 仓储查不到时直接返回 NotFound，省掉每个 Get handler 的 nil 判断
// 第三个参数留空则用通用的 error.not_found
item, err := pbx.ToPBOrNotFound[models.FooModel, pb.FooInfo](ctx, stored, "error.foo_not_found")

// 写路径：PB -> Model
model, err := pbx.FromPB[pb.CreateFooRequest, models.FooModel](ctx, req)
models, err := pbx.FromPBs[pb.FooInfo, models.FooModel](ctx, reqs)

// 复用目标对象，省一次分配；部分更新也用它：
// 先取出完整模型，再用请求 PB 覆盖，未映射到的字段保持原样
err := pbx.FromPBInto[pb.UpdateFooRequest, models.FooModel](ctx, req, stored)
```

`pbx` 只封装「会返回 error 且该 error 应当变成业务错误」的那些：`SafeToPBs` / `SafeFromPBs` 用 `BatchResult` 表达部分成功，映射不到单个 error；`ModelToUpdates` / `PBToUpdates` 不返回 error，两类都请直接调 go-pbmo

i18n 模板变量（`{reason}` 由 WithCause 自动注入，其他占位符手工传入）：

```go
return bizerror.ToGRPCErrorWithTemplate(ctx, "error.quota_exceeded", map[string]any{
 "limit": 100,
})
```

非 gRPC 场景转换为网关 `AppError`，在 `gwx` 子模块：

```go
import "github.com/HitGameAI/apex-bizerror/gwx"

appErr := gwx.NewBizAppError(ctx, "error.conflict")
appErr := gwx.NewBizAppErrorWithCause(ctx, "error.database_error", err) // 注入 {reason}
```

跨服务调用侧提取可读错误消息（兼容 `fmt.Errorf("%w")` 包裹场景）：

```go
msg := bizerror.ExtractRpcErrorMsg(err)
```

对应关系：`ToGRPCError*` 系列用于 gRPC 服务实现内返回；`gwx.NewBizAppError*` 系列用于网关/HTTP 处理器内返回；两者共享同一份映射注册表

## 迁移对照

| 旧 | 新 |
| --- | --- |
| `bizerror.NewBizAppError*` | `gwx.NewBizAppError*` |
| `bizerror.ToPB` / `ToPBs` | `pbx.ToPB` / `pbx.ToPBs` |
| `bizerror.MapBizCodeToGatewayErrorCode` 返回 `gwerrors.ErrorCode` | `bizerror.MapBizCodeToGatewayCode` 返回 `bizerror.GatewayCode`（数值与网关一致，`gwx` 内有测试逐一比对） |
| `CodeMapping.Gateway` 类型 `gwerrors.ErrorCode` | `bizerror.GatewayCode` |
