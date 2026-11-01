/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-08-15 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-08-15 01:16:16
 * @FilePath: \go-bizerror\errors.go
 * @Description: 国际化业务错误（Code + Message，消息经 go-i18n 按请求语言解析）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package bizerror

import (
	"context"
	"errors"

	i18n "github.com/HitGameAI/apex-toolbox/pkg/i18n"
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
	message := i18n.GetMsgByKey(ctx, errorCode)

	if message == "" || message == errorCode {
		message = BuiltinMessage(i18n.GetLanguage(ctx), errorCode)
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
//
// 内置文案兜底时同样要套模板：内置表里存的是带占位符的模板（如
// "Invalid request parameter: {reason}"），若直接返回，服务 locales 未覆盖该
// 错误码时调用方就会收到字面量 {reason}——而内置文案存在的理由恰恰是让新服务
// 零配置也有可读消息
// [EN] The builtin fallback must have the template applied too: the table stores
// templates with placeholders, so returning one verbatim would hand the caller a
// literal {reason} precisely in the zero-config case the table exists for
func NewI18nErrorWithTemplate(ctx context.Context, errorCode string, templateData map[string]interface{}) *Error {
	message := i18n.GetMsgWithMap(ctx, errorCode, templateData)

	if message == "" || message == errorCode {
		if builtin := BuiltinMessage(i18n.GetLanguage(ctx), errorCode); builtin != "" {
			message = i18n.FormatWithTemplateData(builtin, templateData)
		} else {
			message = ""
		}
	}
	if message == "" {
		message = errorCode
	}

	return &Error{
		Code:    errorCode,
		Message: message,
	}
}
