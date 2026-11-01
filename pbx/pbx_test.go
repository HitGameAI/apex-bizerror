/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-08-16 15:51:36
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-08-16 16:05:12
 * @FilePath: \go-bizerror\pbx\pbx_test.go
 * @Description: Model-PB 转换封装测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package pbx

import (
	"context"
	"testing"

	bizerror "github.com/HitGameAI/apex-bizerror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userModel struct {
	ID   int64
	Name string
}

type userPB struct {
	ID   int64
	Name string
}

func TestToPB(t *testing.T) {
	pb, err := ToPB[userModel, userPB](context.Background(), &userModel{ID: 7, Name: "kamalyes"})
	require.NoError(t, err)
	require.NotNil(t, pb)
	assert.Equal(t, int64(7), pb.ID)
	assert.Equal(t, "kamalyes", pb.Name)
}

func TestToPBs(t *testing.T) {
	pbs, err := ToPBs[userModel, userPB](context.Background(), []*userModel{
		{ID: 1, Name: "a"}, {ID: 2, Name: "b"},
	})
	require.NoError(t, err)
	require.Len(t, pbs, 2)
	assert.Equal(t, "b", pbs[1].Name)
}

// TestToPBOrNotFound_Nil 仓储查不到时返回 NotFound，而不是把 nil 当成空对象返回
func TestToPBOrNotFound_Nil(t *testing.T) {
	pb, err := ToPBOrNotFound[userModel, userPB](context.Background(), nil, "")
	assert.Nil(t, pb)
	require.Error(t, err)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Equal(t, "Resource not found", st.Message())
}

// TestToPBOrNotFound_CustomCode 调用方可以指定更具体的错误码
func TestToPBOrNotFound_CustomCode(t *testing.T) {
	bizerror.Register("error.user_not_found", bizerror.CodeMapping{
		GRPC: codes.NotFound, Gateway: bizerror.GatewayCodeNotFound,
	})

	_, err := ToPBOrNotFound[userModel, userPB](context.Background(), nil, "error.user_not_found")
	require.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, st.Code())
}

func TestToPBOrNotFound_Found(t *testing.T) {
	pb, err := ToPBOrNotFound[userModel, userPB](context.Background(), &userModel{ID: 9, Name: "x"}, "")
	require.NoError(t, err)
	require.NotNil(t, pb)
	assert.Equal(t, int64(9), pb.ID)
}

func TestFromPB(t *testing.T) {
	m, err := FromPB[userPB, userModel](context.Background(), &userPB{ID: 3, Name: "c"})
	require.NoError(t, err)
	require.NotNil(t, m)
	assert.Equal(t, int64(3), m.ID)
	assert.Equal(t, "c", m.Name)
}

func TestFromPBs(t *testing.T) {
	ms, err := FromPBs[userPB, userModel](context.Background(), []*userPB{{ID: 4}, {ID: 5}})
	require.NoError(t, err)
	require.Len(t, ms, 2)
	assert.Equal(t, int64(5), ms[1].ID)
}

func TestToPBInto(t *testing.T) {
	var dst userPB
	require.NoError(t, ToPBInto[userModel, userPB](context.Background(), &userModel{ID: 11, Name: "d"}, &dst))
	assert.Equal(t, int64(11), dst.ID)
	assert.Equal(t, "d", dst.Name)
}

// TestFromPBInto_PartialUpdate 部分更新：先取出完整模型，再用请求 PB 覆盖，
// 未映射到的字段保持原样
func TestFromPBInto_PartialUpdate(t *testing.T) {
	type fullModel struct {
		ID       int64
		Name     string
		Internal string // 请求 PB 里没有这个字段
	}
	type patchPB struct {
		Name string
	}

	dst := &fullModel{ID: 1, Name: "旧名", Internal: "不该被清空"}
	require.NoError(t, FromPBInto[patchPB, fullModel](context.Background(), &patchPB{Name: "新名"}, dst))

	assert.Equal(t, "新名", dst.Name)
	assert.Equal(t, "不该被清空", dst.Internal, "请求 PB 未覆盖的字段必须保持原值")
	assert.Equal(t, int64(1), dst.ID)
}
