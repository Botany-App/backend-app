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
	// Contexto global
	ctx = context.Background()
)

// RateLimitMiddleware aplica o rate limit e o jail
func RateLimitMiddleware(redisClient *redis.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			limiter := redis_rate.NewLimiter(redisClient)
			clientIP := r.RemoteAddr

			// Chave para o rate limit e jail
			rateLimitKey := fmt.Sprintf("ratelimit:%s", clientIP)
			jailKey := fmt.Sprintf("jail:%s", clientIP)

			// Verificar se o cliente está na "jail"
			jailStatus, err := redisClient.Get(ctx, jailKey).Result()
			if err != nil && err != redis.Nil {
				http.Error(w, "Erro no sistema de Jail", http.StatusInternalServerError)
				return
			}

			if jailStatus == "1" {
				http.Error(w, "Você está bloqueado temporariamente devido a excesso de requisições", http.StatusForbidden)
				return
			}

			// Aplicar rate limit
			res, err := limiter.Allow(ctx, rateLimitKey, redis_rate.PerMinute(1))
			if err != nil {
				http.Error(w, "Erro no Rate Limiter", http.StatusInternalServerError)
				return
			}

			if res.Allowed == 0 {
				// Incrementa contador de falhas no Redis
				failCountKey := fmt.Sprintf("failcount:%s", clientIP)
				failCount, _ := redisClient.Incr(ctx, failCountKey).Result()

				// Definir um TTL para o contador de falhas (ex.: 1 minuto)
				redisClient.Expire(ctx, failCountKey, time.Minute)

				// Bloquear o cliente se exceder 5 falhas em 1 minuto
				if failCount > 5 {
					// Coloca o cliente na "jail" por 5 minutos
					redisClient.Set(ctx, jailKey, "1", 5*time.Minute)
					http.Error(w, "Você está bloqueado temporariamente devido a excesso de requisições", http.StatusForbidden)
					return
				}

				// Informar o tempo para próxima tentativa
				w.Header().Set("Retry-After", fmt.Sprintf("%d", res.RetryAfter/time.Second))
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			// Resetar contador de falhas em caso de requisição bem-sucedida
			redisClient.Del(ctx, fmt.Sprintf("failcount:%s", clientIP))

			// Passar para o próximo middleware/handler
			next.ServeHTTP(w, r)
		})
	}
}
