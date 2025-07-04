package ezconfig

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

const (
	prefixTag              = "prefix"              // 配置属性的前缀, 默认""
	ignoreInvalidFieldsTag = "ignoreInvalidFields" // 忽略无法转换的无效的属性, 默认false（即忽略那些配置中属性类型无法转换为结构体属性类型的字段）
	autoRefreshTag         = "autoRefresh"         // 当配置发生变化时，是否自动刷新, 默认true
)

// ConfigurationBean 配置类接口，标识结构体为一个配置类
type ConfigurationBean interface {
	configurationBeanAnnotation()
}

var _ ConfigurationBean = (*ConfigurationProperties)(nil)

// ConfigurationProperties ConfigurationBean的唯一实现类，提供注解的作用
// 所有需要标记为配置类的结构体都必须将 ConfigurationProperties 定义为该结构体的内嵌属性，这样才能间接实现 ConfigurationBean 接口
type ConfigurationProperties struct {
}

// 约束了指针实现
func (*ConfigurationProperties) configurationBeanAnnotation() {}

// ConfigurationBeanDefinition 解析ConfigurationProperties注解的定义信息
type ConfigurationBeanDefinition struct {
	Prefix              string // 配置属性的前缀
	IgnoreInvalidFields bool   // 忽略无法转换的无效的属性（即忽略那些配置中属性类型无法转换为结构体属性类型的字段）
	AutoRefresh         bool   // 当配置发生变化时，是否自动刷新
	CompleteName        string // 配置类全路径名称
	Name                string // 配置类名称
}

var rtConfigurationProperties = reflect.TypeOf((*ConfigurationProperties)(nil)).Elem()

// ResolveConfigurationBean 解析配置类 ConfigurationBean 的注解元信息
func ResolveConfigurationBean[T ConfigurationBean](configBean T) ConfigurationBeanDefinition {
	rt := reflect.TypeOf(configBean)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	completeName := rt.PkgPath() + "." + rt.Name()
	var annoStructField *reflect.StructField
	for i := 0; i < rt.NumField(); i++ {
		structField := rt.Field(i)
		if !structField.Anonymous {
			continue
		}

		rtField := structField.Type
		if rtField.Kind() == reflect.Ptr {
			rtField = rtField.Elem()
		}

		if rtField == rtConfigurationProperties {
			annoStructField = &structField
		}
	}
	if annoStructField == nil {
		msg := fmt.Sprintf("configBean [%s] not exist anonymous field of type ConfigurationProperties struct",
			completeName)
		panic(errors.New(msg))
	}

	prefixVal := ""
	ignoreInvalidFieldsVal := false
	autoRefreshVal := true
	if value, ok := annoStructField.Tag.Lookup(prefixTag); ok {
		prefixVal = value
	}
	if value, ok := annoStructField.Tag.Lookup(ignoreInvalidFieldsTag); ok {
		if parseBool, err := strconv.ParseBool(value); err == nil {
			ignoreInvalidFieldsVal = parseBool
		}
	}
	if value, ok := annoStructField.Tag.Lookup(autoRefreshTag); ok {
		if parseBool, err := strconv.ParseBool(value); err == nil {
			autoRefreshVal = parseBool
		}
	}
	return ConfigurationBeanDefinition{
		Prefix:              prefixVal,
		IgnoreInvalidFields: ignoreInvalidFieldsVal,
		AutoRefresh:         autoRefreshVal,
		CompleteName:        completeName,
		Name:                rt.Name(),
	}
}
