package do

import (
	"context"
	"fmt"
)

type ContextHelper[K ~struct{}, V any] struct {
	k K
}

func (h ContextHelper[K, V]) WithValue(ctx context.Context, v V) context.Context {
	return context.WithValue(ctx, h.k, v)
}

func (h ContextHelper[K, V]) Value(ctx context.Context) (v V, ok bool) {
	v, ok = ctx.Value(h.k).(V)
	return
}

func (h ContextHelper[K, V]) MustValue(ctx context.Context) (v V) {
	v, ok := h.Value(ctx)
	if !ok {
		panic(fmt.Errorf("context can't find value of %#v", h.k))
	}
	return
}

type (
	// 时间戳
	TimestampType          struct{}
	TimestampTypeCtxHelper = ContextHelper[TimestampType, int64]

	// 远程地址
	RemoteAddrType          struct{}
	RemoteAddrTypeCtxHelper = ContextHelper[RemoteAddrType, string]

	// 用户
	UserKeyType          struct{}
	UserKeyTypeCtxHelper = ContextHelper[UserKeyType, uint64]

	// 请求
	RequestKeyType          = TraceKey
	RequestKeyTypeCtxHelper = ContextHelper[RequestKeyType, string]

	// 数据权限
	CheckDataPerm          struct{}
	CheckDataPermCtxHelper = ContextHelper[CheckDataPerm, bool]

	// 数据权限join语句
	DataPermJoinType          struct{}
	DataPermJoinTypeCtxHelper = ContextHelper[DataPermJoinType, [][]string]

	// 数据权限条件语句
	DataPermType          struct{}
	DataPermTypeCtxHelper = ContextHelper[DataPermType, string]

	// 用户是否超管角色
	IsAdminType          struct{}
	IsAdminTypeCtxHelper = ContextHelper[IsAdminType, bool]

	// 接口信息
	APIType          struct{}
	APITypeCtxHelper = ContextHelper[APIType, string]

	// 环境信息
	EnvType          struct{}
	EnvTypeCtxHelper = ContextHelper[EnvType, string]

	// 无需保存标记
	NotSaveType      struct{}
	NotSaveCtxHelper = ContextHelper[NotSaveType, bool]
)
