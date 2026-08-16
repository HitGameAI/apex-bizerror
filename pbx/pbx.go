/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-08-16 15:28:03
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-08-16 16:05:12
 * @FilePath: \go-bizerror\pbx\pbx.go
 * @Description: Model-PB 泛型转换便捷封装（唯一引入 go-pbmo 的地方）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

// Package pbx 把 go-pbmo 的 Model-PB 转换与 go-bizerror 的错误码串起来：
// 转换失败统一返回 error.conversion_failed 对应的 gRPC 错误，调用方不必每处都写
// if err != nil { return nil, bizerror.ToGRPCError(ctx, ...) }
//
// 只封装「会返回 error、且该 error 应当变成业务错误」的那些：
//   - SafeToPBs / SafeFromPBs 用 BatchResult 表达部分成功，映射不到单个 error，
//     需要就直接调 go-pbmo
//   - ModelToUpdates / PBToUpdates 不返回 error，没有可转换的东西
//
// 它独立成模块，是为了让 go-bizerror 主模块不引入 go-pbmo 与 protobuf
// 主模块的职责是业务错误码到 gRPC 状态码与 i18n 文案的映射，Model-PB 转换
// 只是顺手提供的便利，不该让所有使用者都背上 50 多个 protobuf 相关包
//
// [EN] Wires go-pbmo's model/PB conversion to go-bizerror's error codes so
// callers stop repeating the same error-mapping block. Only the functions that
// return an error worth turning into a business error are wrapped: the Safe*
// variants express partial success through BatchResult, and the Updates builders
// return no error at all. A separate module so go-bizerror's main module pulls in
// neither go-pbmo nor protobuf
package pbx

import (
	"context"

	bizerror "github.com/kamalyes/go-bizerror"
	"github.com/kamalyes/go-pbmo"
)

// ==================== Model -> PB（读路径） ====================

// ToPB 单个模型转 PB，转换失败返回 error.conversion_failed 的 gRPC 错误
func ToPB[M any, PB any](ctx context.Context, m *M) (*PB, error) {
	result, err := pbmo.ToPB[M, PB](m)
	if err != nil {
		return nil, conversionFailed(ctx)
	}
	return result, nil
}

// ToPBs 模型列表转 PB 列表
func ToPBs[M any, PB any](ctx context.Context, items []*M) ([]*PB, error) {
	result, err := pbmo.ToPBs[M, PB](items)
	if err != nil {
		return nil, conversionFailed(ctx)
	}
	return result, nil
}

// ToPBOrNotFound 单个模型转 PB，模型为 nil 时返回 notFoundCode 对应的 gRPC 错误
//
// 仓储层查不到记录通常返回 (nil, nil)，而 gRPC 接口必须回 NotFound 而不是空对象，
// 于是每个 Get 类 handler 都要手写一次 nil 判断；这里收口掉
// notFoundCode 留给调用方指定，便于给出「用户不存在」这类更具体的文案；
// 传空串则用通用的 error.not_found
//
// [EN] A repository usually returns (nil, nil) when a row is missing, while a
// gRPC endpoint must answer NotFound rather than an empty object, so every Get
// handler repeats the same nil check. notFoundCode is the caller's to choose so
// the message can be specific; an empty string falls back to error.not_found
func ToPBOrNotFound[M any, PB any](ctx context.Context, m *M, notFoundCode string) (*PB, error) {
	if m == nil {
		if notFoundCode == "" {
			notFoundCode = bizerror.BizErrCodeNotFound
		}
		return nil, bizerror.ToGRPCError(ctx, notFoundCode)
	}
	return ToPB[M, PB](ctx, m)
}

// ToPBInto 把模型转换写入已有的 PB 对象，省掉一次分配
// 适用于循环内复用同一个目标对象的热路径
// [EN] Converts into an existing PB, saving an allocation; for hot paths that
// reuse one destination across a loop
func ToPBInto[M any, PB any](ctx context.Context, m *M, dst *PB) error {
	if err := pbmo.ToPBInto[M, PB](m, dst); err != nil {
		return conversionFailed(ctx)
	}
	return nil
}

// ==================== PB -> Model（写路径） ====================

// FromPB 单个 PB 转模型，转换失败返回 error.conversion_failed 的 gRPC 错误
// 写接口把请求 PB 落成库模型时用
// [EN] For write endpoints turning a request PB into a persistence model
func FromPB[PB any, M any](ctx context.Context, pb *PB) (*M, error) {
	result, err := pbmo.FromPB[PB, M](pb)
	if err != nil {
		return nil, conversionFailed(ctx)
	}
	return result, nil
}

// FromPBs PB 列表转模型列表
func FromPBs[PB any, M any](ctx context.Context, pbs []*PB) ([]*M, error) {
	result, err := pbmo.FromPBs[PB, M](pbs)
	if err != nil {
		return nil, conversionFailed(ctx)
	}
	return result, nil
}

// FromPBInto 把 PB 转换写入已有的模型对象
//
// 除了省分配，更常用的场景是部分更新：先从库里取出完整模型，再用请求 PB 覆盖
// 其中有值的字段，未映射到的字段保持原样
// [EN] Besides saving an allocation, the common use is a partial update: load the
// full model, then overlay the fields the request PB carries, leaving unmapped
// fields untouched
func FromPBInto[PB any, M any](ctx context.Context, pb *PB, dst *M) error {
	if err := pbmo.FromPBInto[PB, M](pb, dst); err != nil {
		return conversionFailed(ctx)
	}
	return nil
}

// conversionFailed 统一的转换失败错误
//
// 刻意不透出 go-pbmo 的原始错误：它描述的是字段级反射细节，对调用方无意义，
// 而暴露内部结构名也不合适；需要排查时看 go-pbmo 自己的日志
// [EN] The underlying go-pbmo error is deliberately not surfaced: it describes
// field-level reflection details that mean nothing to the caller and would leak
// internal type names
func conversionFailed(ctx context.Context) error {
	return bizerror.ToGRPCError(ctx, bizerror.BizErrCodeConversionFailed)
}
