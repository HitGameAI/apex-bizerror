/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-08-15 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-08-15 01:16:16
 * @FilePath: \go-bizerror\codes.go
 * @Description: 跨服务通用业务错误码常量（服务专属错误码请在各服务内定义并经 Register 注册映射）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package bizerror

const (
	// ============================================
	// 通用错误
	// ============================================

	BizErrCodeInvalidArgument         = "error.invalid_argument"
	BizErrCodeUnauthenticated         = "error.unauthenticated"
	BizErrCodePermissionDenied        = "error.permission_denied"
	BizErrCodeNotFound                = "error.not_found"
	BizErrCodeConflict                = "error.conflict"
	BizErrCodeValidationFailed        = "error.validation_failed"
	BizErrCodeInvalidParameter        = "error.invalid_parameter"
	BizErrCodeCannotBeSelf            = "error.cannot_be_self"
	BizErrCodeDatabaseError           = "error.database_error"
	BizErrCodeDatabaseOperationFailed = "error.database_operation_failed"
	BizErrCodeConversionFailed        = "error.conversion_failed"

	// ============================================
	// 系统内部错误
	// ============================================

	BizErrCodeInternalError            = "error.internal"
	BizErrCodeRepositoryNotInitialized = "error.repository_not_initialized"
	BizErrCodeUnsupportedOperation     = "error.unsupported_operation"
	BizErrCodeExternalServiceError     = "error.external_service_error"
)
