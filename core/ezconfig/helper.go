package ezconfig

import (
	"log/slog"
	"reflect"
)

var p *properties

func newProperties() *properties {
	prop := new(properties)
	prop.files = make([]string, 0)
	prop.configBeans = make([]ConfigurationBean, 0)
	prop.log = slog.Default()
	return prop
}

// GetByProperties 根据泛型类型获取配置对象
// 由于 ConfigurationBean 的唯一实现类 ConfigurationProperties 指定了指针实现
// 所以泛型T约束了只能是指针对象
func GetByProperties[T ConfigurationBean](prop Properties) T {
	rt := reflect.TypeOf((*T)(nil)).Elem().Elem()
	bean := prop.GetConfigBean(rt)
	if bean == nil {
		var zeroVal T
		return zeroVal
	}
	return any(bean).(T)
}

func Get[T ConfigurationBean]() T {
	return GetByProperties[T](p)
}

// PropertiesBuilder properties构建对象
type PropertiesBuilder struct {
	prop *properties // 配置对象指针
}

func Builder() *PropertiesBuilder {
	builder := new(PropertiesBuilder)
	p = newProperties()
	builder.prop = p
	return builder
}

func (my *PropertiesBuilder) AddFiles(files ...string) *PropertiesBuilder {
	my.prop.files = append(my.prop.files, files...)
	return my
}

func (my *PropertiesBuilder) AddConfigBeans(configBeans ...ConfigurationBean) *PropertiesBuilder {
	my.prop.configBeans = append(my.prop.configBeans, configBeans...)
	return my
}

func (my *PropertiesBuilder) AddWatcher(watchers ...WatchHandler) *PropertiesBuilder {
	my.prop.watchHandlers = append(my.prop.watchHandlers, watchers...)
	return my
}

func (my *PropertiesBuilder) Build() Properties {
	my.prop.init()
	return my.prop
}
