/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-08-16 15:51:36
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-08-16 16:05:12
 * @FilePath: \go-bizerror\gwx\gwx_test.go
 * @Description: 网关错误码数值一致性与 AppError 构造测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package gwx

import (
	"context"
	"errors"
	"testing"

	bizerror "github.com/HitGameAI/apex-bizerror"
	gwerrors "github.com/HitGameAI/apex-rpc-gateway/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGatewayCodeParity 守住 go-bizerror 手写的那组 GatewayCode 数值
//
// 主模块为了不依赖 go-rpc-gateway，把网关错误码的数值抄了一份
// 抄来的常量没有编译期约束，一旦网关侧改了数值，转换出来的 AppError 就会悄悄
// 带上错误的 code——HTTP 状态码跟着错，而且不报任何错
// 本测试是唯一能同时看见两边定义的地方，故在此逐一比对
//
// [EN] The main module copies the gateway's error-code values so it need not
// depend on the gateway. Copied constants carry no compile-time guarantee: if the
// gateway renumbers, the converted AppError silently gets the wrong code and the
// HTTP status follows it. This test is the only place that sees both definitions
func TestGatewayCodeParity(t *testing.T) {
	tests := []struct {
		name string
		mine bizerror.GatewayCode
		gw   gwerrors.ErrorCode
	}{
		{"InternalServerError", bizerror.GatewayCodeInternalServerError, gwerrors.ErrCodeInternalServerError},
		{"Unauthorized", bizerror.GatewayCodeUnauthorized, gwerrors.ErrCodeUnauthorized},
		{"Forbidden", bizerror.GatewayCodeForbidden, gwerrors.ErrCodeForbidden},
		{"BadRequest", bizerror.GatewayCodeBadRequest, gwerrors.ErrCodeBadRequest},
		{"NotFound", bizerror.GatewayCodeNotFound, gwerrors.ErrCodeNotFound},
		{"Conflict", bizerror.GatewayCodeConflict, gwerrors.ErrCodeConflict},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, int(tt.gw), tt.mine.Int(),
				"go-bizerror 的 GatewayCode%s 与网关 ErrCode%s 数值不符，改动任一侧都要同步另一侧", tt.name, tt.name)
		})
	}
}

func TestNewBizAppError(t *testing.T) {
	appErr := NewBizAppError(context.Background(), bizerror.BizErrCodeNotFound)
	require.NotNil(t, appErr)
	assert.Equal(t, gwerrors.ErrCodeNotFound, appErr.Code)
	// 无 i18n 上下文时回退内置文案，details 不应为空
	assert.Equal(t, "Resource not found", appErr.GetDetails())
}

func TestNewBizAppError_UnknownCodeFallsBackToInternal(t *testing.T) {
	appErr := NewBizAppError(context.Background(), "error.nobody_registered_this")
	require.NotNil(t, appErr)
	assert.Equal(t, gwerrors.ErrCodeInternalServerError, appErr.Code)
}

func TestNewBizAppErrorWithTemplate(t *testing.T) {
	appErr := NewBizAppErrorWithTemplate(context.Background(),
		bizerror.BizErrCodeInvalidParameter, map[string]any{"reason": "age 必须为正数"})
	require.NotNil(t, appErr)
	assert.Equal(t, gwerrors.ErrCodeBadRequest, appErr.Code)
	assert.Contains(t, appErr.GetDetails(), "age 必须为正数")
	assert.NotContains(t, appErr.GetDetails(), "{reason}", "占位符必须被填充")
}

func TestNewBizAppErrorWithCause(t *testing.T) {
	// cause 为 nil 时等同 NewBizAppError
	plain := NewBizAppErrorWithCause(context.Background(), bizerror.BizErrCodeNotFound, nil)
	assert.Equal(t, "Resource not found", plain.GetDetails())

	// cause 非 nil 时注入 {reason}
	withCause := NewBizAppErrorWithCause(context.Background(),
		bizerror.BizErrCodeInvalidParameter, errors.New("字段缺失"))
	assert.Equal(t, gwerrors.ErrCodeBadRequest, withCause.Code)
	assert.Contains(t, withCause.GetDetails(), "字段缺失")
}

// TestNewBizAppErrorWithTemplate_NilCtx ctx 为 nil 时不应 panic，返回只带 code 的 AppError
func TestNewBizAppErrorWithTemplate_NilCtx(t *testing.T) {
	appErr := NewBizAppErrorWithTemplate(nil, bizerror.BizErrCodeNotFound, nil)
	require.NotNil(t, appErr)
	assert.Equal(t, gwerrors.ErrCodeNotFound, appErr.Code)
}
