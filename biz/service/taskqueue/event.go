package taskqueue

import (
	"context"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"judgeMore/biz/dal/cache"
	"judgeMore/biz/dal/mysql"
	"judgeMore/biz/service/model"
	"judgeMore/pkg/errno"
	"judgeMore/pkg/taskqueue"
)

// 加入任务队列
func AddScoreEvent(ctx context.Context, key, event_id string) {
	taskQueue.Add(key, taskqueue.QueueTask{Execute: func() error {
		return calculateScore(ctx, event_id)
	}})
}

// 这边 中间报错了怎么办？
func calculateScore(ctx context.Context, event_id string) error {
	exist, err := mysql.IsScoreRecordExist_Event(ctx, event_id)
	if err != nil {
		logger.Errorf("calculateScore:failed to query %v exist info :%v", event_id, err)
		return err
	}
	if exist {
		logger.Errorf(" %v scoreRecord exist info", event_id)
		return nil
	}
	eventInfo, err := mysql.GetEventInfoById(ctx, event_id)
	if err != nil {
		logger.Errorf("calculateScore:failed to query %v info :%v", event_id, err)
		return errno.NewErrNo(errno.InternalDatabaseErrorCode,
			fmt.Sprintf("calculateScore: Redis SET failed: %v", err))
	}
	// cache处获取rule
	exist, err = cache.IsRuleExist(ctx)
	var ruleList []*model.ScoreRule
	if !exist {
		// db 载入 redis
		ruleList, _, err = mysql.GetScoreRule(ctx)
		if err != nil {
			logger.Errorf("calculateScore:failed to query rule info :%v", err)
			return errno.NewErrNo(errno.InternalDatabaseErrorCode,
				fmt.Sprintf("calculateScore: failed to query rule info :%v", err))
		}
		err = cache.RuleToCache(ctx, ruleList)
		if err != nil {
			logger.Errorf("calculateScore:failed to update rule cache :%v", err)
			return errno.NewErrNo(errno.InternalRedisErrorCode,
				fmt.Sprintf("calculateScore: failed update rule cache :%v", err))
		}
	} else {
		ruleList, err = cache.QueryAllRule(ctx)
	}
	var score float64
	score = -1
	// 遍历规则 计算积分
	// 总体原则是 匹配rule的event_level,award_level 如果存在该两项匹配 且rule.reconizeid == event.recognizeid
	// 则应该使用 recognized_id相同的
	// 如果不存在则 只要两者匹配即可
	// scorerule存在以下的情况即
	//	eventlevel awardlevel	recognize_id
	// 1. 国家级     一等奖         0(默认为零)
	// 2. 国家级       一等奖        10086 (代表如果event的recognize_id=10086 则他不能按常规的国家一等计算）

	// 构建两层映射：第一层是 recognizedId，第二层是 eventLevel + awardLevel
	recognizedMap := make(map[string]*model.ScoreRule) // recognizedId -> rule
	levelMap := make(map[string]*model.ScoreRule)      // "eventLevel_awardLevel" -> rule

	for _, rule := range ruleList {
		// 优先处理 recognizedId 不为 0 的规则
		if rule.RecognizedEventId != "0" {
			recognizedMap[rule.RecognizedEventId] = rule
		}
		// 构建 level 映射
		levelKey := fmt.Sprintf("%s_%s", rule.EventLevel, rule.AwardLevel)
		levelMap[levelKey] = rule
	}
	var ruleid string
	if rule, exists := recognizedMap[eventInfo.RecognizeId]; exists {
		score = float64(rule.Integral) * rule.EventWeight
		ruleid = rule.RuleId
	} else {
		// 其次匹配 eventLevel + awardLevel
		levelKey := fmt.Sprintf("%s_%s", eventInfo.EventLevel, eventInfo.AwardLevel)
		if rule, exists := levelMap[levelKey]; exists {
			score = float64(rule.Integral) * rule.EventWeight
			ruleid = rule.RuleId
		}
	}
	if score == -1 {
		logger.Errorf("calculateScore:failed to get event %v score :%v", event_id, "no rule match")
		return errno.NewErrNo(errno.InternalServiceErrorCode, "calculateScore:failed to get event "+event_id+" score :no rule match")
	}
	record := &model.ScoreRecord{
		UserId:        eventInfo.Uid,
		EventId:       eventInfo.EventId,
		RuleId:        ruleid,
		FinalIntegral: score,
	}
	return mysql.CreateNewScoreRecord(ctx, record)
}

func Work(key string) {
	taskQueue.Start()
}
