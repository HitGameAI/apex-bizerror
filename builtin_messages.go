/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-08-15 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-08-15 00:00:00
 * @FilePath: \go-bizerror\builtin_messages.go
 * @Description: 通用错误码内置 i18n 文案（服务 locales 未命中时兜底，服务同名 key 优先覆盖）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package bizerror

// builtinMessages 内置通用错误码文案（语言 → 错误码 → 文案）
// 服务侧 locales/*.json 同名 key 优先，此处仅作缺失兜底，新服务零配置即有可读错误消息
var builtinMessages = map[string]map[string]string{
	"en": {
		BizErrCodeInvalidArgument:          "Invalid argument",
		BizErrCodeUnauthenticated:          "Unauthenticated",
		BizErrCodePermissionDenied:         "Permission denied",
		BizErrCodeNotFound:                 "Resource not found",
		BizErrCodeConflict:                 "Resource conflict",
		BizErrCodeValidationFailed:         "Validation failed",
		BizErrCodeInvalidParameter:         "Invalid request parameter: {reason}",
		BizErrCodeCannotBeSelf:             "Cannot perform this operation on yourself",
		BizErrCodeDatabaseError:            "Database error",
		BizErrCodeDatabaseOperationFailed:  "Database operation failed",
		BizErrCodeConversionFailed:         "Data conversion failed",
		BizErrCodeInternalError:            "Internal error",
		BizErrCodeRepositoryNotInitialized: "Repository not initialized",
		BizErrCodeUnsupportedOperation:     "Unsupported operation",
		BizErrCodeExternalServiceError:     "External service error",
	},
	"zh": {
		BizErrCodeInvalidArgument:          "无效的参数",
		BizErrCodeUnauthenticated:          "未认证",
		BizErrCodePermissionDenied:         "权限不足",
		BizErrCodeNotFound:                 "资源未找到",
		BizErrCodeConflict:                 "资源冲突",
		BizErrCodeValidationFailed:         "数据校验失败",
		BizErrCodeInvalidParameter:         "请求参数无效：{reason}",
		BizErrCodeCannotBeSelf:             "不能对自身执行该操作",
		BizErrCodeDatabaseError:            "数据库错误",
		BizErrCodeDatabaseOperationFailed:  "数据库操作失败",
		BizErrCodeConversionFailed:         "数据转换失败",
		BizErrCodeInternalError:            "内部错误",
		BizErrCodeRepositoryNotInitialized: "仓库未初始化",
		BizErrCodeUnsupportedOperation:     "不支持的操作",
		BizErrCodeExternalServiceError:     "外部服务错误",
	},
}

// BuiltinMessage 查询内置文案：指定语言未收录时回退 en，仍无返回空串
//
// 导出是为了让 gwx 子模块在网关 locales 未命中时也能取到同一份兜底文案，
// 避免两边各维护一张表
// [EN] Exported so the gwx submodule can reach the same fallback table when the
// gateway's locales miss a key, instead of keeping a second copy
func BuiltinMessage(language, errorCode string) string {
	if language == "" {
		language = "en"
	}
	if msg, ok := builtinMessages[language][errorCode]; ok {
		return msg
	}
	return builtinMessages["en"][errorCode]
}
