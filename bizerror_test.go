/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-08-16 15:51:36
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-08-16 16:05:12
 * @FilePath: \go-bizerror\bizerror_test.go
 * @Description: 错误码映射、i18n 兜底与 gRPC 转换测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package bizerror

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMapping_Defaults(t *testing.T) {
	tests := []struct {
		bizCode string
		grpc    codes.Code
		gateway GatewayCode
	}{
		{BizErrCodeInternalError, codes.Internal, GatewayCodeInternalServerError},
		{BizErrCodeInvalidArgument, codes.InvalidArgument, GatewayCodeBadRequest},
		{BizErrCodeUnauthenticated, codes.Unauthenticated, GatewayCodeUnauthorized},
		{BizErrCodePermissionDenied, codes.PermissionDenied, GatewayCodeForbidden},
		{BizErrCodeNotFound, codes.NotFound, GatewayCodeNotFound},
		{BizErrCodeConflict, codes.AlreadyExists, GatewayCodeConflict},
		{BizErrCodeCannotBeSelf, codes.FailedPrecondition, GatewayCodeBadRequest},
		{BizErrCodeUnsupportedOperation, codes.Unimplemented, GatewayCodeInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.bizCode, func(t *testing.T) {
			m := Mapping(tt.bizCode)
			assert.Equal(t, tt.grpc, m.GRPC)
			assert.Equal(t, tt.gateway, m.Gateway)
			assert.Equal(t, tt.grpc, MapBizCodeToGRPCCode(tt.bizCode))
			assert.Equal(t, tt.gateway, MapBizCodeToGatewayCode(tt.bizCode))
		})
	}
}

// TestMapping_NotFoundSuffix 未显式注册的 *_not_found 走后缀规则而非兜底 Internal
func TestMapping_NotFoundSuffix(t *testing.T) {
	m := Mapping("user_not_found")
	assert.Equal(t, codes.NotFound, m.GRPC)
	assert.Equal(t, GatewayCodeNotFound, m.Gateway)
}

// TestMapping_UnknownFallsBackToInternal 既未注册也不匹配后缀规则时兜底 Internal
func TestMapping_UnknownFallsBackToInternal(t *testing.T) {
	m := Mapping("error.something_nobody_registered")
	assert.Equal(t, codes.Internal, m.GRPC)
	assert.Equal(t, GatewayCodeInternalServerError, m.Gateway)
}

func TestRegister(t *testing.T) {
	const code = "error.test_register_only"
	Register(code, CodeMapping{codes.ResourceExhausted, GatewayCodeConflict})

	m := Mapping(code)
	assert.Equal(t, codes.ResourceExhausted, m.GRPC)
	assert.Equal(t, GatewayCodeConflict, m.Gateway)

	// 后注册覆盖先注册
	Register(code, CodeMapping{codes.Aborted, GatewayCodeBadRequest})
	assert.Equal(t, codes.Aborted, Mapping(code).GRPC)

	// 空错误码被忽略，不会污染表
	Register("", CodeMapping{codes.Unknown, GatewayCodeNotFound})
	assert.Equal(t, codes.Internal, Mapping("").GRPC, "空 key 不应被注册，查询应走兜底")
}

func TestRegisterMany(t *testing.T) {
	RegisterMany(map[string]CodeMapping{
		"error.test_many_a": {codes.Canceled, GatewayCodeBadRequest},
		"error.test_many_b": {codes.DeadlineExceeded, GatewayCodeInternalServerError},
	})
	assert.Equal(t, codes.Canceled, Mapping("error.test_many_a").GRPC)
	assert.Equal(t, codes.DeadlineExceeded, Mapping("error.test_many_b").GRPC)
}

