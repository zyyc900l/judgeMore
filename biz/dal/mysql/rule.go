package mysql

import (
	"context"
	"gorm.io/gorm"
	"judgeMore/biz/service/model"
	"judgeMore/pkg/constants"
	"judgeMore/pkg/errno"
)

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
func AddScoreRule(ctx context.Context, re *model.ScoreRule) (*model.ScoreRule, error) {
	r := &EventRule{
		EventLevel:        re.EventLevel,
		AwardLevel:        re.AwardLevel,
		RecognizedEventId: re.RecognizedEventId,
		RuleDesc:          re.RuleDesc,
		EventWeight:       re.EventWeight,
		Integral:          re.Integral,
		IsEditable:        re.IsEditable,
	}
	err := db.WithContext(ctx).
		Table(constants.TableRule).
		Create(&r).
		Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "mysql:failed to create new Rule"+err.Error())
	}
	return buildRule(r), nil
}
func DeleteRule(ctx context.Context, id string) error {
	err := db.WithContext(ctx).
		Table(constants.TableRule).
		Where("rule_id = ?", id).
		Update("is_active", 0).
		Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "mysql:failed to delete rule event:"+err.Error())
	}
	return nil
}
func UpdateRule(ctx context.Context, r *model.ScoreRule) (*model.ScoreRule, error) {
	err := db.WithContext(ctx).Table(constants.TableRule).Transaction(
		func(tx *gorm.DB) error {
			var err error
			if r.RuleDesc != "" {
				err = tx.Where("rule_id = ?", r.RuleId).
					Update("rule_desc", r.RuleDesc).
					Error
				if err != nil {
					return errno.NewErrNo(errno.InternalDatabaseErrorCode, "update rule error :"+err.Error())
				}
			}
			if r.Integral != 0 {
				err = tx.Where("rule_id = ?", r.RuleId).
					Update("integral", r.Integral).
					Error
				if err != nil {
					return errno.NewErrNo(errno.InternalDatabaseErrorCode, "update rule error :"+err.Error())
				}
			}
			if r.EventWeight != 0 {
				err = tx.Where("rule_id = ?", r.RuleId).
					Update("event_weight", r.EventWeight).
					Error
				if err != nil {
					return errno.NewErrNo(errno.InternalDatabaseErrorCode, "update rule error :"+err.Error())
				}
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	var rule *EventRule
	err = db.WithContext(ctx).Table(constants.TableRule).Where("rule_id = ?", r.RuleId).First(&rule).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "update rule error :"+err.Error())
	}
	return buildRule(rule), nil
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
		IsActive:          rule.IsActive,
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
