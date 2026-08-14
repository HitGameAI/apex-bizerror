/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-08-15 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-08-15 00:00:00
 * @FilePath: \go-bizerror\grpc.go
 * @Description: 业务错误码 → gRPC status 错误转换
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package bizerror

import (
	"context"

	"google.golang.org/grpc/status"
)

// ToGRPCError 将业务错误码转换为携带 i18n 消息的 gRPC 错误
func ToGRPCError(ctx context.Context, bizErrorCode string) error {
	errCode := NewI18nError(ctx, bizErrorCode)
	grpcCode := MapBizCodeToGRPCCode(bizErrorCode)
	return status.Error(grpcCode, errCode.Message)
}

// ToGRPCErrorWithTemplate 将业务错误码 + 模板变量转换为 gRPC 错误
func ToGRPCErrorWithTemplate(ctx context.Context, bizErrorCode string, templateData map[string]interface{}) error {
	errCode := NewI18nErrorWithTemplate(ctx, bizErrorCode, templateData)
	grpcCode := MapBizCodeToGRPCCode(bizErrorCode)
	return status.Error(grpcCode, errCode.Message)
}

// ToGRPCErrorWithCause 将业务错误码转换为 gRPC 错误，并将底层错误信息注入 {reason} 模板占位；
// cause 为 nil 时退化为 ToGRPCError
func ToGRPCErrorWithCause(ctx context.Context, bizErrorCode string, cause error) error {
	if cause == nil {
		return ToGRPCError(ctx, bizErrorCode)
	}
	return ToGRPCErrorWithTemplate(ctx, bizErrorCode, map[string]interface{}{"reason": cause.Error()})
}
