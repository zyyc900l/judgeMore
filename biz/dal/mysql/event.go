package mysql

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"judgeMore/biz/service/model"
	"judgeMore/pkg/constants"
	"judgeMore/pkg/errno"
)

func IsEventExist(ctx context.Context, event_id string) (bool, error) {
	var eventInfo *Event
	err := db.WithContext(ctx).
		Table(constants.TableEvent).
		Where("event_id = ?", event_id).
		First(&eventInfo).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { //没找到了用户不存在
			return false, nil
		}
		return false, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query event: %v", err)
	}
	return true, nil
}
func CreateNewEvent(ctx context.Context, event *model.Event) (string, error) {
	eventInfo := &Event{
		UserId:         event.Uid,
		RecognizedId:   event.RecognizeId,
		EventName:      event.EventName,
		AutoExtracted:  event.AutoExtracted,
		EventOrganizer: event.EventOrganizer,
		EventLevel:     event.EventLevel,
		MaterialUrl:    event.MaterialUrl,
		MaterialStatus: "待审核",
		AwardLevel:     event.AwardLevel, //提取的内容并没有这项
		AwardAt:        event.AwardTime,
	}
	err := db.WithContext(ctx).
		Table(constants.TableEvent).
		Create(&eventInfo).
		Error
	if err != nil {
		return "", err
	}
	return eventInfo.EventId, nil
}

// 该函数调用前检验存在性
func GetEventInfoById(ctx context.Context, event_id string) (*model.Event, error) {
	var eventInfo *Event
	err := db.WithContext(ctx).
		Table(constants.TableEvent).
		Where("event_id = ?", event_id).
		First(&eventInfo).
		Error
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query event: %v", err)
	}
	return &model.Event{
		Uid:            eventInfo.UserId,
		EventId:        eventInfo.EventId,
		AwardLevel:     eventInfo.AwardLevel,
		EventLevel:     eventInfo.EventLevel,
		EventName:      eventInfo.EventName,
		EventOrganizer: eventInfo.EventOrganizer,
		MaterialUrl:    eventInfo.MaterialUrl,
		MaterialStatus: eventInfo.MaterialStatus,
		AutoExtracted:  eventInfo.AutoExtracted,
		AwardTime:      eventInfo.AwardAt,
		UpdateAT:       eventInfo.UpdatedAt.Unix(),
		CreateAT:       eventInfo.CreatedAt.Unix(),
		DeleteAT:       0,
	}, nil
}

func GetEventInfoByStuId(ctx context.Context, stu_id string) ([]*model.Event, int64, error) {
	var eventInfos []*Event
	var count int64
	err := db.WithContext(ctx).
		Table(constants.TableEvent).
		Where("user_id = ?", stu_id).
		Find(&eventInfos).
		Count(&count).
		Error
	if err != nil {
		return nil, -1, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed query stu event: %v", err)
	}
	return buildEventList(eventInfos), count, err
}
func UpdateEventStatus(ctx context.Context, event_id string, status int64) (*model.Event, error) {
	updateFields := make(map[string]interface{})
	switch status {
	case 1:
		updateFields["material_status"] = "已审核"
		break
	case 2:
		updateFields["material_status"] = "驳回"
		break
	}
	err := db.WithContext(ctx).
		Table(constants.TableEvent).
		Where("event_id = ?", event_id).
		Updates(updateFields).
		Error
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update userInfo: %v", err)
	}
	return GetEventInfoById(ctx, event_id)
}

func buildEvent(data *Event) *model.Event {
	return &model.Event{
		EventId:        data.EventId,
		Uid:            data.UserId,
		AwardContent:   data.AwardContent,
		AwardLevel:     data.AwardLevel,
		EventLevel:     data.EventLevel,
		EventName:      data.EventName,
		EventOrganizer: data.EventOrganizer,
		MaterialUrl:    data.MaterialUrl,
		MaterialStatus: data.MaterialStatus,
		AutoExtracted:  data.AutoExtracted,
		AwardTime:      data.AwardAt,
		UpdateAT:       data.UpdatedAt.Unix(),
		CreateAT:       data.CreatedAt.Unix(),
		DeleteAT:       0,
	}
}
func buildEventList(data []*Event) []*model.Event {
	resp := make([]*model.Event, 0)
	for _, v := range data {
		s := buildEvent(v)
		resp = append(resp, s)
	}
	return resp
}
