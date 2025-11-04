package mysql

import (
	"context"
	"errors"
	"judgeMore/biz/service/model"
	"judgeMore/pkg/constants"
	"judgeMore/pkg/errno"

	"gorm.io/gorm"
)

type GetCollegeInfoFunc func(ctx context.Context) ([]*model.College, int64, error)
type IsCollegeExistFunc func(ctx context.Context, college_id string) (bool, error)

// 对外暴露的函数变量（默认指向真实实现,用于测试）
var (
	GetCollegeInfo GetCollegeInfoFunc = RealGetCollegeInfo
	IsCollegeExist IsCollegeExistFunc = RealIsCollegeExist
)

func RealGetCollegeInfo(ctx context.Context) ([]*model.College, int64, error) {
	var collegeInfos []*College
	var count int64
	err := db.WithContext(ctx).
		Table(constants.TableCollege).
		Find(&collegeInfos).
		Count(&count).
		Error
	if err != nil {
		return nil, -1, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed query college: %v", err)
	}
	return buildCollegeInfoList(collegeInfos), count, err
}

func buildCollegeInfo(data *College) *model.College {
	return &model.College{
		CollegeId:   data.CollegeId,
		CollegeName: data.CollegeName,
	}
}
func buildCollegeInfoList(data []*College) []*model.College {
	resp := make([]*model.College, 0)
	for _, v := range data {
		s := buildCollegeInfo(v)
		resp = append(resp, s)
	}
	return resp
}
func RealIsCollegeExist(ctx context.Context, college_id string) (bool, error) {
	var collegeInfo *College
	err := db.WithContext(ctx).
		Table(constants.TableCollege).
		Where("college_id = ?", college_id).
		First(&collegeInfo).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { //没找到了用户不存在
			return false, nil
		}
		return false, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query college: %v", err)
	}
	return true, nil
}
