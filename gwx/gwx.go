/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-08-16 15:28:03
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-08-16 16:05:12
 * @FilePath: \go-bizerror\gwx\gwx.go
 * @Description: 业务错误码到网关 AppError 的适配（唯一引入 go-rpc-gateway 的地方）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

// Package gwx 把业务错误码转换为 go-rpc-gateway 的 AppError
//
// 它独立成模块，是为了让 go-bizerror 主模块不引入 go-rpc-gateway：
// 主模块原先直接引用网关的 middleware 与 errors 两个包，而 middleware 一路连到
// cpool，再连到 ClickHouse 与 AWS S3 SDK——407 行代码拖进 686 个外部包
// 拆开后，只做 gRPC 错误转换的服务不必再承担网关那一整套
//
// [EN] A separate module so go-bizerror's main module pulls in no gateway code:
// it used to reference the gateway's middleware and errors packages directly,
// and middleware reaches cpool and from there ClickHouse and the AWS S3 SDK,
// so 407 lines dragged in 686 external packages
package gwx

import (
	"context"

	bizerror "github.com/kamalyes/go-bizerror"
	i18n "github.com/kamalyes/go-i18n"
	gwerrors "github.com/kamalyes/go-rpc-gateway/errors"
)

// NewBizAppError 将业务错误码转换为网关 AppError
func NewBizAppError(ctx context.Context, bizerrCode string) *gwerrors.AppError {
	return NewBizAppErrorWithTemplate(ctx, bizerrCode, nil)
}

// NewBizAppErrorWithTemplate 使用模板数据将业务错误码转换为网关 AppError
// 网关 locales 未命中文案时回退 go-bizerror 的内置通用文案，避免空 details
//
// 这里直接用 gwerrors.NewError 加 go-i18n 查文案，没有走网关的
// middleware.NewLocalizedAppErrorWithMap：那个函数本身只是这两步的组合，
// 却会把整个 middleware 包连带进来
// [EN] Builds the AppError from gwerrors.NewError plus a go-i18n lookup instead of
// the gateway's middleware.NewLocalizedAppErrorWithMap, which is exactly those two
// steps but drags the whole middleware package in
func NewBizAppErrorWithTemplate(ctx context.Context, bizerrCode string, templateData map[string]any) *gwerrors.AppError {
	code := gwerrors.ErrorCode(bizerror.MapBizCodeToGatewayCode(bizerrCode))
	appErr := gwerrors.NewError(code, "")

	if ctx == nil || bizerrCode == "" {
		return appErr
	}

	var localized string
	if templateData == nil {
		localized = i18n.GetMsgByKey(ctx, bizerrCode)
	} else {
		localized = i18n.GetMsgWithMap(ctx, bizerrCode, templateData)
	}
	if localized != "" && localized != bizerrCode {
		return appErr.WithDetails(localized)
	}

	// 服务 locales 未命中，回退内置通用文案，保证 details 不为空
	//
	// 内置表里存的是带占位符的模板（如 "Invalid request parameter: {reason}"），
	// 必须把 templateData 套上再写入——否则 locales 未覆盖该错误码时，调用方会
	// 收到字面量 {reason}，而内置文案存在的理由恰恰是让新服务零配置也有可读消息
	// [EN] The builtin table stores templates with placeholders, so templateData
	// must be applied before use; otherwise the caller gets a literal {reason}
	// exactly in the zero-config case the table exists for
	if msg := bizerror.BuiltinMessage(i18n.GetLanguage(ctx), bizerrCode); msg != "" {
		return appErr.WithDetails(i18n.FormatWithTemplateData(msg, templateData))
	}
	return appErr
}

// NewBizAppErrorWithCause 将业务错误码转换为网关 AppError，并把底层 error 文案
// 注入模板占位符 {reason}/{{reason}}/{{.reason}}，便于在 i18n 文案中追加具体原因
// cause 为 nil 时等同于 NewBizAppError
func NewBizAppErrorWithCause(ctx context.Context, bizerrCode string, cause error) *gwerrors.AppError {
	if cause == nil {
		return NewBizAppError(ctx, bizerrCode)
	}
	return NewBizAppErrorWithTemplate(ctx, bizerrCode, map[string]any{
		"reason": cause.Error(),
	})
}
