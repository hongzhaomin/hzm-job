package service

import (
	"errors"
	"fmt"
	"github.com/hongzhaomin/hzm-job/admin/dao"
	"github.com/hongzhaomin/hzm-job/admin/internal/consts"
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/internal/tool"
	"github.com/hongzhaomin/hzm-job/admin/po"
	"github.com/hongzhaomin/hzm-job/admin/vo"
	"github.com/hongzhaomin/hzm-job/admin/vo/req"
)

type HzmUserService struct {
	hzmUserDao               dao.HzmUserDao
	hzmExecutorDao           dao.HzmExecutorDao
	hzmUserDataPermissionDao dao.HzmUserDataPermissionDao
}

func (my *HzmUserService) Add(param req.User) error {
	oldUser, err := my.hzmUserDao.FindByUserName(param.UserName)
	if err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return consts.ServerError
	}
	if oldUser != nil {
		return errors.New(fmt.Sprintf("账户[%s]已存在", *param.UserName))
	}

	// 新增用户
	password := tool.MD5(*param.Password)
	user := po.HzmUser{
		UserName: param.UserName,
		Password: &password,
		Role:     (*byte)(param.Role),
		Email:    param.Email,
	}
	if err = my.hzmUserDao.Save(&user); err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return consts.ServerError
	}

	// 更新用户数据权限
	dataPerms := param.DataPerms
	if po.CommonUser.Is(user.Role) && len(dataPerms) > 0 {
		if err = my.updateUserDataPerms(user.Id, dataPerms); err != nil {
			global.SingletonPool().Log.Error(err.Error())
			return consts.ServerError
		}
	}
	return nil
}

func (my *HzmUserService) Edit(param req.User) error {
	userId := param.Id
	user, err := my.hzmUserDao.FindById(userId)
	if err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return consts.ServerError
	}
	if user == nil {
		return errors.New("账户不存在")
	}

	// 校验同名用户是否存在
	sameNameUser, err := my.hzmUserDao.FindByUserName(param.UserName)
	if err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return consts.ServerError
	}
	if sameNameUser != nil && *sameNameUser.Id != *userId {
		return errors.New(fmt.Sprintf("账户[%s]已存在", *param.UserName))
	}

	upgradeVersion := false
	// 编辑用户
	user.UserName = param.UserName
	if !param.Role.Is(user.Role) {
		// 角色变更，升级版本
		upgradeVersion = true
	}
	user.Role = (*byte)(param.Role)
	user.Email = param.Email
	// 密码为空则不更新
	if param.Password != nil && *param.Password != "" {
		newPassword := tool.MD5(*param.Password)
		if *user.Password != newPassword {
			user.Password = &newPassword
			// 密码更新，升级token版本
			upgradeVersion = true
		}
	}
	if err = my.hzmUserDao.Update(user, upgradeVersion); err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return consts.ServerError
	}

	// 更新用户数据权限
	dataPerms := param.DataPerms
	if po.CommonUser.Is(user.Role) && len(dataPerms) > 0 {
		if err = my.updateUserDataPerms(user.Id, dataPerms); err != nil {
			global.SingletonPool().Log.Error(err.Error())
			return consts.ServerError
		}
	} else {
		// 也许从普通用户改成管理员了，所以要把之前普通用户的权限数据删除掉
		if err = my.hzmUserDataPermissionDao.DeleteByUserIds(*userId); err != nil {
			global.SingletonPool().Log.Error(err.Error())
			return consts.ServerError
		}
	}
	return nil
}

func (my *HzmUserService) updateUserDataPerms(userId *int64, dataPerms []*vo.DataPermsTransfer) error {
	// 删除以前的权限数据
	if err := my.hzmUserDataPermissionDao.DeleteByUserIds(*userId); err != nil {
		return err
	}

	// 新增新配置的权限数据
	newDataPerms := tool.BeanConv[vo.DataPermsTransfer, po.HzmUserDataPermission](dataPerms,
		func(dataPerm *vo.DataPermsTransfer) (*po.HzmUserDataPermission, bool) {
			return &po.HzmUserDataPermission{
				UserId:     userId,
				ExecutorId: &dataPerm.Value,
			}, true
		})
	if err := my.hzmUserDataPermissionDao.SaveBatch(newDataPerms); err != nil {
		return err
	}
	return nil
}

func (my *HzmUserService) DeleteBatch(userIds []int64) error {
	if err := my.hzmUserDao.DeleteBatch(userIds); err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return consts.ServerError
	}
	// 删除用户的权限数据
	if err := my.hzmUserDataPermissionDao.DeleteByUserIds(userIds...); err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return consts.ServerError
	}
	return nil
}

