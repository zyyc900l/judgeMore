package pack

import (
	resp "judgeMore/biz/model/model"
	"judgeMore/biz/service/model"
	"strconv"
)

func Rule(rule *model.ScoreRule) *resp.Rule {
	return &resp.Rule{
		RuleID:            rule.RuleId,
		RecognizedEventID: rule.RecognizedEventId,
		EventLevel:        rule.EventLevel,
		EventWeight:       rule.EventWeight,
		AwardLevel:        rule.AwardLevel,
		Integral:          rule.Integral,
		RuleDesc:          rule.RuleDesc,
		IsEditable:        rule.IsEditable,
		CreatedAt:         strconv.FormatInt(rule.CreateAT, 10),
		UpdatedAt:         strconv.FormatInt(rule.UpdateAT, 10),
		DeletedAt:         strconv.FormatInt(rule.DeleteAT, 10),
	}
}
func RuleList(rule []*model.ScoreRule, total int64) *resp.RuleList {
	result := make([]*resp.Rule, 0)
	for _, v := range rule {
		result = append(result, Rule(v))
	}
	return &resp.RuleList{
		Item:  result,
		Total: total,
	}
}
