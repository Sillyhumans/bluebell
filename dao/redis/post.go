package redis

import (
	"bluebell/models"
	"fmt"

	"github.com/go-redis/redis"
)

func SetPostScore(postID int64, score float64) error {
	_, err := client.ZAdd(KeyPostScoreZSet, redis.Z{Score: score, Member: postID}).Result()
	fmt.Printf("score:%f\n", score)
	return err
}

func GetPostOrderScore(p *models.ParaPostList) (ids []string, err error) {
	// 从redis获取id
	key := Prefix + KeyPostScoreZSet
	start := (p.Page - 1) * p.Size
	end := start*p.Size - 1
	ids, err = client.ZRevRange(key, start, end).Result()
	if err != nil {
		return
	}
	return
}