func (my *HzmUserService) PageUsers(param req.UserPage) (int64, []*vo.User) {
	count, users, err := my.hzmUserDao.Page(param)
	if err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return 0, nil
	}
	voUsers := tool.BeanConv[po.HzmUser, vo.User](users, func(user *po.HzmUser) (*vo.User, bool) {
		return &vo.User{
			Id:       user.Id,
			UserName: user.UserName,
			Role:     (*po.UserRole)(user.Role),
			RoleName: po.GetNameByRole(user.Role),
			Email:    user.Email,
		}, true
	})
	return count, voUsers
}

func (my *HzmUserService) DataPermsAll() []*vo.DataPermsTransfer {
	executors, err := my.hzmExecutorDao.FindAll()
	if err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return nil
	}
	return tool.BeanConv[po.HzmExecutor, vo.DataPermsTransfer](executors,
		func(executor *po.HzmExecutor) (*vo.DataPermsTransfer, bool) {
			return &vo.DataPermsTransfer{
				Value:    *executor.Id,
				Title:    fmt.Sprintf("%s[%s]", *executor.Name, *executor.AppKey),
				Disabled: false,
				Checked:  false,
			}, true
		})
}

func (my *HzmUserService) DataPermsByUserId(userId int64) *vo.UserDataPerms {
	// 所有权限数据
	allDataPerms := my.DataPermsAll()
	// 查询用户以配置的权限数据
	executorIds, err := my.hzmUserDataPermissionDao.FindExecutorIdsByUserId(userId)
	if err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return nil
	}
	executors, err := my.hzmExecutorDao.FindByIds(executorIds)
	if err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return nil
	}
	selectedDataPerms := tool.BeanConv[po.HzmExecutor, vo.DataPermsTransfer](executors,
		func(executor *po.HzmExecutor) (*vo.DataPermsTransfer, bool) {
			return &vo.DataPermsTransfer{
				Value:    *executor.Id,
				Title:    fmt.Sprintf("%s[%s]", *executor.Name, *executor.AppKey),
				Disabled: false,
				Checked:  false,
			}, true
		})
	return &vo.UserDataPerms{
		AllDataPerms:      allDataPerms,
		SelectedDataPerms: selectedDataPerms,
	}
}

func (my *HzmUserService) EditPassword(userId int64, param req.EditPasswordParam) error {
	if param.NewPassword != param.AgainPassword {
		return errors.New("两次输入的新密码不一致")
	}
	user, err := my.hzmUserDao.FindById(&userId)
	if err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return consts.ServerError
	}
	if user == nil {
		return errors.New("当前用户不存在")
	}
	if *user.Password != tool.MD5(param.OldPassword) {
		return errors.New("密码输入错误")
	}
	newPassword := tool.MD5(param.NewPassword)
	user.Password = &newPassword
	if err = my.hzmUserDao.Update(user, *user.Password != newPassword); err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return consts.ServerError
	}
	return nil
}

func (my *HzmUserService) Login(param req.Login) (*vo.LoginUser, error) {
	user, err := my.hzmUserDao.FindByUserName(&param.UserName)
	if user == nil {
		if err != nil {
			global.SingletonPool().Log.Error(err.Error())
			return nil, consts.ServerError
		}
		return nil, errors.New("用户不存在")
	}
	if *user.Password != tool.MD5(param.Password) {
		return nil, errors.New("密码不正确")
	}

	token, err := tool.GenerateToken(*user.Id, *user.TokenVersion)
	if err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return nil, consts.ServerError
	}

	return &vo.LoginUser{
		Id:          user.Id,
		UserName:    user.UserName,
		AccessToken: &token,
	}, nil
}

func (my *HzmUserService) FindUserById(userId int64) *vo.User {
	user, err := my.hzmUserDao.FindById(&userId)
	if user == nil {
		if err != nil {
			global.SingletonPool().Log.Error(err.Error())
		}
		return nil
	}
	return &vo.User{
		Id:       user.Id,
		UserName: user.UserName,
		Role:     (*po.UserRole)(user.Role),
		RoleName: po.GetNameByRole(user.Role),
		Email:    user.Email,
	}
}

func (my *HzmUserService) LoginOut(userId int64) error {
	// 升级token版本
	if err := my.hzmUserDao.UpgradeTokenVersion(userId); err != nil {
		global.SingletonPool().Log.Error(err.Error())
		return consts.ServerError
	}
	return nil
}
