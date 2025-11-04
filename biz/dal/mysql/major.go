package mysql

import (
	"context"
	"judgeMore/biz/service/model"
	"judgeMore/pkg/constants"
	"judgeMore/pkg/errno"
)

type GetMajorInfoByCollegeIdFunc func(ctx context.Context, college_id string) ([]*model.Major, int64, error)
type CreateMajorFunc func(ctx context.Context, major *Major) error

// 对外暴露的函数变量（默认指向真实实现）
var (
	GetMajorInfoByCollegeId GetMajorInfoByCollegeIdFunc = RealGetMajorInfoByCollegeId
	CreateMajor             CreateMajorFunc             = RealCreateMajor
)

func RealGetMajorInfoByCollegeId(ctx context.Context, college_id string) ([]*model.Major, int64, error) {
	var majorInfos []*Major
	var count int64
	err := db.WithContext(ctx).
		Table(constants.TableMajor).
		Where("college_id = ?", college_id).
		Find(&majorInfos).
		Count(&count).
		Error
	if err != nil {
		return nil, -1, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed query stu event: %v", err)
	}
	return BuildMajorInfoList(majorInfos), count, err
}

func BuildMajorInfo(data *Major) *model.Major {
	return &model.Major{
		MajorId:   data.MajorId,
		MajorName: data.MajorName,
		CollegeId: data.CollegeId,
	}
}
func BuildMajorInfoList(data []*Major) []*model.Major {
	resp := make([]*model.Major, 0)
	for _, v := range data {
		s := BuildMajorInfo(v)
		resp = append(resp, s)
	}
	return resp
}

func RealCreateMajor(ctx context.Context, major *Major) error {
	err := db.WithContext(ctx).
		Table(constants.TableMajor).
		Create(major).
		Error
	if err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to create major: %v", err)
	}
	// gorm会自动将数据库生成的自增ID回填到major对象中
	return nil
}
