package utils

//redis
import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

// Rd 定义一个全局变量
var client *redis.Client
var ctx = context.Background()

type Redis struct{}

func RedisInit() (err error) {
	client = redis.NewClient(&redis.Options{
		//Addr:     "127.0.0.1:6379", // 本地环境
		Addr:     "r-2zeuu64hewjynuchp8pd.redis.rds.aliyuncs.com:6379", // 测试环境
		Password: "BgS6q5PV7WECWMU3QZvuQCWqJjc@nU",
		DB:       1, // redis一共16个库，指定其中一个库即可
	})
	_, err = client.Ping(ctx).Result()
	return err
}

/*------------------------------------ 字符 操作 ------------------------------------*/

// Set 设置 key的值
func Set(key, value string) bool {
	result, err := client.Set(ctx, key, value, 0).Result()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return result == "OK"
}

// SetEX 设置 key的值并指定过期时间
func SetEX(key, value string, ex time.Duration) bool {
	result, err := client.Set(ctx, key, value, ex).Result()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return result == "OK"
}

// Get 获取 key的值
func Get(key string) (bool, string) {
	result, err := client.Get(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
		return false, ""
	}
	return true, result
}

// GetSet 设置新值获取旧值
func GetSet(key, value string) (bool, string) {
	oldValue, err := client.GetSet(ctx, key, value).Result()
	if err != nil {
		fmt.Println(err)
		return false, ""
	}
	return true, oldValue
}

// Incr key值每次加一 并返回新值
func Incr(key string) int64 {
	val, err := client.Incr(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
	}
	return val
}

// IncrBy key值每次加指定数值 并返回新值
func IncrBy(key string, incr int64) int64 {
	val, err := client.IncrBy(ctx, key, incr).Result()
	if err != nil {
		fmt.Println(err)
	}
	return val
}

// IncrByFloat key值每次加指定浮点型数值 并返回新值
func IncrByFloat(key string, incrFloat float64) float64 {
	val, err := client.IncrByFloat(ctx, key, incrFloat).Result()
	if err != nil {
		fmt.Println(err)
	}
	return val
}

// Decr key值每次递减 1 并返回新值
func Decr(key string) int64 {
	val, err := client.Decr(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
	}
	return val
}

// DecrBy key值每次递减指定数值 并返回新值
func DecrBy(key string, incr int64) int64 {
	val, err := client.DecrBy(ctx, key, incr).Result()
	if err != nil {
		fmt.Println(err)
	}
	return val
}

// Del 删除 key
func Del(key string) bool {
	result, err := client.Del(ctx, key).Result()
	if err != nil {
		return false
	}
	return result == 1
}

// Expire 设置 key的过期时间
func Expire(key string, ex time.Duration) bool {
	result, err := client.Expire(ctx, key, ex).Result()
	if err != nil {
		return false
	}
	return result
}

/*------------------------------------ set 操作 ------------------------------------*/

// SAdd 添加元素到集合中
func SAdd(key string, data ...interface{}) bool {
	err := client.SAdd(ctx, key, data).Err()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// SCard 获取集合元素个数
func SCard(key string) int64 {
	size, err := client.SCard(ctx, "key").Result()
	if err != nil {
		fmt.Println(err)
	}
	return size
}

// SIsMember 判断元素是否在集合中
func SIsMember(key string, data interface{}) bool {
	ok, err := client.SIsMember(ctx, key, data).Result()
	if err != nil {
		fmt.Println(err)
	}
	return ok
}

// SMembers 获取集合所有元素
func SMembers(key string) []string {
	es, err := client.SMembers(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
	}
	return es
}

// SRem 删除 key集合中的 data元素
func SRem(key string, data ...interface{}) bool {
	_, err := client.SRem(ctx, key, data).Result()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// SPopN 随机返回集合中的 count个元素，并且删除这些元素
func SPopN(key string, count int64) []string {
	vales, err := client.SPopN(ctx, key, count).Result()
	if err != nil {
		fmt.Println(err)
	}
	return vales
}

/*------------------------------------ zset 操作 ------------------------------------*/

// ZAdd 将一个 member 元素及其 score 值加入到有序集 key 当中。
func ZAdd(key string, score int64, value interface{}) int64 {
	member := &redis.Z{
		Score:  float64(score),
		Member: value,
	}
	result, err := client.ZAdd(ctx, key, member).Result()
	if err != nil {
		fmt.Println(err)
	}
	return result
}

// ZRem 移除有序集 key 中的一个成员，不存在的成员将被忽略。
func ZRem(key string, member interface{}) int64 {
	result, err := client.ZRem(ctx, key, member).Result()
	if err != nil {
		fmt.Println(err)
	}
	return result
}

// ZRevRange 返回有序集中，指定区间内的成员。其中成员的位置按分数值递减(从大到小)来排列。具有相同分数值的成员按字典序(lexicographical order )来排列。
// 以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。或 以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。
func ZRevRange(key string, from, to int64) []string {
	result, err := client.ZRevRange(ctx, key, from, to).Result()
	if err != nil {
		fmt.Println(err)
	}
	return result
}

// ZRevByScoreWithScores 返回有序集中指定分数区间内的所有的成员。有序集成员按分数值递减(从大到小)的次序排列。
// 具有相同分数值的成员按字典序来排列
func ZRevByScoreWithScores(key string, min, max float64, offset, count int64) []redis.Z {
	var tMax string
	if max == 0 {
		tMax = "+inf"
	} else {
		tMax = strconv.FormatFloat(max, 'f', 1, 64)
	}
	t := &redis.ZRangeBy{
		Min:    strconv.FormatFloat(min, 'f', 1, 64),
		Max:    tMax,
		Offset: offset,
		Count:  count,
	}
	result, err := client.ZRevRangeByScoreWithScores(ctx, key, t).Result()
	if err != nil {
		fmt.Println(err)
	}
	return result
}
