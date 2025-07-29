package req

type BasePage struct {
	CurrentPage int `json:"currentPage" form:"currentPage,default=1"` // 当前页
	PageSize    int `json:"pageSize" form:"pageSize,default=10"`      // 每页显示条数，默认 10
}

func (p *BasePage) Limit() int {
	return p.PageSize
}

func (p *BasePage) Start() int {
	return (p.CurrentPage - 1) * p.PageSize
}
