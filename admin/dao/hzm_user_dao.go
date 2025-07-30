package dao

import (
	"errors"
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/po"
	"github.com/hongzhaomin/hzm-job/admin/vo/req"
	"gorm.io/gorm"
	"strconv"
)

type HzmUserDao struct{}

func (my *HzmUserDao) Save(user *po.HzmUser) error {
	return global.SingletonPool().Mysql.
		Select("UserName", "Password", "Role", "Email").
		Create(user).
		Error
}

func (my *HzmUserDao) Update(user *po.HzmUser, upgradeVersion bool) error {
	return global.SingletonPool().Mysql.Transaction(func(tx *gorm.DB) error {
		err := tx.Save(user).Error
		if err != nil {
			return err
		}
		if upgradeVersion {
			return tx.Model(&po.HzmUser{}).
				Where("valid = 1 and id = ?", user.Id).
				Update("token_version", gorm.Expr("token_version + 1")).
				Error
		}
		return nil
	})
}

func (my *HzmUserDao) UpgradeTokenVersion(id int64) error {
	return global.SingletonPool().Mysql.
		Where("valid = 1 and id = ?", id).
		Model(&po.HzmUser{}).
		Update("token_version", gorm.Expr("token_version + 1")).
		Error
}

func (my *HzmUserDao) FindById(id *int64) (*po.HzmUser, error) {
	if id == nil {
		return nil, nil
	}
	var user po.HzmUser
	err := global.SingletonPool().Mysql.
		Where("valid = 1 and id = ?", id).
		First(&user).
		Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (my *HzmUserDao) FindByUserName(userName *string) (*po.HzmUser, error) {
	if userName == nil || *userName == "" {
		return nil, nil
	}
	var user po.HzmUser
	err := global.SingletonPool().Mysql.
		Where("valid = 1 and user_name = ?", userName).
		First(&user).
		Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (my *HzmUserDao) Page(param req.UserPage) (int64, []*po.HzmUser, error) {
	// 构造条件
	db := global.SingletonPool().Mysql
	db = db.Where("valid = ?", 1)
	if param.UserName != "" {
		db = db.Where("user_name LIKE ?", "%"+param.UserName+"%")
	}
	if param.Role != "" {
		role, _ := strconv.Atoi(param.Role)
		db = db.Where("role = ?", role)
	}

	var count int64
	db.Model(po.HzmUser{}).Count(&count)
	if count == 0 {
		return 0, nil, nil
	}

	var users []*po.HzmUser
	err := db.Offset(param.Start()).Limit(param.Limit()).Find(&users).Error
	return count, users, err
}

func (my *HzmUserDao) DeleteBatch(ids []int64) error {
	if len(ids) <= 0 {
		return nil
	}
	return global.SingletonPool().Mysql.
		Unscoped().
		Where("valid = 1 and id in (?)", ids).
		Delete(&po.HzmUser{}).
		Error
}
