package main

import (
	"context"
	"fmt"

	"github.com/donnol/do"
)

// ===== BookServiceMock =====

type BookServiceMock struct {
	WithUserFunc func(ctx context.Context, id uint64) Book
}

var (
	bookServiceMockCommonProxyContext = do.ProxyContext{
		PkgPath:       "github.com/donnol/do/example/proxy",
		InterfaceName: "BookService",
	}

	// represent BookService.WithUser: func(ctx context.Context, id uint64) Book
	BookServiceMockWithUserProxyContext = func() (pctx do.ProxyContext) {
		pctx = bookServiceMockCommonProxyContext
		pctx.MethodName = "WithUser"
		return
	}()

	BookServiceMockProxyContextAll = []do.ProxyContext{
		BookServiceMockWithUserProxyContext,
	}
)

// GetBookServiceProxy 获取接口代理；若使用泛型，需传入typeParams，其值为类型参数的字符串字面量；若想进一步修改方法行为，可以使用RegisterProxyMethod函数注入自定义方法实现；如果想要为每个实例单独注入方法，则使用第二个返回值对象来设置
func GetBookServiceProxy(base BookService, typeParams ...string) (BookService, *do.ProxyCtxFuncStore) {
	if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	_gen_innerCtxMap := do.NewProxyCtxMap()
	return &BookServiceMock{
		WithUserFunc: func(ctx context.Context, id uint64) Book {
			var _gen_ctx = BookServiceMockWithUserProxyContext

			_gen_stop := do.ProxyTraceBegin(_gen_ctx, ctx, id)
			defer func() {
				_gen_stop()
			}()

			var _gen_r0 Book

			var _gen_actual_cf do.ProxyCtxFunc

			_gen_inner_cf, _gen_inner_ok := _gen_innerCtxMap.Lookup(_gen_ctx, typeParams...)
			_gen_cf, _gen_ok := do.GlobalProxyCtxMap().Lookup(_gen_ctx, typeParams...)
			if _gen_inner_ok {
				_gen_actual_cf = _gen_inner_cf
			} else if _gen_ok {
				_gen_actual_cf = _gen_cf
			}

			if _gen_actual_cf != nil {
				_gen_params := []any{}

				_gen_params = append(_gen_params, ctx)
				_gen_params = append(_gen_params, id)

				_gen_res := _gen_actual_cf(_gen_ctx, base.WithUser, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(Book)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.WithUser(ctx, id)
			}

			return _gen_r0
		},
	}, _gen_innerCtxMap
}

func (mockRecv *BookServiceMock) WithUser(ctx context.Context, id uint64) Book {
	return mockRecv.WithUserFunc(ctx, id)
}

// ===== BookStoreMock =====

type BookStoreMock struct {
	ByIdFunc func(ctx context.Context, id uint64) Book

	HookFunc func()
}

var (
	bookStoreMockCommonProxyContext = do.ProxyContext{
		PkgPath:       "github.com/donnol/do/example/proxy",
		InterfaceName: "BookStore",
	}

	// represent BookStore.ById: func(ctx context.Context, id uint64) Book
	BookStoreMockByIdProxyContext = func() (pctx do.ProxyContext) {
		pctx = bookStoreMockCommonProxyContext
		pctx.MethodName = "ById"
		return
	}()

	// represent BookStore.Hook: func()
	BookStoreMockHookProxyContext = func() (pctx do.ProxyContext) {
		pctx = bookStoreMockCommonProxyContext
		pctx.MethodName = "Hook"
		return
	}()

	BookStoreMockProxyContextAll = []do.ProxyContext{
		BookStoreMockByIdProxyContext,
		BookStoreMockHookProxyContext,
	}
)

// GetBookStoreProxy 获取接口代理；若使用泛型，需传入typeParams，其值为类型参数的字符串字面量；若想进一步修改方法行为，可以使用RegisterProxyMethod函数注入自定义方法实现；如果想要为每个实例单独注入方法，则使用第二个返回值对象来设置
func GetBookStoreProxy(base BookStore, typeParams ...string) (BookStore, *do.ProxyCtxFuncStore) {
	if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	_gen_innerCtxMap := do.NewProxyCtxMap()
	return &BookStoreMock{
		ByIdFunc: func(ctx context.Context, id uint64) Book {
			var _gen_ctx = BookStoreMockByIdProxyContext

			_gen_stop := do.ProxyTraceBegin(_gen_ctx, ctx, id)
			defer func() {
				_gen_stop()
			}()

			var _gen_r0 Book

			var _gen_actual_cf do.ProxyCtxFunc

			_gen_inner_cf, _gen_inner_ok := _gen_innerCtxMap.Lookup(_gen_ctx, typeParams...)
			_gen_cf, _gen_ok := do.GlobalProxyCtxMap().Lookup(_gen_ctx, typeParams...)
			if _gen_inner_ok {
				_gen_actual_cf = _gen_inner_cf
			} else if _gen_ok {
				_gen_actual_cf = _gen_cf
			}

			if _gen_actual_cf != nil {
				_gen_params := []any{}

				_gen_params = append(_gen_params, ctx)
				_gen_params = append(_gen_params, id)

				_gen_res := _gen_actual_cf(_gen_ctx, base.ById, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(Book)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.ById(ctx, id)
			}

			return _gen_r0
		},

		HookFunc: func() {
			var _gen_ctx = BookStoreMockHookProxyContext

			_gen_stop := do.ProxyTraceBegin(_gen_ctx)
			defer func() {
				_gen_stop()
			}()

			var _gen_actual_cf do.ProxyCtxFunc

			_gen_inner_cf, _gen_inner_ok := _gen_innerCtxMap.Lookup(_gen_ctx, typeParams...)
			_gen_cf, _gen_ok := do.GlobalProxyCtxMap().Lookup(_gen_ctx, typeParams...)
			if _gen_inner_ok {
				_gen_actual_cf = _gen_inner_cf
			} else if _gen_ok {
				_gen_actual_cf = _gen_cf
			}

			if _gen_actual_cf != nil {
				_gen_params := []any{}

				_gen_actual_cf(_gen_ctx, base.Hook, _gen_params)

			} else {
				base.Hook()
			}

		},
	}, _gen_innerCtxMap
}

func (mockRecv *BookStoreMock) ById(ctx context.Context, id uint64) Book {
	return mockRecv.ByIdFunc(ctx, id)
}

func (mockRecv *BookStoreMock) Hook() {
	mockRecv.HookFunc()
}

// ===== UserStoreMock =====

type UserStoreMock struct {
	ByIdFunc func(ctx context.Context, id uint64) User
}

var (
	userStoreMockCommonProxyContext = do.ProxyContext{
		PkgPath:       "github.com/donnol/do/example/proxy",
		InterfaceName: "UserStore",
	}

	// represent UserStore.ById: func(ctx context.Context, id uint64) User
	UserStoreMockByIdProxyContext = func() (pctx do.ProxyContext) {
		pctx = userStoreMockCommonProxyContext
		pctx.MethodName = "ById"
		return
	}()

	UserStoreMockProxyContextAll = []do.ProxyContext{
		UserStoreMockByIdProxyContext,
	}
)

// GetUserStoreProxy 获取接口代理；若使用泛型，需传入typeParams，其值为类型参数的字符串字面量；若想进一步修改方法行为，可以使用RegisterProxyMethod函数注入自定义方法实现；如果想要为每个实例单独注入方法，则使用第二个返回值对象来设置
func GetUserStoreProxy(base UserStore, typeParams ...string) (UserStore, *do.ProxyCtxFuncStore) {
	if base == nil {
		panic(fmt.Errorf("base cannot be nil"))
	}
	_gen_innerCtxMap := do.NewProxyCtxMap()
	return &UserStoreMock{
		ByIdFunc: func(ctx context.Context, id uint64) User {
			var _gen_ctx = UserStoreMockByIdProxyContext

			_gen_stop := do.ProxyTraceBegin(_gen_ctx, ctx, id)
			defer func() {
				_gen_stop()
			}()

			var _gen_r0 User

			var _gen_actual_cf do.ProxyCtxFunc

			_gen_inner_cf, _gen_inner_ok := _gen_innerCtxMap.Lookup(_gen_ctx, typeParams...)
			_gen_cf, _gen_ok := do.GlobalProxyCtxMap().Lookup(_gen_ctx, typeParams...)
			if _gen_inner_ok {
				_gen_actual_cf = _gen_inner_cf
			} else if _gen_ok {
				_gen_actual_cf = _gen_cf
			}

			if _gen_actual_cf != nil {
				_gen_params := []any{}

				_gen_params = append(_gen_params, ctx)
				_gen_params = append(_gen_params, id)

				_gen_res := _gen_actual_cf(_gen_ctx, base.ById, _gen_params)

				_gen_tmpr0, _gen_exist := _gen_res[0].(User)
				if _gen_exist {
					_gen_r0 = _gen_tmpr0
				}

			} else {
				_gen_r0 = base.ById(ctx, id)
			}

			return _gen_r0
		},
	}, _gen_innerCtxMap
}

func (mockRecv *UserStoreMock) ById(ctx context.Context, id uint64) User {
	return mockRecv.ByIdFunc(ctx, id)
}
