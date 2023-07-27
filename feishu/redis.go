package feishu

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

func CreateRedisInstance(opType, key string, values ...string) (string, error) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.3.169:6379",
		Password: "123456",
		DB:       20,
	})
	ctx := context.Background()
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Failed to ping redis:", err)
		return err.Error(), err
	}
	fmt.Println("redis.go成功连接到redis: ", pong)
	switch opType {
	case "get":
		val, err := rdb.Get(ctx, key).Result()
		//		fmt.Println("redis.go现在输出的是val", val)
		if err != nil {
			return "", fmt.Errorf("redis.go获取值出错: %s", err)
		}
		return val, nil
	case "set":
		if len(values) == 0 {
			return "", fmt.Errorf("缺少值参数")
		}
		err := rdb.Set(ctx, key, values[0], 0).Err()

		if err != nil {
			return "", fmt.Errorf("设置值出错: %s", err)
		}
		return "", nil
	default:
		return "", fmt.Errorf("不支持的操作类型: %s", opType)
	}
}
