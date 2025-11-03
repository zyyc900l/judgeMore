package service

import (
	"context"
	"fmt"
	"judgeMore/biz/dal/mysql"
	"judgeMore/biz/service/model"
	"judgeMore/pkg/errno"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

type AdminService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewAdminService(ctx context.Context, c *app.RequestContext) *AdminService {
	return &AdminService{
		ctx: ctx,
		c:   c,
	}
}

func (svc *AdminService) QueryColleges() ([]*model.College, int64, error) {

	collegeInfoList, count, err := mysql.GetCollegeInfo(svc.ctx)
	if err != nil {
		return nil, count, fmt.Errorf("get event Info failed: %w", err)
	}
	return collegeInfoList, count, nil
}

func (svc *AdminService) QueryMajorByCollegeId(college_id string) ([]*model.Major, int64, error) {
	//检查
	exist, err := mysql.IsCollegeExist(svc.ctx, college_id)
	if err != nil {
		return nil, int64(-1), fmt.Errorf("check college exist failed: %w", err)
	}
	if !exist {
		return nil, int64(-1), errno.NewErrNo(errno.ServiceEventExistCode, "college not exist")
	}
	majorInfoList, count, err := mysql.GetMajorInfoByCollegeId(svc.ctx, college_id)
	if err != nil {
		return nil, count, fmt.Errorf("get major Info failed: %w", err)
	}
	return majorInfoList, count, nil
}

func (svc *AdminService) UploadMajor(major_name string, college_id string) (string, error) {
	// 检查
	exist, err := mysql.IsCollegeExist(svc.ctx, college_id)
	if err != nil {
		return "", fmt.Errorf("check college exist failed: %w", err)
	}
	if !exist {
		return "", errno.NewErrNo(errno.ServiceEventExistCode, "college not exist")
	}

	// 构造数据库实体
	major := &mysql.Major{
		MajorName: major_name,
		CollegeId: college_id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	// 保存到数据库
	err = mysql.CreateMajor(svc.ctx, major)
	if err != nil {
		return "", fmt.Errorf("create major failed: %w", err)
	}
	// 返回数据库生成的自增ID
	return major.MajorId, nil
}
