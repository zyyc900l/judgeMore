package service

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"judgeMore/biz/dal/mysql"
)

type FeedbackService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewFeedbackService(ctx context.Context, c *app.RequestContext) *FeedbackService {
	return &FeedbackService{
		ctx: ctx,
		c:   c,
	}
}

func (svc *FeedbackService) NewFeedback(t, content string) error {
	user_id := GetUserIDFromContext(svc.c)
	feedback := &mysql.Feedback{
		UserID:  user_id,
		Type:    t,
		Content: content,
	}
	return mysql.NewFeedback(svc.ctx, feedback)
}
func (svc *FeedbackService) QueryFeedback() ([]*mysql.Feedback, error) {
	return mysql.QueryFeedback(svc.ctx)
}
