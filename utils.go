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
//
// 用 errors.As 直接找链上最内层实现 GRPCStatus 的那个错误，而不是 status.FromError
//
// 原因：grpc-go 的 status.FromError 在遇到被包裹的 status 错误时，虽然能 unwrap
// 拿到内层 status，却会把外层 error 的完整文案覆盖进 Message
// 于是 fmt.Errorf("调用下游失败: %w", status.Error(codes.NotFound, "缺少用户"))
// 提取出来是 "调用下游失败: rpc error: code = NotFound desc = 缺少用户"——
// 把 gRPC 的信封文本漏给了调用方，而这个函数的用途恰恰是拿到干净的 "缺少用户"
// 旧实现因此有两处问题：第一个分支先命中并返回带信封的串，使后面那段专门用来
// 兼容 %w 包裹的代码成为死代码
//
// [EN] Uses errors.As to reach the innermost error implementing GRPCStatus rather
// than status.FromError: when grpc-go unwraps a wrapped status error it overwrites
// Message with the outer error's full text, so the extracted message leaks the
// "rpc error: code = ..." envelope — exactly what this helper exists to strip
// The previous implementation's first branch matched first and returned the
// enveloped string, leaving its %w-compatibility branch unreachable
func ExtractRpcErrorMsg(err error) string {
	if err == nil {
		return ""
	}

	var withStatus interface{ GRPCStatus() *status.Status }
	if errors.As(err, &withStatus) {
		if st := withStatus.GRPCStatus(); st != nil {
			return st.Message()
		}
	}

	return strings.TrimSpace(err.Error())
}
