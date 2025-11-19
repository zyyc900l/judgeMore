package service

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"judgeMore/biz/dal/mysql"
	"judgeMore/biz/service/model"
	"judgeMore/biz/service/taskqueue"
	"judgeMore/pkg/constants"
	"judgeMore/pkg/crypt"
	"judgeMore/pkg/errno"
	"judgeMore/pkg/utils"
	"strconv"
)

type MaintainService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewMaintainService(ctx context.Context, c *app.RequestContext) *MaintainService {
	return &MaintainService{
		ctx: ctx,
		c:   c,
	}
}

// 查找所有学院的信息
func (svc *MaintainService) QueryColleges(page_num, page_size int64) ([]*model.College, int64, error) {
	if utils.VerifyPageParam(page_num, page_size) {
		return nil, -1, errno.NewErrNo(errno.ParamVerifyErrorCode, "Page Param invalid")
	}
	collegeInfoList, err := QueryAllCollege(svc.ctx)
	var count int64
	count = int64(len(collegeInfoList))
	if err != nil {
		return nil, count, err
	}
	// 分页返回
	count = int64(len(collegeInfoList))
	startIndex := (page_num - 1) * page_size
	endIndex := startIndex + page_size
	if startIndex > count {
		return nil, 0, nil
	}
	if endIndex > count {
		endIndex = count
	}
	return collegeInfoList[startIndex:endIndex], count, nil
}

func (svc *MaintainService) QueryMajorByCollegeId(college_id int64, page_num, page_size int64) ([]*model.Major, int64, error) {
	if utils.VerifyPageParam(page_num, page_size) {
		return nil, -1, errno.NewErrNo(errno.ParamVerifyErrorCode, "Page Param invalid")
	}
	// 存在性检查
	exist, err := mysql.IsCollegeExist(svc.ctx, college_id)
	if err != nil {
		return nil, -1, fmt.Errorf("check college exist failed: %w", err)
	}
	if !exist {
		return nil, -1, errno.NewErrNo(errno.ServiceCollegeNotExistCode, "college not exist")
	}
	majorInfoList, count, err := mysql.GetMajorInfoByCollegeId(svc.ctx, college_id)
	if err != nil {
		return nil, count, err
	}
	// 分页返回
	count = int64(len(majorInfoList))
	startIndex := (page_num - 1) * page_size
	endIndex := startIndex + page_size
	if startIndex > count {
		return nil, 0, nil
	}
	if endIndex > count {
		endIndex = count
	}
	return majorInfoList[startIndex:endIndex], count, nil
}

func (svc *MaintainService) UploadMajor(major_name string, college_id int64) (int64, error) {
	// 检查
	exist, err := mysql.IsCollegeExist(svc.ctx, college_id)
	if err != nil {
		return -1, err
	}
	if !exist {
		return -1, errno.NewErrNo(errno.ServiceEventNotExistCode, "college not exist")
	}
	// 保存到数据库
	major_id, err := mysql.CreateMajor(svc.ctx, &model.Major{MajorName: major_name, CollegeId: college_id})
	if err != nil {
		return -1, err
	}
	// 返回数据库生成的自增ID
	taskqueue.AddUpdateCacheMajorTask(svc.ctx, constants.StructKey)
	return major_id, nil
}

func (svc *MaintainService) UploadCollege(collegeName string) (int64, error) {
	collegeId, err := mysql.CreateNewCollege(svc.ctx, collegeName)
	if err != nil {
		return -1, err
	}
	// 返回数据库生成的自增ID
	taskqueue.AddUpdateCacheCollegeTask(svc.ctx, constants.StructKey)
	return collegeId, nil
}

func (svc *MaintainService) AddUser(u *model.User) (string, error) {
	// 判断role
	if u.Role != "counselor" && u.Role != "admin" {
		return "", errno.NewErrNo(errno.ParamVerifyErrorCode, "role error :only admin or counselor")
	}
	u.Status = 1
	// 判断用户存在与否
	exist, err := mysql.IsUserExist(svc.ctx, &model.User{Uid: u.Uid})
	if err != nil {
		return "", err
	}
	if exist {
		return "", errno.NewErrNo(errno.ServiceUserExistCode, "user exist")
	}
	//
	u.Password, err = crypt.PasswordHash(u.Password)
	if err != nil {
		return "", fmt.Errorf("hash password error" + err.Error())
	}
	id, err := mysql.CreateUser(svc.ctx, u)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (svc *MaintainService) AddAdminRelation(relation *model.Relation) error {
	// 传进来的是学院与专业名称
	// 需要查找对应的ID
	if relation.CollegeName != "" {
		collegeList, err := QueryAllCollege(svc.ctx)
		if err != nil {
			return err
		}
		for _, collge := range collegeList {
			if collge.CollegeName == relation.CollegeName {
				relation.CollegeId = strconv.FormatInt(collge.CollegeId, 10)
				break
			}
		}
		if relation.CollegeId == "" { //查找后没找到
			return errno.NewErrNo(errno.ServiceCollegeNotExistCode, "College not exist")
		}
	} else {
		majorList, err := QueryAllMajor(svc.ctx)
		if err != nil {
			return err
		}
		for _, major := range majorList {
			if major.MajorName == relation.MajorName {
				relation.MajorId = strconv.FormatInt(major.MajorId, 10)
				break
			}
		}
		if relation.MajorId == "" { //查找后没找到
			return errno.NewErrNo(errno.ServiceCollegeNotExistCode, "Major not exist")
		}
	}
	// 这时已经正确匹配了存在的专业，插入数据库
	err := mysql.CreateNewRelation(svc.ctx, relation)
	if err != nil {
		return err
	}
	taskqueue.AddUpdateInsertStuTask(svc.ctx, constants.StructKey, relation)
	return nil
}

func (svc *MaintainService) QueryRecognizedReward(page_num, page_size int64) ([]*model.RecognizedEvent, int64, error) {
	if page_num <= 0 || page_size <= 0 {
		return nil, -1, errno.NewErrNo(errno.ParamVerifyErrorCode, "page param no invalid")
	}
	data, err := QueryAllRecognizedReward(svc.ctx)
	if err != nil {
		return nil, -1, err
	}
	count := int64(len(data))
	startIndex := (page_num - 1) * page_size
	endIndex := startIndex + page_size
	if startIndex > count {
		return nil, 0, nil
	}
	if endIndex > count {
		endIndex = count
	}
	return data[startIndex:endIndex], count, nil
}

func (svc *MaintainService) NewRecognizedEvent(re *model.RecognizedEvent) (*model.RecognizedEvent, error) {
	// 无法对新增信息做保证有效的存在性查询
	re.IsActive = true
	info, err := mysql.AddRecognizedEvent(svc.ctx, re)
	if err != nil {
		return nil, err
	}
	taskqueue.AddUpdateRecognizedTask(svc.ctx, constants.StructKey)
	taskqueue.AddUpdateElasticTask(svc.ctx, constants.REKey)
	return info, nil
}

func (svc *MaintainService) DeleteRecognizedEvent(id string) error {
	exist, err := IsRecognizedEventExist(svc.ctx, id)
	if err != nil {
		return err
	}
	if !exist {
		return errno.NewErrNo(errno.ServiceRecognizedNotExistCode, "Recognized Event Not Exist")
	}
	err = mysql.DeleteRecognized(svc.ctx, id)
	if err != nil {
		return err
	}
	taskqueue.AddUpdateRecognizedTask(svc.ctx, constants.StructKey)
	taskqueue.AddUpdateElasticTask(svc.ctx, constants.REKey)
	return nil
}
