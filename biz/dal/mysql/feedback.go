package mysql

import (
	"context"
	"judgeMore/pkg/constants"
	"judgeMore/pkg/errno"
)

func NewFeedback(ctx context.Context, feedback *Feedback) error {
	err := db.WithContext(ctx).Table(constants.TableFeedback).Create(&feedback).Error
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, "Create Feedback Error:"+err.Error())
	}
	return nil
}
func QueryFeedback(ctx context.Context) ([]*Feedback, error) {
	var result []*Feedback
	err := db.WithContext(ctx).Table(constants.TableFeedback).Find(&result).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "query Feedback Error:"+err.Error())
	}
	return result, nil
}
