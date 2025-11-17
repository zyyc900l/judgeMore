package pack

import (
	resp "judgeMore/biz/model/model"
	"judgeMore/biz/service/model"
	"strconv"
)

func RecognizedEvent(event *model.RecognizedEvent) *resp.RecognizeReward {
	return &resp.RecognizeReward{
		RecognizeRewardID: event.RecognizedEventId,
		College:           event.College,
		EventName:         event.RecognizedEventName,
		Organizer:         event.Organizer,
		EventTime:         event.RecognizedEventTime,
		RelatedMajors:     event.RelatedMajors,
		ApplicableMajors:  event.ApplicableMajors,
		RecognitionBasis:  event.RecognitionBasis,
		RecognizedLevel:   event.RecognizedLevel,
		IsActive:          event.IsActive,
		CreatedAt:         strconv.FormatInt(event.CreateAT, 10),
		UpdatedAt:         strconv.FormatInt(event.UpdateAT, 10),
		DeletedAt:         strconv.FormatInt(event.DeleteAT, 10),
	}
}
func RecognizedEventList(event []*model.RecognizedEvent, total int64) *resp.RecognizeRewardList {
	result := make([]*resp.RecognizeReward, 0)
	for _, v := range event {
		result = append(result, RecognizedEvent(v))
	}
	return &resp.RecognizeRewardList{
		Item:  result,
		Total: total,
	}
}
