package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisConnector struct {
	*redis.Client
}

func NewRedisConnector() *RedisConnector {
	redisDataBase := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	fmt.Println(redisDataBase)	
	return &RedisConnector{ redisDataBase } 
}

func (rdb *RedisConnector) SetData(key, value string) (error) {
	
	ctx := context.Background();
	err := rdb.Set(ctx, key, value, 0).Err();

	if err != nil {
		return err;
	}

	return nil
}


func (rdb *RedisConnector) GetData(key string) (string, error) {
	
	ctx := context.Background();
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}


func (rdb *RedisConnector) DeleteData(key string) (error) {
	
	ctx := context.Background();
	err := rdb.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}


func (rdb *RedisConnector) DoesDataExists(key string) (*int64, error) {
	println(key)
	println("key")
	ctx := context.Background();
	val, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return &val, nil
}