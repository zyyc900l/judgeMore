package mysql

import (
	"context"
	"judgeMore/biz/service/model"
	"judgeMore/pkg/constants"
	"judgeMore/pkg/errno"
)

func QueryRecognizedEvent(ctx context.Context) ([]*model.RecognizedEvent, int64, error) {
	var reconize_event []*RecognizedEvent
	var total int64
	err := db.WithContext(ctx).
		Table(constants.TableReconizedEvent).
		Find(&reconize_event).
		Count(&total).
		Error
	if err != nil {
		return nil, -1, errno.NewErrNo(errno.InternalDatabaseErrorCode, "mysql:failed to query reconized_event"+err.Error())
	}
	return buildRecognizedList(reconize_event), total, nil
}

func buildRecognizeEvent(data *RecognizedEvent) *model.RecognizedEvent {
	return &model.RecognizedEvent{
		RecognizedEventId:   data.RecognizedEventId,
		RecognizedLevel:     data.RecognizedLevel,
		RecognizedEventName: data.RecognizedEventName,
		RecognizedEventTime: data.RecognizedEventTime,
		RecognitionBasis:    data.RecognitionBasis,
		College:             data.College,
		Organizer:           data.Organizer,
		RelatedMajors:       data.RelatedMajors,
		ApplicableMajors:    data.ApplicableMajors,
		IsActive:            data.IsActive,
		UpdateAT:            data.UpdatedAt.Unix(),
		CreateAT:            data.CreatedAt.Unix(),
		DeleteAT:            0,
	}
}
func buildRecognizedList(data []*RecognizedEvent) []*model.RecognizedEvent {
	r := make([]*model.RecognizedEvent, 0)
	for _, v := range data {
		r = append(r, buildRecognizeEvent(v))
	}
	return r
}
