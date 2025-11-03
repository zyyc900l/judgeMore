package mysql

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"judgeMore/biz/service/model"
	"judgeMore/pkg/constants"
	"judgeMore/pkg/errno"
)

func IsScoreRecordExist(ctx context.Context, scoreId string) (bool, error) {
	var Info *ScoreResult
	err := db.WithContext(ctx).
		Table(constants.TableScore).
		Where("result_id = ?", scoreId).
		First(&Info).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query score exist : %v", err)
	}
	return true, nil
}
func IsScoreRecordExist_Event(ctx context.Context, eventId string) (bool, error) {
	var Info *ScoreResult
	err := db.WithContext(ctx).
		Table(constants.TableScore).
		Where("event_id = ?", eventId).
		First(&Info).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query score exist: %v", err)
	}
	return true, nil
}
func QueryScoreRecordByEventId(ctx context.Context, eventId string) (*model.ScoreRecord, error) {
	var Info *ScoreResult
	err := db.WithContext(ctx).
		Table(constants.TableScore).
		Where("event_id = ?", eventId).
		First(&Info).
		Error
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query ScoreRecord by event_id: %v", err)
	}
	return buildScore(Info), nil
}

func QueryScoreRecordByScoreId(ctx context.Context, scoreId string) (*model.ScoreRecord, error) {
	var Info *ScoreResult
	err := db.WithContext(ctx).
		Table(constants.TableScore).
		Where("result_id = ?", scoreId).
		First(&Info).
		Error
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query ScoreRecord by result_id: %v", err)
	}
	return buildScore(Info), nil
}
func QueryScoreRecordByStuId(ctx context.Context, stuId string) ([]*model.ScoreRecord, int64, error) {
	var Info []*ScoreResult
	var count int64
	err := db.WithContext(ctx).
		Table(constants.TableScore).
		Where("user_id = ?", stuId).
		Find(&Info).
		Count(&count).
		Error
	if err != nil {
		return nil, -1, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query ScoreRecord by Stu_id: %v", err)
	}
	return buildScoreList(Info), count, nil
}

// 调用前用event_id查找一遍
func CreateNewScoreRecord(ctx context.Context, record *model.ScoreRecord) error {
	var r *ScoreResult
	r = &ScoreResult{
		UserId:        record.UserId,
		RuleId:        record.RuleId,
		EventId:       record.EventId,
		FinalIntegral: record.FinalIntegral,
		Status:        "正常",
	}
	err := db.WithContext(ctx).
		Table(constants.TableScore).
		Create(&r).
		Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "CreateScoreRecord error"+err.Error())
	}
	return nil
}

func GetScoreRule(ctx context.Context) ([]*model.ScoreRule, int64, error) {
	var total int64
	var rules []*EventRule
	err := db.WithContext(ctx).
		Table(constants.TableRule).
		Find(&rules).
		Count(&total).
		Error
	if err != nil {
		return nil, -1, errno.NewErrNo(errno.InternalDatabaseErrorCode, "mysql:query scoreRule failed:"+err.Error())
	}
	return buildRuleList(rules), total, nil
}

func buildRule(rule *EventRule) *model.ScoreRule {
	return &model.ScoreRule{
		RuleId:            rule.RuleId,
		RecognizedEventId: rule.RecognizedEventId,
		EventLevel:        rule.EventLevel,
		EventWeight:       rule.EventWeight,
		Integral:          rule.Integral,
		RuleDesc:          rule.RuleDesc,
		IsEditable:        rule.IsEditable,
		AwardLevel:        rule.AwardLevel,
		AwardLevelWeight:  rule.AwardLevelWeight,
		UpdateAT:          rule.UpdatedAt.Unix(),
		CreateAT:          rule.CreatedAt.Unix(),
		DeleteAT:          0,
	}
}
func buildRuleList(rules []*EventRule) []*model.ScoreRule {
	list := make([]*model.ScoreRule, 0)
	for _, rule := range rules {
		list = append(list, buildRule(rule))
	}
	return list
}

func buildScore(r *ScoreResult) *model.ScoreRecord {
	return &model.ScoreRecord{
		RuleId:        r.RuleId,
		UserId:        r.UserId,
		ResultId:      r.ResultId,
		AppealId:      r.AppealId,
		EventId:       r.EventId,
		Status:        r.Status,
		FinalIntegral: r.FinalIntegral,
		UpdateAT:      r.UpdatedAt.Unix(),
		CreateAT:      r.CreatedAt.Unix(),
		DeleteAT:      0,
	}
}
func buildScoreList(rules []*ScoreResult) []*model.ScoreRecord {
	list := make([]*model.ScoreRecord, 0)
	for _, rule := range rules {
		list = append(list, buildScore(rule))
	}
	return list
}
