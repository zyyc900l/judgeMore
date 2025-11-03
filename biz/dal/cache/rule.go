package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"judgeMore/biz/service/model"
	"judgeMore/pkg/errno"
	"time"
)

// 积分计算规则 属于高频访问，存一份redis作为业务处理

func IsRuleExist(ctx context.Context) (bool, error) {
	keys, err := scoreCa.Keys(ctx, "rule_*").Result()
	if err != nil {
		return false, errno.NewErrNo(errno.InternalRedisErrorCode, "get rule keys from redis error:"+err.Error())
	}
	if len(keys) == 0 {
		return false, nil
	}
	return true, nil
}

func RuleToCache(ctx context.Context, rule []*model.ScoreRule) error {
	for _, r := range rule {
		key := fmt.Sprintf("rule_%v", r.RuleId)
		// 使用 JSON 序列化
		info, err := json.Marshal(r)
		if err != nil {
			return errno.NewErrNo(errno.InternalServiceErrorCode, "marshal rule to json error:"+err.Error())
		}
		expiration := 72 * time.Hour
		err = scoreCa.Set(ctx, key, info, expiration).Err()
		if err != nil {
			return errno.NewErrNo(errno.InternalRedisErrorCode, "write rule to cache error:"+err.Error())
		}
	}
	return nil
}

// 调用前检验rule存在
func QueryAllRule(ctx context.Context) ([]*model.ScoreRule, error) {
	keys, err := scoreCa.Keys(ctx, "rule_*").Result()
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalRedisErrorCode, "Query rule fail:get rule keys error:"+err.Error())
	}
	pipe := scoreCa.Pipeline()
	for _, key := range keys {
		pipe.Get(ctx, key)
	}
	cmds, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, errno.NewErrNo(errno.InternalRedisErrorCode, "Query rule fail:"+err.Error())
	}
	rules := make([]*model.ScoreRule, 0)
	for _, cmd := range cmds {
		getCmd, ok := cmd.(*redis.StringCmd)
		if !ok {
			continue
		}
		data, err := getCmd.Result()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			continue
		}
		var rule model.ScoreRule
		err = json.Unmarshal([]byte(data), &rule)
		if err != nil {
			continue
		}
		rules = append(rules, &rule)
	}
	return rules, nil
}
