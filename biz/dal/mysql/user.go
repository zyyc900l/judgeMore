package mysql

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"judgeMore/biz/service/model"
	"judgeMore/pkg/constants"
	"judgeMore/pkg/errno"
)

func IsUserExist(ctx context.Context, user *model.User) (bool, error) {
	var userInfo *User
	err := db.WithContext(ctx).
		Table(constants.TableUser).
		Where("role_id = ?", user.Uid).
		First(&userInfo).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { //没找到了说明用户不存在
			return false, nil
		}
		return false, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to query user: %v", err)
	}
	return true, nil
}
func CreateUser(ctx context.Context, user *model.User) (string, error) {
	userInfo := &User{
		UserName: user.UserName,
		Password: user.Password,
		Email:    user.Email,
		RoleId:   user.Uid,
		UserRole: user.Role,
		Status:   user.Status, //初始状态未激活
	}
	err := db.WithContext(ctx).
		Table(constants.TableUser).
		Create(userInfo).
		Error
	if err != nil {
		return "", errno.NewErrNo(errno.InternalDatabaseErrorCode, "Create User Error:"+err.Error())
	}
	return userInfo.RoleId, nil
}

// 该函数调用前检验存在性
func GetUserInfoByRoleId(ctx context.Context, role_id string) (*model.User, error) {
	var userInfo *User
	err := db.WithContext(ctx).
		Table(constants.TableUser).
		Where("role_id = ?", role_id).
		First(&userInfo).
		Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "query user Info error:"+err.Error())
	}
	return &model.User{
		Uid:      userInfo.RoleId,
		UserName: userInfo.UserName,
		Grade:    userInfo.Grade,
		Major:    userInfo.Major,
		College:  userInfo.College,
		Password: userInfo.Password,
		Status:   userInfo.Status,
		Email:    userInfo.Email,
		Role:     userInfo.UserRole,
		UpdateAT: userInfo.UpdatedAt.Unix(),
		CreateAT: userInfo.CreatedAt.Unix(),
		DeleteAT: 0,
	}, nil
}
func UpdateUserPassword(ctx context.Context, user *model.User) error {
	err := db.WithContext(ctx).
		Table(constants.TableUser).
		Where("role_id = ?", user.Uid).
		Update("password", user.Password).
		Error
	if err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "Update User Password"+err.Error())
	}
	return nil
}

func UpdateInfoByRoleId(ctx context.Context, role_id string, element ...string) (*model.User, error) {
	updateFields := make(map[string]interface{})
	for i, value := range element {
		if value == "" {
			continue // 跳过空值
		}
		switch i {
		case 0:
			updateFields["major"] = value
		case 1:
			updateFields["college"] = value
		case 2:
			updateFields["grade"] = value
		}
	}
	// 存在多值更新 以事务提交保证原子性
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Table(constants.TableUser).
			Where("role_id = ?", role_id).
			Updates(updateFields).
			Error
	})
	if err != nil {
		return nil, errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to update userInfo: %v", err)
	}

	return GetUserInfoByRoleId(ctx, role_id)
}
func QueryUserByUserName(ctx context.Context, name string) ([]*model.User, error) {
	var user []*User
	err := db.WithContext(ctx).
		Table(constants.TableUser).Where("user_name = ?", name).
		Find(&user).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "query userinfo by username error: "+err.Error())
	}
	return buildUserList(user), nil
}
func QueryUserByCollege(ctx context.Context, college string) ([]*model.User, error) {
	var user []*User
	err := db.WithContext(ctx).
		Table(constants.TableUser).Where("college = ?", college).
		Find(&user).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "query userinfo by username error: "+err.Error())
	}
	return buildUserList(user), nil
}
func QueryUserByRole(ctx context.Context, role string) ([]*model.User, error) {
	var user []*User
	err := db.WithContext(ctx).
		Table(constants.TableUser).Where("user_role = ?", role).
		Find(&user).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "query userinfo by userrole error: "+err.Error())
	}
	return buildUserList(user), nil
}
func QueryUserByMajor(ctx context.Context, major string) ([]*model.User, error) {
	var user []*User
	err := db.WithContext(ctx).
		Table(constants.TableUser).Where("major = ?", major).
		Find(&user).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "query userinfo by major error: "+err.Error())
	}
	return buildUserList(user), nil
}
func QueryUserByUserGrade(ctx context.Context, grade string) ([]*model.User, error) {
	var user []*User
	err := db.WithContext(ctx).
		Table(constants.TableUser).Where("grade = ?", grade).
		Find(&user).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "query userinfo by grade error: "+err.Error())
	}
	return buildUserList(user), nil
}
func QueryAllStu(ctx context.Context) ([]*model.User, error) {
	var user []*User
	err := db.WithContext(ctx).
		Table(constants.TableUser).Where("user_role = ?", "student").
		Find(&user).Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "query userinfo by major error: "+err.Error())
	}
	return buildUserList(user), nil
}
func QueryUserIdByMajor(ctx context.Context, major, grade string) ([]string, error) {
	var stuIds []string
	err := db.WithContext(ctx).
		Table(constants.TableUser).
		Where("major = ? and grade = ?", major, grade).
		Pluck("role_id", &stuIds).
		Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "query user_id error "+err.Error())
	}
	return stuIds, nil
}
func QueryUserIdByCollege(ctx context.Context, college string) ([]string, error) {
	var stuIds []string
	err := db.WithContext(ctx).
		Table(constants.TableUser).
		Where("college", college).
		Pluck("role_id", &stuIds).
		Error
	if err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, "query user_id error "+err.Error())
	}
	return stuIds, nil
}
func ActivateUser(ctx context.Context, uid string) error {
	err := db.WithContext(ctx).
		Table(constants.TableUser).
		Where("role_id = ?", uid).
		Update("status", 1).
		Error
	if err != nil {
		return errno.Errorf(errno.InternalDatabaseErrorCode, "mysql: failed to activate user: %v", err)
	}
	return nil
}

func buildUserInfo(userInfo *User) *model.User {
	return &model.User{
		Uid:      userInfo.RoleId,
		UserName: userInfo.UserName,
		Grade:    userInfo.Grade,
		Major:    userInfo.Major,
		College:  userInfo.College,
		Password: userInfo.Password,
		Status:   userInfo.Status,
		Email:    userInfo.Email,
		Role:     userInfo.UserRole,
		UpdateAT: userInfo.UpdatedAt.Unix(),
		CreateAT: userInfo.CreatedAt.Unix(),
		DeleteAT: 0,
	}
}
func buildUserList(userInfo []*User) []*model.User {
	re := make([]*model.User, 0)
	for _, v := range userInfo {
		re = append(re, buildUserInfo(v))
	}
	return re
}
