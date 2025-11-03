package service

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"judgeMore/biz/dal/mysql"
	"judgeMore/biz/service/model"
	"judgeMore/pkg/errno"
)

type ScoreService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewScoreService(ctx context.Context, c *app.RequestContext) *ScoreService {
	return &ScoreService{
		ctx: ctx,
		c:   c,
	}
}

func (svc *ScoreService) QueryScoreRecordByScoreId(score_id string) (*model.ScoreRecord, error) {
	exist, err := mysql.IsScoreRecordExist(svc.ctx, score_id)
	if err != nil {
		return nil, fmt.Errorf("check score exist failed: %w", err)
	}
	if !exist {
		return nil, errno.NewErrNo(errno.ServiceEventExistCode, "event not exist")
	}
	recordInfo, err := mysql.QueryScoreRecordByScoreId(svc.ctx, score_id)
	if err != nil {
		return nil, fmt.Errorf("get record Info failed: %w", err)
	}
	return recordInfo, nil
}

func (svc *ScoreService) QueryScoreRecordByEventId(event_id string) (*model.ScoreRecord, error) {
	exist, err := mysql.IsScoreRecordExist_Event(svc.ctx, event_id)
	if err != nil {
		return nil, fmt.Errorf("check score exist failed: %w", err)
	}
	if !exist {
		return nil, errno.NewErrNo(errno.ServiceEventExistCode, "event not exist")
	}
	recordInfo, err := mysql.QueryScoreRecordByEventId(svc.ctx, event_id)
	if err != nil {
		return nil, fmt.Errorf("get record Info failed: %w", err)
	}
	return recordInfo, nil
}
func (svc *ScoreService) QueryScoreRecordByStuId(stu_id string) ([]*model.ScoreRecord, int64, error) {
	exist, err := mysql.IsUserExist(svc.ctx, &model.User{Uid: stu_id})
	if err != nil {
		return nil, -1, fmt.Errorf("check user exist failed: %w", err)
	}
	if !exist {
		return nil, -1, errno.NewErrNo(errno.ServiceUserExistCode, "user not exist")
	}
	recordInfoList, count, err := mysql.QueryScoreRecordByStuId(svc.ctx, stu_id)
	if err != nil {
		return nil, count, fmt.Errorf("get record Info failed: %w", err)
	}
	return recordInfoList, count, nil
}
