/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-08-15 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-08-15 01:16:16
 * @FilePath: \go-bizerror\pbmo.go
 * @Description: Model-PB 泛型转换便捷封装（go-pbmo 按字段名自动映射，失败统一归为转换错误）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package bizerror

import (
	"context"

	"github.com/kamalyes/go-pbmo"
)

// ToPB 单个模型转 PB（转换失败统一返回 error.conversion_failed 的 gRPC 错误）
func ToPB[M any, PB any](ctx context.Context, m *M) (*PB, error) {
	result, err := pbmo.ToPB[M, PB](m)
	if err != nil {
		return nil, ToGRPCError(ctx, BizErrCodeConversionFailed)
	}
	return result, nil
}

// ToPBs 模型列表转 PB 列表
func ToPBs[M any, PB any](ctx context.Context, items []*M) ([]*PB, error) {
	result, err := pbmo.ToPBs[M, PB](items)
	if err != nil {
		return nil, ToGRPCError(ctx, BizErrCodeConversionFailed)
	}
	return result, nil
}
