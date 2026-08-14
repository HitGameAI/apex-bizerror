/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-08-15 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-08-15 01:16:16
 * @FilePath: \go-bizerror\registry.go
 * @Description: 业务错误码 → gRPC/网关错误码映射注册中心（通用码内置默认映射，服务专属码启动期注入）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package bizerror

import (
	"strings"
	"sync"

	gwerrors "github.com/kamalyes/go-rpc-gateway/errors"
	"google.golang.org/grpc/codes"
)

// CodeMapping 业务错误码映射（gRPC 状态码 + 网关 AppError 错误码）
type CodeMapping struct {
	GRPC    codes.Code
	Gateway gwerrors.ErrorCode
}

var (
	mappingMu sync.RWMutex
	mappings  = map[string]CodeMapping{}
)

func init() {
	// 通用错误码默认映射（与各服务原有 switch 分支保持一致）
	RegisterMany(map[string]CodeMapping{
		BizErrCodeInternalError:           {codes.Internal, gwerrors.ErrCodeInternalServerError},
		BizErrCodeInvalidArgument:         {codes.InvalidArgument, gwerrors.ErrCodeBadRequest},
		BizErrCodeValidationFailed:        {codes.InvalidArgument, gwerrors.ErrCodeBadRequest},
		BizErrCodeInvalidParameter:        {codes.InvalidArgument, gwerrors.ErrCodeBadRequest},
		BizErrCodeUnauthenticated:         {codes.Unauthenticated, gwerrors.ErrCodeUnauthorized},
		BizErrCodePermissionDenied:        {codes.PermissionDenied, gwerrors.ErrCodeForbidden},
		BizErrCodeNotFound:                {codes.NotFound, gwerrors.ErrCodeNotFound},
		BizErrCodeConflict:                {codes.AlreadyExists, gwerrors.ErrCodeConflict},
		BizErrCodeCannotBeSelf:            {codes.FailedPrecondition, gwerrors.ErrCodeBadRequest},
		BizErrCodeUnsupportedOperation:    {codes.Unimplemented, gwerrors.ErrCodeInternalServerError},
		BizErrCodeDatabaseError:           {codes.Internal, gwerrors.ErrCodeInternalServerError},
		BizErrCodeDatabaseOperationFailed: {codes.Internal, gwerrors.ErrCodeInternalServerError},
		BizErrCodeConversionFailed:        {codes.Internal, gwerrors.ErrCodeInternalServerError},
		BizErrCodeExternalServiceError:    {codes.Internal, gwerrors.ErrCodeInternalServerError},
	})
}

// Register 注册单个业务错误码映射（幂等，后注册覆盖先注册；空错误码忽略）
func Register(bizCode string, mapping CodeMapping) {
	if bizCode == "" {
		return
	}
	mappingMu.Lock()
	defer mappingMu.Unlock()
	mappings[bizCode] = mapping
}

// RegisterMany 批量注册业务错误码映射（服务专属码在启动期集中注入）
func RegisterMany(batch map[string]CodeMapping) {
	for bizCode, mapping := range batch {
		Register(bizCode, mapping)
	}
}

// Mapping 查询错误码映射：显式注册优先，其次 _not_found 后缀规则，兜底 Internal
func Mapping(bizCode string) CodeMapping {
	mappingMu.RLock()
	mapping, ok := mappings[bizCode]
	mappingMu.RUnlock()
	if ok {
		return mapping
	}

	if bizCode == BizErrCodeNotFound || strings.HasSuffix(bizCode, "_not_found") {
		return CodeMapping{codes.NotFound, gwerrors.ErrCodeNotFound}
	}

	return CodeMapping{codes.Internal, gwerrors.ErrCodeInternalServerError}
}

// MapBizCodeToGRPCCode 业务错误码映射为 gRPC 状态码
func MapBizCodeToGRPCCode(bizCode string) codes.Code {
	return Mapping(bizCode).GRPC
}

// MapBizCodeToGatewayErrorCode 业务错误码映射为网关 AppError 错误码
func MapBizCodeToGatewayErrorCode(bizCode string) gwerrors.ErrorCode {
	return Mapping(bizCode).Gateway
}
