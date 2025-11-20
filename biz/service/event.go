package service

import (
	"context"
	"fmt"
	"github.com/antlabs/strsim"
	"github.com/cloudwego/hertz/pkg/app"
	"judgeMore/biz/dal/mysql"
	"judgeMore/biz/service/model"
	"judgeMore/biz/service/taskqueue"
	"judgeMore/config"
	"judgeMore/pkg/constants"
	"judgeMore/pkg/errno"
	"judgeMore/pkg/oss"
	"judgeMore/pkg/utils"
	"mime/multipart"
	"path/filepath"
	"time"
)

type EventService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewEventService(ctx context.Context, c *app.RequestContext) *EventService {
	return &EventService{
		ctx: ctx,
		c:   c,
	}
}
func (svc *EventService) QueryEventByEventId(event_id string) (*model.Event, error) {
	exist, err := mysql.IsEventExist(svc.ctx, event_id)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errno.NewErrNo(errno.ServiceEventNotExistCode, "event not exist")
	}
	eventInfo, err := mysql.GetEventInfoById(svc.ctx, event_id)
	if err != nil {
		return nil, err
	}
	return eventInfo, nil
}

func (svc *EventService) QueryEventByStuId() ([]*model.Event, int64, error) {
	stu_id := GetUserIDFromContext(svc.c)
	exist, err := mysql.IsUserExist(svc.ctx, &model.User{Uid: stu_id})
	if err != nil {
		return nil, -1, err
	}
	if !exist {
		return nil, -1, errno.NewErrNo(errno.ServiceUserNotExistCode, "user not exist")
	}
	eventInfoList, count, err := mysql.GetEventInfoByStuId(svc.ctx, stu_id)
	if err != nil {
		return nil, count, err
	}
	return eventInfoList, count, nil
}
func (svc *EventService) UpdateEventStatus(event_id string, status int64) (*model.Event, error) {
	// 常规检验
	admin_id := GetUserIDFromContext(svc.c)
	exist, err := mysql.IsEventExist(svc.ctx, event_id)
	if err != nil {
		return nil, fmt.Errorf("check event exist failed: %w", err)
	}
	if !exist {
		return nil, errno.NewErrNo(errno.ServiceEventNotExistCode, "event not exist")
	}
	// 检验传上来的status
	if status != 1 && status != 2 {
		return nil, errno.NewErrNo(errno.ParamVerifyErrorCode, "status should be 1 or 2")
	}
	// 判断是否有权限审核
	eventInfo, err := mysql.GetEventInfoById(svc.ctx, event_id)
	if err != nil {
		return nil, err
	}
	if eventInfo.MaterialStatus == "已审核" || eventInfo.MaterialStatus == "驳回" {
		return nil, errno.NewErrNo(errno.ServiceRepeatAction, "the martial have been checked")
	}
	exist, err = mysql.IsAdminRelationExist(svc.ctx, admin_id, eventInfo.Uid)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errno.NewErrNo(errno.ServiceNoAuthToDo, "No permission to check the stu's event")
	}
	info, err := mysql.UpdateEventStatus(svc.ctx, event_id, status)
	if err != nil {
		return nil, err
	}
	// 将计算积分加入任务队列，异步处理
	if status == 1 {
		taskqueue.AddScoreTask(svc.ctx, constants.EventKey, event_id)
	}
	return info, nil
}

func (svc *EventService) UploadEventFile(file *multipart.FileHeader) (string, error) {
	stu_id := GetUserIDFromContext(svc.c)
	// 检测文件类型
	err := oss.IsImage(file)
	if err != nil {
		return "", fmt.Errorf("check image failed: %w", err)
	}
	// 识别图片信息
	// 暂存本地
	fileName := fmt.Sprintf("%v_%v", stu_id, time.Now().Unix())
	err = oss.SaveFile(file, constants.StorePath, fileName)
	if err != nil {
		return "", fmt.Errorf("save file failed: %w", err)
	}
	filePath := filepath.Join(constants.StorePath, fileName)
	Info, err := utils.CallGLM4VWithImage(svc.ctx, filePath, config.OpenAI.ApiKey)
	if err != nil {
		return "", fmt.Errorf("call glm4v with image failed: %w", err)
	}
	if Info.Success == constants.AIErrorMessage {
		return "", errno.NewErrNo(errno.ServiceImageNotAwardCode, "image not a award or certificate")
	}
	eventInfo := &model.Event{
		Uid:            stu_id,
		EventName:      Info.EventName,
		EventOrganizer: Info.EventSponsor,
		AwardContent:   Info.AwardLevel,
		AwardTime:      Info.EventTime,
		AutoExtracted:  true,
	}
	// 材料云端存储
	url, err := oss.Upload(filePath, fileName, stu_id, constants.OssOrigin)
	if err != nil {
		return "", fmt.Errorf("upload file failed: %w", err)
	}
	eventInfo.MaterialUrl = url
	// 存入数据库 不进行check则无法完善相关材料。但不匹配也应该存下该材料
	err = CheckEvent(svc.ctx, eventInfo)
	if err != nil {
		return "", err
	}
	event_id, err := mysql.CreateNewEvent(svc.ctx, eventInfo)
	if err != nil {
		return "", err
	}
	return event_id, nil
}

