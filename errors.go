/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-08-15 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-08-15 01:16:16
 * @FilePath: \go-bizerror\errors.go
 * @Description: 国际化业务错误（Code + Message，消息经 go-rpc-gateway 中间件按请求语言解析）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package bizerror

import (
	"context"
	"errors"

	gwerrors "github.com/kamalyes/go-rpc-gateway/errors"
	"github.com/kamalyes/go-rpc-gateway/middleware"
)

// Error 标准业务错误（i18n 解析后的 Code + Message）
type Error struct {
	Code    string
	Message string
}

// New 创建普通错误（等价标准库 errors.New）
func New(message string) error {
	return errors.New(message)
}

// NewI18nError 按错误码解析当前请求语言的国际化消息；
// 解析优先级：服务 locales（可覆盖内置文案）> 包内置通用文案 > 错误码本身，保证 Code 始终可读
func NewI18nError(ctx context.Context, errorCode string) *Error {
	message := middleware.GetMsgByKey(ctx, errorCode)

	if message == "" || message == errorCode {
		message = builtinMessage(middleware.GetLanguage(ctx), errorCode)
	}
	if message == "" {
		message = errorCode
	}

	return &Error{
		Code:    errorCode,
		Message: message,
	}
}

// NewI18nErrorWithTemplate 按错误码 + 模板变量解析国际化消息（如 {reason} 占位）
func NewI18nErrorWithTemplate(ctx context.Context, errorCode string, templateData map[string]interface{}) *Error {
	message := middleware.GetMsgWithMap(ctx, errorCode, templateData)

	if message == "" || message == errorCode {
		message = builtinMessage(middleware.GetLanguage(ctx), errorCode)
	}
	if message == "" {
		message = errorCode
	}

	return &Error{
		Code:    errorCode,
		Message: message,
	}
}

// NewBizAppError 将业务错误码转换为网关 AppError
func NewBizAppError(ctx context.Context, bizerrCode string) *gwerrors.AppError {
	return NewBizAppErrorWithTemplate(ctx, bizerrCode, nil)
}

// NewBizAppErrorWithTemplate 使用模板数据将业务错误码转换为网关 AppError
// 网关 locales 未命中文案时回退包内置通用文案，避免空 details
func NewBizAppErrorWithTemplate(ctx context.Context, bizerrCode string, templateData map[string]any) *gwerrors.AppError {
	appErr := middleware.NewLocalizedAppErrorWithMap(ctx, MapBizCodeToGatewayErrorCode(bizerrCode), bizerrCode, templateData)
	if appErr != nil && appErr.GetDetails() == "" {
		if msg := builtinMessage(middleware.GetLanguage(ctx), bizerrCode); msg != "" {
			appErr.WithDetails(msg)
		}
	}
	return appErr
}

// NewBizAppErrorWithCause 将业务错误码转换为网关 AppError，并把底层 error 文案
// 注入到模板占位符 {reason}/{{reason}}/{{.reason}}，便于在 i18n 文案中追加具体原因
// 若 cause 为 nil，行为等同于 NewBizAppError
func NewBizAppErrorWithCause(ctx context.Context, bizerrCode string, cause error) *gwerrors.AppError {
	if cause == nil {
		return NewBizAppError(ctx, bizerrCode)
	}
	return NewBizAppErrorWithTemplate(ctx, bizerrCode, map[string]any{
		"reason": cause.Error(),
	})
}
