package pack

import (
	resp "judgeMore/biz/model/model"
	"judgeMore/biz/service/model"
	"strconv"
)

func ScoreRecord(data *model.ScoreRecord) *resp.ScoreRecord {
	return &resp.ScoreRecord{
		ScoreID:    data.ResultId,
		EventID:    data.EventId,
		UserID:     data.UserId,
		RuleID:     data.RuleId,
		AppealID:   data.AppealId,
		FinalScore: data.FinalIntegral,
		Status:     data.Status,
		CreatedAt:  strconv.FormatInt(data.CreateAT, 10),
		UpdatedAt:  strconv.FormatInt(data.UpdateAT, 10),
		DeletedAt:  strconv.FormatInt(data.DeleteAT, 10),
	}
}

func ScoreRecordList(data []*model.ScoreRecord, total int64) *resp.ScoreRecordList {
	info := make([]*resp.ScoreRecord, 0)
	var sum float64
	for _, v := range data {
		sum += v.FinalIntegral
		info = append(info, ScoreRecord(v))
	}
	return &resp.ScoreRecordList{
		Items: info,
		Total: total,
		Sum:   sum,
	}
}
