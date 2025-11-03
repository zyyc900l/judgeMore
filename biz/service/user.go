package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"judgeMore/biz/dal/cache"
	"judgeMore/biz/dal/mysql"
	"judgeMore/biz/service/model"
	"judgeMore/pkg/constants"
	"judgeMore/pkg/crypt"
	"judgeMore/pkg/errno"
	"strings"
)

type UserService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewUserService(ctx context.Context, c *app.RequestContext) *UserService {
	return &UserService{
		ctx: ctx,
		c:   c,
	}
}

func (svc *UserService) Login(req *model.User) (*model.User, error) {
	exist, err := mysql.IsUserExist(svc.ctx, req)
	if err != nil {
		return nil, fmt.Errorf("check user exist failed: %w", err)
	}
	if !exist {
		return nil, errno.NewErrNo(errno.ServiceUserExistCode, "user not exist")
	}
	// 密码检验
	userInfo, err := mysql.GetUserInfoByRoleId(svc.ctx, req.Uid)
	if err != nil {
		return nil, fmt.Errorf("get user Info failed: %w", err)
	}
	// 激活检验
	if userInfo.Status == 0 {
		return nil, errno.NewErrNo(errno.ServiceUserDeathCode, "user not active ")
	}
	if !crypt.VerifyPassword(req.Password, userInfo.Password) {
		return nil, errno.Errorf(errno.ServiceUserPasswordError, "password not match")
	}
	return userInfo, nil
}

func (svc *UserService) Register(req *model.User) (string, error) {
	exist, err := mysql.IsUserExist(svc.ctx, req)
	if err != nil {
		return "", fmt.Errorf("check user exist failed: %w", err)
	}
	if exist {
		return "", errno.NewErrNo(errno.ServiceUserExistCode, "user already exist")
	}
	req.Password, err = crypt.PasswordHash(req.Password)
	if err != nil {
		return "", fmt.Errorf("hash password failed: %w", err)
	}
	//验证邮箱
	err = svc.SendEmail(svc.ctx, req)
	if err != nil {
		return "", fmt.Errorf("send email failed: %w", err)
	}
	// 创建账户
	uid, err := mysql.CreateUser(svc.ctx, req)
	if err != nil {
		return "", fmt.Errorf("create user failed: %w", err)
	}
	return uid, nil
}
func (svc *UserService) QueryUserInfo(u *model.User) (UserInfo *model.User, err error) {
	//存在性检验
	exist, err := mysql.IsUserExist(svc.ctx, u)
	if err != nil {
		return nil, fmt.Errorf("check user exist failed: %w", err)
	}
	if !exist {
		return nil, errno.NewErrNo(errno.ServiceUserExistCode, "user not exist")
	}
	userInfo, err := mysql.GetUserInfoByRoleId(svc.ctx, u.Uid)
	if err != nil {
		return nil, fmt.Errorf("get user Info failed: %w", err)
	}
	return userInfo, nil
}

func (svc *UserService) UpdateUserInfo(u *model.User) (UserInfo *model.User, err error) {
	// 由于uid读取自token 所以不做存在性检验
	id := GetUserIDFromContext(svc.c)
	u.Uid = id
	userInfo, err := svc.UpdateUser(svc.ctx, u)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

func (svc *UserService) SendEmail(ctx context.Context, user *model.User) error {
	// 首先进行验证 学号即Uid 与fzu邮箱强绑定
	Correct := strings.HasSuffix(user.Email, constants.EmailSuffix) && len(user.Email) == constants.EmailLength
	if !Correct {
		return errno.NewErrNo(errno.ServiceEmailIncorrectCode, "Uid do not match email")
	}
	key := fmt.Sprintf("Email:%s", user.Email)
	code, err := cache.PutCodeToCache(svc.ctx, key)
	if err != nil {
		return err
	}
	//err = utils.MailSendCode(user.Email, code)
	//if err != nil {
	//	return err
	//}
	hlog.Info(code)
	return nil
}
func (svc *UserService) UpdateUser(ctx context.Context, user *model.User) (*model.User, error) {
	// 由于这里需要对需更新的内容做选择 在svc处处理
	var updateParams []string

	if user.Major != "" {
		updateParams = append(updateParams, user.Major)
	}
	if user.College != "" {
		updateParams = append(updateParams, user.College)
	}
	if user.Grade != "" {
		updateParams = append(updateParams, user.Grade)
	}
	// 如果有需要更新的字段才执行
	if len(updateParams) > 0 {
		return mysql.UpdateInfoByRoleId(svc.ctx, user.Uid, updateParams...)
	}

	return nil, errno.Errorf(errno.InternalServiceErrorCode, "no element to update")

}

func (svc *UserService) VerifyEmail(data *model.EmailAuth) (err error) {
	// 判断存不存在
	key := fmt.Sprintf("Email:%s", data.Email)
	exist := cache.IsKeyExist(svc.ctx, key)
	if !exist {
		return errors.New("code expired")
	}
	emailParts := strings.Split(data.Email, "@")
	localPart := emailParts[0]
	data.Uid = localPart
	code, err := cache.GetCodeCache(svc.ctx, key)
	if err != nil {
		return err
	}
	if code != data.Code {
		return errors.New("code not match")
	}
	// 更新user表的信息
	err = mysql.ActivateUser(svc.ctx, data.Uid)
	if err != nil {
		return err
	}
	err = cache.DeleteCodeCache(svc.ctx, key)
	if err != nil {
		return err
	}
	return nil
}
