package middleware

import (
	"net/http"
	"strconv"
	"time"

	"Taurus/pkg/httpx"
	"Taurus/pkg/redisx"
	"Taurus/pkg/util"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 通过http header中的token解析来认证
		token := r.Header.Get("token")
		setJwtToTrace(r, token)
		if token == "" {
			httpx.SendResponse(w, http.StatusUnauthorized, "Jwt Token is empty", nil)
			return
		}

		// 解析token中包含的相关信息（有效载荷）
		claims, err := util.ParseToken(token)
		if err != nil {
			httpx.SendResponse(w, http.StatusUnauthorized, "Jwt Token is parse error", nil)
			return
		}

		ua := r.Header.Get("User-Agent")
		val, err := redisx.Redis.HGet(r.Context(), strconv.Itoa(int(claims.Uid)), ua)

		if err != nil { // 说明该token是其他User-Agent的token（比如说电脑端的token的map key 是User-Agent，当然不能用来登录手机端）
			httpx.SendResponse(w, http.StatusUnauthorized, "Jwt Token form redis error", nil)
			return
		}

		if token != val { // 请求携带的token与redis中存储的token不一致，说明是旧的token
			httpx.SendResponse(w, http.StatusUnauthorized, "Jwt Token is invalid", nil)
			return
		}
		// 处理过期token
		if time.Now().Unix() > claims.ExpiresAt {
			httpx.SendResponse(w, http.StatusUnauthorized, "Jwt Token is expired", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func setJwtToTrace(r *http.Request, token string) {
	if span := trace.SpanFromContext(r.Context()); span.SpanContext().IsValid() {
		span.SetAttributes(attribute.String("Jwt", token))
	}
}

// ------------------  例子 ------------------
/*
登录成功的时候需要签发JWT, 分发token（其他功能需要身份验证，给前端存储的）,
token返回给前端以后， 前端需要存储下来，每次请求要带过来, 服务端需要将token存储到Redis, 方便后面校验

// -----> 登录成功，存token <-----
token, err := util.GenerateToken(user.ID, user.UserName) // 生产token
if err != nil {
	return httpx.Response{
		Status: http.StatusInternalServerError,
		Msg:    "token签发失败！",
		Error:  err.Error(),
	}
}
// 签发token后，存储到redis中（为了保证token唯一有效）
ua := r.Header.Get("User-Agent") // key用户user-agent可以保证换了浏览器token失效
m := map[string]string{ua: token}
redisx.Redis.HSet(r.Context(), strconv.FormatUint(uint64(user.ID), 10), m)
return response.Response{
	Status: http.StatusOK,
	Msg:    "登录成功！",
	Data:   map[string]string{"token": token},
}

// -----> 中间件校验token <------
// 判断该token是不是最新token（从redis里查）
ua := r.Header.Get("User-Agent")
// 存的时候  key = userid  value = map["User-Agent"]token, 取的时候  取 UID 对于的 map里面的key="User-Agent"对应的值
val, err := redisx.Redis.HGet(r.Context(), strconv.Itoa(int(claims.Uid)), ua).Result()
*/
