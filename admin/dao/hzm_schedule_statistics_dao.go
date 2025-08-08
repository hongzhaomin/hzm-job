package dao

import (
	"errors"
	"github.com/hongzhaomin/hzm-job/admin/internal/global"
	"github.com/hongzhaomin/hzm-job/admin/po"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type HzmScheduleStatisticsDao struct{}

func (my *HzmScheduleStatisticsDao) FindByGeStartDay(startDay string) ([]*po.HzmScheduleStatistics, error) {
	if startDay == "" {
		return nil, nil
	}
	var result []*po.HzmScheduleStatistics
	err := global.SingletonPool().Mysql.
		Where("valid = 1 and day >= ?", startDay).
		Find(&result).
		Error
	return result, err
}

func (my *HzmScheduleStatisticsDao) FindByDayIfAbsentCreate(day string) (*po.HzmScheduleStatistics, error) {
	if day == "" {
		return nil, nil
	}

	var result *po.HzmScheduleStatistics
	err := global.SingletonPool().Mysql.Transaction(func(tx *gorm.DB) error {
		var sta po.HzmScheduleStatistics
		// SELECT * FROM `hzm_schedule_statistics` WHERE valid = 1 and day = ? ORDER BY `hzm_schedule_statistics`.`id` LIMIT 1 FOR UPDATE
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("valid = 1 and day = ?", day).
			First(&sta).
			Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 数据不存在，则创建任务日志
				d, _ := time.ParseInLocation(time.DateTime, day+" 00:00:00", time.Local)
				newSta := &po.HzmScheduleStatistics{
					Day: &d,
				}
				if err = tx.Select("Day").Create(newSta).Error; err != nil {
					return err
				}
				result = newSta
				return nil
			}
			return err
		}
		result = &sta
		return nil
	})
	return result, err
}

func (my *HzmScheduleStatisticsDao) Increment(id, total, success, fail int64) error {
	return global.SingletonPool().Mysql.
		Where("valid = 1 and id = ?", id).
		Model(&po.HzmScheduleStatistics{}).
		UpdateColumns(map[string]interface{}{
			"total_num":   gorm.Expr("total_num + ?", total),
			"success_num": gorm.Expr("success_num + ?", success),
			"fail_num":    gorm.Expr("fail_num + ?", fail),
		}).Error
}

func (my *HzmScheduleStatisticsDao) Update(sta *po.HzmScheduleStatistics) error {
	return global.SingletonPool().Mysql.
		Where("valid = 1 and id = ?", sta.Id).
		Model(&po.HzmScheduleStatistics{}).
		Updates(map[string]interface{}{
			"total_num":   sta.TotalNum,
			"success_num": sta.SuccessNum,
			"fail_num":    sta.FailNum,
		}).Error
}