func TestBuiltinMessage(t *testing.T) {
	assert.Equal(t, "资源未找到", BuiltinMessage("zh", BizErrCodeNotFound))
	assert.Equal(t, "Resource not found", BuiltinMessage("en", BizErrCodeNotFound))
	// 语言为空时按 en
	assert.Equal(t, "Resource not found", BuiltinMessage("", BizErrCodeNotFound))
	// 未收录的语言回退 en
	assert.Equal(t, "Resource not found", BuiltinMessage("fr", BizErrCodeNotFound))
	// 未收录的错误码返回空串，由调用方决定兜底
	assert.Empty(t, BuiltinMessage("zh", "error.nobody_wrote_this"))
}

// TestNewI18nError_FallsBackToBuiltin ctx 里没有 i18n 上下文时回退内置文案，
// 这是「新服务零配置即有可读错误消息」的保证
func TestNewI18nError_FallsBackToBuiltin(t *testing.T) {
	e := NewI18nError(context.Background(), BizErrCodeNotFound)
	require.NotNil(t, e)
	assert.Equal(t, BizErrCodeNotFound, e.Code)
	assert.Equal(t, "Resource not found", e.Message)
}

// TestNewI18nError_UnknownCodeKeepsCodeReadable 内置文案也没有时回退错误码本身，
// Message 永不为空
func TestNewI18nError_UnknownCodeKeepsCodeReadable(t *testing.T) {
	const code = "error.totally_unknown"
	e := NewI18nError(context.Background(), code)
	assert.Equal(t, code, e.Code)
	assert.Equal(t, code, e.Message)
}

func TestNewI18nErrorWithTemplate(t *testing.T) {
	e := NewI18nErrorWithTemplate(context.Background(), BizErrCodeInvalidParameter,
		map[string]interface{}{"reason": "字段 age 必须为正数"})
	assert.Equal(t, BizErrCodeInvalidParameter, e.Code)
	// 内置 en 文案是 "Invalid request parameter: {reason}"，占位符应被填充
	assert.Contains(t, e.Message, "字段 age 必须为正数")
	assert.NotContains(t, e.Message, "{reason}")
}

func TestToGRPCError(t *testing.T) {
	err := ToGRPCError(context.Background(), BizErrCodeNotFound)
	require.Error(t, err)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Equal(t, "Resource not found", st.Message())
}

func TestToGRPCErrorWithCause(t *testing.T) {
	// cause 为 nil 时退化为 ToGRPCError
	plain := ToGRPCErrorWithCause(context.Background(), BizErrCodeNotFound, nil)
	st, _ := status.FromError(plain)
	assert.Equal(t, "Resource not found", st.Message())

	// cause 非 nil 时注入 {reason}
	withCause := ToGRPCErrorWithCause(context.Background(), BizErrCodeInvalidParameter,
		errors.New("连接被重置"))
	st2, _ := status.FromError(withCause)
	assert.Equal(t, codes.InvalidArgument, st2.Code())
	assert.Contains(t, st2.Message(), "连接被重置")
}

func TestExtractRpcErrorMsg(t *testing.T) {
	assert.Empty(t, ExtractRpcErrorMsg(nil))

	// gRPC status 错误取 Message
	assert.Equal(t, "boom", ExtractRpcErrorMsg(status.Error(codes.Internal, "boom")))

	// 被 fmt.Errorf("%w") 包裹的 status 错误仍能取到
	wrapped := fmt.Errorf("调用下游失败: %w", status.Error(codes.NotFound, "缺少用户"))
	assert.Equal(t, "缺少用户", ExtractRpcErrorMsg(wrapped))

	// 普通错误回退 Error() 并去除首尾空白
	assert.Equal(t, "普通错误", ExtractRpcErrorMsg(errors.New("  普通错误  ")))
}

func TestNew(t *testing.T) {
	err := New("自定义消息")
	require.Error(t, err)
	assert.Equal(t, "自定义消息", err.Error())
}

func TestGatewayCode_Int(t *testing.T) {
	assert.Equal(t, 1004, GatewayCodeInternalServerError.Int())
	assert.Equal(t, 3002, GatewayCodeNotFound.Int())
}
