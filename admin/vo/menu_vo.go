package vo

type MenuInfo struct {
	Id       int64      `json:"id"`       // 菜单id
	ParentId int64      `json:"parentId"` // 菜单父id
	Title    string     `json:"title"`    // 名称
	Href     string     `json:"href"`     // 路由地址
	Icon     string     `json:"icon"`     // icon图标
	Target   string     `json:"target"`   // 同a标签的target属性
	Children []MenuInfo `json:"children"` // 子菜单列表
}
