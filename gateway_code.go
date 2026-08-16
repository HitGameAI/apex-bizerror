/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-08-16 15:28:03
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-08-16 16:05:12
 * @FilePath: \go-bizerror\gateway_code.go
 * @Description: 网关错误码（本包自有，与 go-rpc-gateway/errors.ErrorCode 数值一致）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package bizerror

// GatewayCode 网关层错误码
//
// 为什么本包要自己定义一份
//
// 本包只有 407 行，职责是「业务错误码 -> gRPC 状态码 / 网关错误码 / i18n 文案」的
// 映射与转换，却因为直接引用 go-rpc-gateway 的类型而拖进 686 个外部包：
// middleware 那条边一路连到 cpool，再连到 ClickHouse 与 AWS S3 SDK
// 错误码定义属于底层契约，本该由下层持有、上层引用，方向原先是反的
//
// 这里的数值与 go-rpc-gateway/errors.ErrorCode 逐一对齐，两边都是 int 底层类型，
// 需要互转时一次显式转换即可（见 gwx 子包）
//
// [EN] This package is 407 lines whose job is mapping business error codes to
// gRPC status codes, gateway codes and i18n messages, yet referencing the
// gateway's types dragged in 686 external packages: the middleware edge reaches
// cpool and from there ClickHouse and the AWS S3 SDK. Error-code definitions are
// a lower-layer contract, so they belong here and the gateway references them,
// not the other way round. The values match errors.ErrorCode one for one and
// both have int as their underlying type, so a single explicit conversion
// bridges them (see the gwx subpackage).
type GatewayCode int

// 与 go-rpc-gateway/errors 的 ErrCode* 常量数值一一对应
// 改动这里的数值必须同步 go-rpc-gateway，gwx 子包有断言守住这一点
// [EN] One-to-one with the gateway's ErrCode* constants; the gwx subpackage
// asserts they stay in step.
const (
	// GatewayCodeInternalServerError 服务内部错误，对应 errors.ErrCodeInternalServerError
	GatewayCodeInternalServerError GatewayCode = 1004
	// GatewayCodeUnauthorized 未认证，对应 errors.ErrCodeUnauthorized
	GatewayCodeUnauthorized GatewayCode = 2001
	// GatewayCodeForbidden 权限不足，对应 errors.ErrCodeForbidden
	GatewayCodeForbidden GatewayCode = 2002
	// GatewayCodeBadRequest 请求错误，对应 errors.ErrCodeBadRequest
	GatewayCodeBadRequest GatewayCode = 3001
	// GatewayCodeNotFound 资源未找到，对应 errors.ErrCodeNotFound
	GatewayCodeNotFound GatewayCode = 3002
	// GatewayCodeConflict 资源冲突，对应 errors.ErrCodeConflict
	GatewayCodeConflict GatewayCode = 9004
)

// Int 返回底层整数值，便于转换为网关自己的错误码类型
func (c GatewayCode) Int() int { return int(c) }