func (svc *EventService) UpdateEventLevel(event_id string, level string, appeal_id string) error {
	// 材料存在
	exist, err := mysql.IsEventExist(svc.ctx, event_id)
	if err != nil {
		return err
	}
	if !exist {
		return errno.NewErrNo(errno.ServiceEventNotExistCode, "event not exist")
	}
	// 材料是经过申诉
	exist, err = mysql.IsAppealExistByAppealId(svc.ctx, appeal_id)
	if err != nil {
		return err
	}
	if !exist {
		return errno.NewErrNo(errno.ServiceEventUnChangedCode, "appeal not exist cannot change Event")
	}
	// 联查的话要查很远 但目前这样没办法判断申诉是否已完成 后期加上权限的问题会爆发一系列问题
	recordInfo, err := mysql.QueryScoreRecordByEventId(svc.ctx, event_id)
	if err != nil {
		return err
	}
	if recordInfo.AppealId != appeal_id {
		return errno.NewErrNo(errno.ServiceEventUnChangedCode, "appeal not match the event")
	}
	user_id := GetUserIDFromContext(svc.c)
	exist, err = mysql.IsAdminRelationExist(svc.ctx, user_id, recordInfo.UserId)
	if err != nil {
		return err
	}
	if !exist {
		return errno.NewErrNo(errno.ServiceNoAuthToDo, "No permission to update the stu's appeal")
	}
	err = mysql.UpdateEventLevel(svc.ctx, event_id, level)
	if err != nil {
		return err
	}
	// 再次异步计算
	taskqueue.AddSyncScoreTask(svc.ctx, constants.EventKey, event_id)
	return nil
}

func CheckEvent(ctx context.Context, eventInfo *model.Event) error {
	req := &model.ViewRecognizedRewardReq{
		EventName:     &eventInfo.EventName,
		OrganizerName: &eventInfo.EventOrganizer,
	}
	Event, err := SearchRecognizedEvent(ctx, req)
	if err != nil {
		return err
	}
	// 设置最低相似度阈值
	const minSimilarity = 0.5
	var bestMatch *model.RecognizedEvent
	var highestSimilarity float64
	// 遍历所有事件，计算相似度
	for _, v := range Event {
		similarity := strsim.Compare(v.RecognizedEventName, eventInfo.EventName)
		// 如果相似度高于阈值且是当前最高相似度
		if similarity >= minSimilarity && similarity > highestSimilarity {
			highestSimilarity = similarity
			bestMatch = v
		}
	}
	// 如果找到了符合条件的匹配
	if bestMatch != nil {
		eventInfo.RecognizeId = bestMatch.RecognizedEventId
		eventInfo.EventLevel = bestMatch.RecognizedLevel                     //直接根据认定赛事表来确定，可以不用做匹配
		eventInfo.AwardLevel = utils.AppraisalReward(eventInfo.AwardContent) //这里再做一次模糊鉴定
	} else {
		return errno.NewErrNo(errno.ServiceEventNotMatchCode, "reward not match")
	}
	return nil
}

func (svc *EventService) QueryBelongStuEvent(status string) ([]*model.Event, int64, error) {
	// 这边由token提取 前面jwt中间件会将学生token拦在外面 保证权限够高
	if status != "待审核" && status != "已审核" && status != "驳回" {
		return nil, 0, errno.NewErrNo(errno.InternalDatabaseErrorCode, "error status type")
	}
	user_id := GetUserIDFromContext(svc.c)
	stuList, err := mysql.QueryStuByAdmin(svc.ctx, user_id)
	if err != nil {
		return nil, -1, err
	}
	if len(stuList) == 0 {
		return nil, 0, nil
	}
	eventInfoList := make([]*model.Event, 0)
	var totalCount int64 = 0

	for _, v := range stuList {
		events, _, err := mysql.GetEventInfoByStuId(svc.ctx, v) // events 是 []*model.Event
		if err != nil {
			return nil, -1, err
		}
		// 过滤符合状态的事件
		for _, event := range events {
			if event.MaterialStatus == status {
				eventInfoList = append(eventInfoList, event)
				totalCount++
			}
		}
	}
	return eventInfoList, totalCount, nil
}
