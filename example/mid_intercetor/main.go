package main

import (
	"context"
	"errors"
	"fmt"
)

// 处理 函数
type Handler func(ctx context.Context, data interface{}) (interface{}, error)

func DoHandler(ctx context.Context, data interface{}) (interface{}, error) {
	fmt.Println("doHandler start")
	fmt.Println("doHandler data:", data)
	fmt.Println("doHandler end")
	return data, nil
}

// 中间件 (装饰器模式), 装饰器的传入和传出就是处理函数handler
type Middleware func(Handler) Handler

func LoginMiddleware(next Handler) Handler {
	return func(ctx context.Context, data interface{}) (interface{}, error) {
		fmt.Println("login middleware start")
		data, err := next(ctx, data)
		fmt.Println("login middleware end")
		return " login middleware " + data.(string), err
	}
}

func AuthMiddleware(next Handler) Handler {
	return func(ctx context.Context, data interface{}) (interface{}, error) {
		fmt.Println("auth middleware start")
		data, err := next(ctx, data)
		fmt.Println("auth middleware end")
		return " auth middleware " + data.(string), err
	}
}

// 创建中间件方法， 更灵活
func CreateMiddleware(flag bool) Middleware {
	return func(next Handler) Handler {
		if flag {
			fmt.Println("-------> middleware flag is true")
		} else {
			fmt.Println("-------> middleware flag is false")
		}

		/*
			return func(ctx context.Context, data interface{}) (interface{}, error) {
				fmt.Println("this is a middleware with create factory")
				return next(ctx, data)
			}
		*/
		return next
	}
}

// 中间件链 调用顺序从后往前
func chainMiddleware(middlewares ...Middleware) Middleware {
	return func(next Handler) Handler {
		chain := next
		for i := len(middlewares) - 1; i >= 0; i-- {
			chain = middlewares[i](chain)
		}
		return chain
	}
}

// 拦截器（观察者模式）, 定义拦截器，注意： 拦截器其实就是handler，并且入参多了 handler整个函数
type Interceptor func(ctx context.Context, data interface{}, next Handler) (interface{}, error)

func AuthInterceptor(ctx context.Context, data interface{}, next Handler) (interface{}, error) {
	fmt.Println("auth interceptor start")
	// 模拟认证失败，被拦截
	if data.(string) == "test" {
		return nil, errors.New("auth failed")
	}
	data, err := next(ctx, data)
	fmt.Println("auth interceptor end")
	return data, err
}

func CreateRateLimitInterceptor(limit int) Interceptor {
	return func(ctx context.Context, data interface{}, next Handler) (interface{}, error) {
		fmt.Println("create rate limit interceptor start limit:", limit)
		data, err := next(ctx, data)
		fmt.Println("create rate limit interceptor end limit:", limit)
		return data, err
	}
}

func chainInterceptor(interceptors ...Interceptor) Interceptor {
	return func(ctx context.Context, data interface{}, next Handler) (interface{}, error) {
		chain := next
		for i := len(interceptors) - 1; i >= 0; i-- {

			chain = func(next Handler, interceptor Interceptor) Handler {
				return func(ctx context.Context, data interface{}) (interface{}, error) {
					return interceptor(ctx, data, next)
				}
			}(chain, interceptors[i])
		}
		return chain(ctx, data)
	}
}

// 将中间件和拦截器综合一下

func chainMiddlewareAndInterceptor(middlewares []Middleware, interceptors []Interceptor) Interceptor {

	return func(ctx context.Context, data interface{}, next Handler) (interface{}, error) {
		chain := next
		for i := len(middlewares) - 1; i >= 0; i-- {
			chain = middlewares[i](chain)
		}

		// 处理拦截器
		for i := len(interceptors) - 1; i >= 0; i-- {
			chain = func(next Handler, interceptor Interceptor) Handler {
				return func(ctx context.Context, data interface{}) (interface{}, error) {
					return interceptor(ctx, data, next)
				}
			}(chain, interceptors[i])
		}

		return chain(ctx, data)
	}

}

func main() {

	/*
		// 中间件链
		chain := chainMiddleware(LoginMiddleware, AuthMiddleware, CreateMiddleware(true))
		handler := chain(DoHandler)
		data, err := handler(context.Background(), "test")
		if err != nil {
			fmt.Println("error:", err)
		}
		fmt.Println("----> data:", data)
	*/

	/*
		// 拦截器
		chainInterceptor := chainInterceptor(AuthInterceptor, CreateRateLimitInterceptor(10))
		data, err := chainInterceptor(context.Background(), "test", DoHandler)
		if err != nil {
			fmt.Println("error:", err)
		}
		fmt.Println("----> data:", data)
	*/

	// 将中间件和拦截器综合一下
	chainMiddlewareAndInterceptor := chainMiddlewareAndInterceptor(
		[]Middleware{LoginMiddleware, AuthMiddleware, CreateMiddleware(true)},
		[]Interceptor{AuthInterceptor, CreateRateLimitInterceptor(10)})
	data, err := chainMiddlewareAndInterceptor(context.Background(), "online", DoHandler)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("----> data:", data)
}
