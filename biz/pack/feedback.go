package pack

import (
	"judgeMore/biz/dal/mysql"
	resp "judgeMore/biz/model/model"
)

func Feedback(feedback *mysql.Feedback) *resp.Feedback {
	return &resp.Feedback{
		Type:    feedback.UserID,
		Content: feedback.Content,
		UserID:  feedback.UserID,
	}
}
func FeedbackList(f []*mysql.Feedback) *resp.FeedbackList {
	r := make([]*resp.Feedback, 0)
	for _, v := range f {
		r = append(r, Feedback(v))
	}
	return &resp.FeedbackList{
		Item:  r,
		Total: int64(len(f)),
	}
}
