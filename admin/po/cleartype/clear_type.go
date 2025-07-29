package cleartype

import "time"

type ClearType int

func (my *ClearType) GetDesc() string {
	return clearTypeMap[*my]
}

func (my *ClearType) GetQueryParams() (*time.Time, *int) {
	switch *my {
	case OneMonthAgo:
		timeBefore := time.Now().Add(-time.Hour * 24 * 30)
		return &timeBefore, nil
	case ThreeMonthAgo:
		timeBefore := time.Now().Add(-time.Hour * 24 * 30 * 3)
		return &timeBefore, nil
	case SixMonthAgo:
		timeBefore := time.Now().Add(-time.Hour * 24 * 30 * 6)
		return &timeBefore, nil
	case OneYearAgo:
		timeBefore := time.Now().Add(-time.Hour * 24 * 30 * 12)
		return &timeBefore, nil
	case OneThousandAgo:
		count := 1000
		return nil, &count
	case TenThousandAgo:
		count := 10000
		return nil, &count
	case ThirtyThousandAgo:
		count := 30000
		return nil, &count
	case OneHundredThousandAgo:
		count := 100000
		return nil, &count
	case All:
		return nil, nil
	default:
		panic("清理策略不合法")
	}
}

func ConvClearType(typ int) *ClearType {
	clearType := ClearType(typ)
	if _, ok := clearTypeMap[clearType]; ok {
		return &clearType
	}
	return nil
}

func Values() []ClearType {
	return []ClearType{
		OneMonthAgo,
		ThreeMonthAgo,
		SixMonthAgo,
		OneYearAgo,
		OneThousandAgo,
		TenThousandAgo,
		ThirtyThousandAgo,
		OneHundredThousandAgo,
		All,
	}
}

const (
	OneMonthAgo ClearType = iota + 1
	ThreeMonthAgo
	SixMonthAgo
	OneYearAgo
	OneThousandAgo
	TenThousandAgo
	ThirtyThousandAgo
	OneHundredThousandAgo
	All
)

var clearTypeMap = map[ClearType]string{
	OneMonthAgo:           "清理近一个月之前的数据",
	ThreeMonthAgo:         "清理近三个月之前的数据",
	SixMonthAgo:           "清理近六个月之前的数据",
	OneYearAgo:            "清理近一年之前的数据",
	OneThousandAgo:        "清理近一千条之前的数据",
	TenThousandAgo:        "清理近一万条之前的数据",
	ThirtyThousandAgo:     "清理近三万条之前的数据",
	OneHundredThousandAgo: "清理近十万条之前的数据",
	All:                   "清理所有数据",
}
