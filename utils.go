/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-08-15 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-08-15 01:16:16
 * @FilePath: \go-bizerror\utils.go
 * @Description: 错误信息提取工具
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package bizerror

import (
	"errors"
	"strings"

	"google.golang.org/grpc/status"
)

// ExtractRpcErrorMsg 从 RPC 错误中提取可读消息（无 status 时回退 error 文案）
func ExtractRpcErrorMsg(err error) string {
	if err == nil {
		return ""
	}

	if rpcErr, ok := status.FromError(err); ok {
		return rpcErr.Message()
	}

	// 兼容包裹了 RPC 错误的 fmt.Errorf("%w") 场景
	var wrapper interface{ Unwrap() error }
	if errors.As(err, &wrapper) && wrapper.Unwrap() != nil {
		if rpcErr, ok := status.FromError(wrapper.Unwrap()); ok {
			return rpcErr.Message()
		}
	}

	return strings.TrimSpace(err.Error())
}
