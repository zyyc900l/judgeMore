package pack

import (
	resp "judgeMore/biz/model/model"
	"judgeMore/biz/service/model"
	"strconv"
)

func Event(data *model.Event) *resp.Event {
	return &resp.Event{
		UserID:         data.Uid,
		EventID:        data.EventId,
		MaterialStatus: data.MaterialStatus,
		RecognizeID:    data.RecognizeId,
		MaterialURL:    data.MaterialUrl,
		EventName:      data.EventName,
		EventOrganizer: data.EventOrganizer,
		EventLevel:     data.EventLevel,
		AwardLevel:     data.AwardLevel,
		AwardContent:   data.AwardContent,
		AutoExtracted:  data.AutoExtracted,
		AwardTime:      data.AwardTime,
		CreatedAt:      strconv.FormatInt(data.CreateAT, 10),
		UpdatedAt:      strconv.FormatInt(data.UpdateAT, 10),
		DeletedAt:      strconv.FormatInt(data.DeleteAT, 10),
	}
}
func EventList(data []*model.Event, total int64) *resp.EventList {
	info := make([]*resp.Event, 0)
	for _, v := range data {
		info = append(info, Event(v))
	}
	return &resp.EventList{
		Items: info,
		Total: total,
	}
}
