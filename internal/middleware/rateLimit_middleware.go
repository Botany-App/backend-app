package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
)

var (
	ctx = context.Background()
)

func RateLimitMiddleware(redisClient *redis.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			limiter := redis_rate.NewLimiter(redisClient)
			clientIP := r.RemoteAddr

			rateLimitKey := fmt.Sprintf("ratelimit:%s", clientIP)
			jailKey := fmt.Sprintf("jail:%s", clientIP)

			jailStatus, err := redisClient.Get(ctx, jailKey).Result()
			if err != nil && err != redis.Nil {
				http.Error(w, "Erro no sistema de Jail", http.StatusInternalServerError)
				return
			}

			if jailStatus == "1" {
				http.Error(w, "Você está bloqueado temporariamente devido a excesso de requisições", http.StatusForbidden)
				return
			}

			res, err := limiter.Allow(ctx, rateLimitKey, redis_rate.PerMinute(100))
			if err != nil {
				http.Error(w, "Erro no Rate Limiter", http.StatusInternalServerError)
				return
			}

			if res.Allowed == 0 {
				failCountKey := fmt.Sprintf("failcount:%s", clientIP)
				failCount, _ := redisClient.Incr(ctx, failCountKey).Result()

				redisClient.Expire(ctx, failCountKey, time.Minute)

				if failCount > 5 {
					redisClient.Set(ctx, jailKey, "1", 5*time.Minute)
					http.Error(w, "Você está bloqueado temporariamente devido a excesso de requisições", http.StatusForbidden)
					return
				}

				w.Header().Set("Retry-After", fmt.Sprintf("%d", res.RetryAfter/time.Second))
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			redisClient.Del(ctx, fmt.Sprintf("failcount:%s", clientIP))

			next.ServeHTTP(w, r)
		})
	}
}
