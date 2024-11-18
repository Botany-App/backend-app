package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/go-redis/redis/v8"
)

const cacheDuration = 6 * time.Hour

func GetFromCache[T any](rd *redis.Client, key string, fetch func() (T, error)) (T, error) {
	var result T
	ctx := context.Background()

	cachedData, err := rd.Get(ctx, key).Result()
	if err == redis.Nil || cachedData == "null" || cachedData == "" {
		log.Println("Cached data is `null` or empty. Fetching from database.")
		fetchedData, fetchErr := fetch()
		if fetchErr != nil {
			return result, fetchErr
		}

		// Verificar se fetchedData é nil
		if reflect.ValueOf(fetchedData).IsZero() {
			log.Println("Fetched data is nil or zero value. Skipping cache set.")
			return result, nil // Retorna result que é o zero value para T
		}

		serializedData, err := json.Marshal(fetchedData)
		if err != nil {
			return result, fmt.Errorf("failed to serialize data for caching: %w", err)
		}

		err = rd.Set(ctx, key, serializedData, cacheDuration).Err()
		if err != nil {
			log.Printf("Failed to cache data: %v", err)
		}
		return fetchedData, nil
	} else if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(cachedData), &result)
	if err != nil {
		log.Printf("Failed to deserialize data from cache: %v", err)
		return result, nil
	}
	return result, nil
}
