package sipcron

// SimpleCron 极简表达式
// @yearly (or @annually)		Run once a year at midnight of 1 January						0 0 1 1 *
// @monthly						Run once a month at midnight of the first day of the month		0 0 1 * *
// @weekly						Run once a week at midnight on Sunday							0 0 * * 0
// @daily (or @midnight)		Run once a day at midnight										0 0 * * *
// @hourly						Run once an hour at the beginning of the hour					0 * * * *
type SimpleCron string

func (my *SimpleCron) GetDesc() string {
	return simpleCronMap[*my]
}

func Values() []SimpleCron {
	return []SimpleCron{
		Yearly,
		Monthly,
		Weekly,
		Daily,
		Hourly,
	}
}

const (
	Yearly  SimpleCron = "@yearly"  // 每年1月1日午夜运行一次
	Monthly SimpleCron = "@monthly" // 每月第一天午夜运行一次
	Weekly  SimpleCron = "@weekly"  // 每周在周日午夜运行一次
	Daily   SimpleCron = "@daily"   // 每天午夜运行一次
	Hourly  SimpleCron = "@hourly"  // 每小时运行一次
)

var simpleCronMap = map[SimpleCron]string{
	Yearly:  "每年1月1日午夜运行一次",
	Monthly: "每月第一天午夜运行一次",
	Weekly:  "每周在周日午夜运行一次",
	Daily:   "每天午夜运行一次",
	Hourly:  "每小时运行一次",
}
