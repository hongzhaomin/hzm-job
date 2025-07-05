package ezconfig

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/hongzhaomin/hzm-job/core/internal/tools"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"reflect"
	"slices"
	"strings"
	"time"
)

// aliasTag 可以给 ConfigurationBean 中的属性设置别名
const aliasTag = "alias"

// Properties 配置核心接口，也是easyconfig的核心
type Properties interface {

	// Get 同viper.Get()方法
	Get(key string) any

	// GetString 同viper.GetString()方法
	GetString(key string) string

	// GetInt 同viper.GetInt()方法
	GetInt(key string) int

	// GetInt64 同viper.GetInt64()方法
	GetInt64(key string) int64

	// GetFloat64 同viper.GetFloat64()方法
	GetFloat64(key string) float64

	// GetBool 同viper.GetBool()方法
	GetBool(key string) bool

	// ContainsKey 判断key是否存在
	ContainsKey(key string) bool

	// GetConfigBean 根据反射类型查询配置类对象
	GetConfigBean(rt reflect.Type) ConfigurationBean

	// GetConfigBeans 返回所有配置类对象
	GetConfigBeans() []ConfigurationBean

	// ResolveAndSetConfigBeans 解析配置类并做属性赋值
	ResolveAndSetConfigBeans(configBeans ...ConfigurationBean) error
}

var _ Properties = (*properties)(nil)

type properties struct {
	files              []string                                          // 文件名列表（含扩展名）
	configBeans        []ConfigurationBean                               // 配置类列表（必须为指针对象）
	watchHandlers      []WatchHandler                                    // 监视处理器列表
	file2Viper         map[string]*viper.Viper                           // 文件名对应viper.Viper对象map
	conf2DefinitionMap map[ConfigurationBean]ConfigurationBeanDefinition // 配置类对应配置定义对象map
	mergedViper        *viper.Viper                                      // 合并的viper.Viper
	log                *slog.Logger                                      // 日志对象
}

func (my *properties) Get(key string) any { return my.mergedViper.Get(key) }

func (my *properties) GetString(key string) string { return my.mergedViper.GetString(key) }

func (my *properties) GetInt(key string) int { return my.mergedViper.GetInt(key) }

func (my *properties) GetInt64(key string) int64 { return my.mergedViper.GetInt64(key) }

func (my *properties) GetFloat64(key string) float64 { return my.mergedViper.GetFloat64(key) }

func (my *properties) GetBool(key string) bool { return my.mergedViper.GetBool(key) }

func (my *properties) ContainsKey(key string) bool {
	return my.mergedViper.IsSet(key)
}

func (my *properties) GetConfigBean(rt reflect.Type) ConfigurationBean {
	if rt == nil {
		return nil
	}

	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	for _, bean := range my.configBeans {
		rtBean := reflect.Indirect(reflect.ValueOf(bean)).Type()
		// 反射类型可以直接使用 == 比较
		if rt == rtBean {
			return bean
		}
	}
	return nil
}

func (my *properties) GetConfigBeans() []ConfigurationBean {
	configBeanLen := len(my.configBeans)
	if configBeanLen <= 0 {
		return nil
	}
	configs := make([]ConfigurationBean, 0, configBeanLen)
	return append(configs, my.configBeans...)
}

func (my *properties) refreshMergedViper() (oldVp *viper.Viper) {
	oldVp = my.mergedViper

	my.mergedViper = viper.New()
	for _, vp := range my.file2Viper {
		err := my.mergedViper.MergeConfigMap(vp.AllSettings())
		if err != nil {
			my.log.Error("get merged viper error: ", err)
		}
	}
	return
}

func (my *properties) init() {
	my.log.Debug("start init easyconfig >>>>>>>>>")
	files := my.files
	if len(files) <= 0 {
		panic(errors.New("files is empty"))
	}

	// 读取配置文件
	my.file2Viper = make(map[string]*viper.Viper, len(files))
	for _, f := range files {
		if f == "" {
			panic(errors.New("file is empty"))
		}
		my.log.Debug("reading config file %s", f)
		vp := my.getViperAndReadInConfig(f)
		my.file2Viper[f] = vp
	}
	// 刷新合并的viper.Viper
	my.refreshMergedViper()
	my.log.Debug("read config files completed")

	// 添加默认监视处理器，并放在第一个
	my.watchHandlers = append([]WatchHandler{my.defaultWatchHandler}, my.watchHandlers...)

	// 注册监视器
	my.registryWatcher()
	my.log.Debug("registry watchers completed")

	// 解析配置类并做属性赋值
	my.conf2DefinitionMap = make(map[ConfigurationBean]ConfigurationBeanDefinition, 16)
	if err := my.ResolveAndSetConfigBeans(my.configBeans...); err != nil {
		panic(err)
	}

	my.log.Debug("end init easyconfig >>>>>>>>>")
}

// ResolveAndSetConfigBeans 解析配置类并做属性赋值
func (my *properties) ResolveAndSetConfigBeans(configBeans ...ConfigurationBean) error {
	if len(configBeans) == 0 {
		return nil
	}

	// 解析配置类信息
	my.configBeans = append(my.configBeans, configBeans...)
	for _, cnf := range configBeans {
		if reflect.ValueOf(cnf).Kind() != reflect.Ptr {
			return errors.New("easyconfig bean must be ptr")
		}
		if cnf == nil {
			return errors.New("easyconfig bean is empty")
		}
		definition := ResolveConfigurationBean(cnf)
		my.log.Debug("resolve config bean [%s] completed", definition.CompleteName)
		my.conf2DefinitionMap[cnf] = definition
	}

	// 配置类属性赋值
	my.log.Debug("binding config to beans")
	my.setConfigBeans(configBeans)
	my.log.Debug("bind config to beans completed")
	return nil
}

func (my *properties) getViperAndReadInConfig(filePath string) *viper.Viper {
	vp := viper.New()
	vp.SetConfigFile(filePath)
	if err := vp.ReadInConfig(); err != nil {
		msg := fmt.Sprintf("viper.ReadInConfig() err: %v\n", err)
		panic(errors.New(msg))
	}
	return vp
}

func (my *properties) setConfigBeans(configBeans []ConfigurationBean) {
	if len(configBeans) <= 0 {
		return
	}

	for _, configBean := range configBeans {
		definition := my.conf2DefinitionMap[configBean]
		rvConfigBean := reflect.ValueOf(configBean)
		if rvConfigBean.Type().Kind() != reflect.Ptr {
			// 不是指针类型，跳过
			continue
		}

		if !definition.AutoRefresh {
			// 不自动刷新，第一次赋值，后面配置变更将不会刷新配置
			// 在此打一个warn日志
			my.log.Warn("config bean [%s] will not be auto refreshed", definition.CompleteName)
		}

		// 忽略无法转换的无效的属性（即忽略那些配置中属性类型无法转换为结构体属性类型的字段）
		ignoreInvalidFields := definition.IgnoreInvalidFields
		//// 忽略未知的属性（即忽略那些配置中有但结构体中没有的字段）
		//ignoreUnknownFields := definition.IgnoreUnknownFields

		key := definition.Prefix
		var mapValue map[string]any
		if key == "" {
			mapValue = my.mergedViper.AllSettings()
		} else {
			mapValue = my.mergedViper.GetStringMap(key)
		}
		if mapValue == nil {
			if !ignoreInvalidFields {
				panic(errors.New(fmt.Sprintf("easyconfig value is empty for key [%s]", key)))
			}
			return
		}

		my.doSetConfigBean(configBean, definition, mapValue, key, false)
	}
}

func (my *properties) doSetConfigBean(configBean any, definition ConfigurationBeanDefinition,
	mapValue map[string]any, viperKey string, ignoreEmpty bool) {
	// 忽略无法转换的无效的属性（即忽略那些配置中属性类型无法转换为结构体属性类型的字段）
	ignoreInvalidFields := definition.IgnoreInvalidFields

	var rv reflect.Value
	switch val := configBean.(type) {
	case reflect.Value:
		rv = reflect.Indirect(val)
	default:
		rv = reflect.Indirect(reflect.ValueOf(val))
		if rv.Kind() != reflect.Struct {
			panic(errors.New("source and target must be struct kind"))
		}
	}

	// 开始遍历结构体属性
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		structField := rt.Field(i)
		if !structField.IsExported() {
			// 不可导出，跳过
			continue
		}

		// 字段反射类型
		rtField := structField.Type
		isPtr := false
		if rtField.Kind() == reflect.Ptr {
			rtField = rtField.Elem()
			isPtr = true
		}
		// 字段反射值
		rvField := rv.Field(i)

		if structField.Anonymous {
			// 如果是内嵌字段，并且不是 注解类ConfigurationProperties，则当成父类注入配置
			isConfigurationProperties := rtField == rtConfigurationProperties
			if rtField.Kind() == reflect.Struct && !isConfigurationProperties {
				if isPtr {
					// 如果是指针类型的结构体，则先创建指针结构体对象，再递归，然后再把指针对象赋值到当前rvField上即可
					// 否则，会报错，因为此时的属性值是nil，也就是rvField的值是无效的反射值
					newRvField := reflect.New(rtField)
					rvField.Set(newRvField)
				}
				my.doSetConfigBean(rvField, definition, mapValue, viperKey, ignoreEmpty)
			} else {
				if !isConfigurationProperties {
					// 其他类型的内嵌字段，全部忽略
					my.log.Warn("ignored field [%s] in config bean [%s]", structField.Name, definition.CompleteName)
				}
			}
			continue
		}

		// 字段名称
		fileName := structField.Name
		// 根据字段名称获取viper配置中的配置值，如果设置了别名，则以别名为准
		mapKey := strings.ToLower(fileName)
		if alias, ok := structField.Tag.Lookup(aliasTag); ok {
			mapKey = strings.ToLower(alias)
		}
		nextViperKey := viperKey + "." + mapKey
		fieldMapVal := mapValue[mapKey]
		if fieldMapVal == nil {
			if !ignoreInvalidFields && !ignoreEmpty {
				panic(errors.New(fmt.Sprintf("easyconfig value is empty for key [%s.%s]", viperKey, mapKey)))
			}
			continue
		}

		switch rtField.Kind() {
		case reflect.Map:
			// 配置值的反射值
			rvMapVal := reflect.Indirect(reflect.ValueOf(fieldMapVal))
			if rvMapVal.Kind() != reflect.Map {
				if !ignoreInvalidFields {
					panic(errors.New(fmt.Sprintf("easyconfig value is not convert to field type for key [%s.%s]",
						viperKey, mapKey)))
				}
				continue
			}
			// 获取map属性的值反射类型
			rtFieldMapValue := rtField.Elem()
			// 根据map属性类型，反射创建map
			rvFieldMap := reflect.MakeMapWithSize(rtField, rvMapVal.Len())
			for _, mk := range rvMapVal.MapKeys() {
				// 获取配置值map中的value反射值
				rvMapValEle := rvMapVal.MapIndex(mk)
				// 将配置值map中的value转换成属性map的value类型的 反射值
				rvFieldMapValue, err := tools.ReflectConvert4Any(rtFieldMapValue, reflect.Indirect(rvMapValEle).Interface())
				if err != nil {
					if !ignoreInvalidFields {
						panic(errors.New(fmt.Sprintf("easyconfig value is not convert to field type for key [%s.%s]: err: %s",
							viperKey, mapKey, err.Error())))
					}
					break
				}
				// map的value值类型是否为指针
				if rtFieldMapValue.Kind() == reflect.Ptr {
					rvFieldMap.SetMapIndex(mk, rvFieldMapValue)
				} else {
					rvFieldMap.SetMapIndex(mk, rvFieldMapValue.Elem())
				}
			}
			// map是否为指针map
			if isPtr {
				// 创建一个map的指针对象
				rvFieldMapPtr := reflect.New(rtField)
				// 将获得的值赋值到这个指针map对象上
				rvFieldMapPtr.Elem().Set(rvFieldMap)
				// 最后将这个map指针对象赋值到rvField上
				rvField.Set(rvFieldMapPtr)
			} else {
				rvField.Set(rvFieldMap)
			}
		case reflect.Slice, reflect.Array:
			// 获取切片属性的值反射类型
			rtFieldArrValue := rtField.Elem()
			rtFieldMapValueIsPtr := false
			if rtFieldArrValue.Kind() == reflect.Ptr {
				rtFieldArrValue = rtFieldArrValue.Elem()
				rtFieldMapValueIsPtr = true
			}
			supportKinds := []reflect.Kind{reflect.Int8, reflect.Int16, reflect.Int, reflect.Int32, reflect.Int64,
				reflect.Uint8, reflect.Uint16, reflect.Uint, reflect.Uint32, reflect.Uint64,
				reflect.Float32, reflect.Float64, reflect.String, reflect.Bool}
			if !slices.Contains(supportKinds, rtFieldArrValue.Kind()) {
				slog.Warn(fmt.Sprintf("easyconfig value is not support [%s] slice for key [%s.%s]",
					rtFieldArrValue.Kind(), viperKey, mapKey))
				continue
			}
			// 配置值的反射值
			rvMapVal := reflect.Indirect(reflect.ValueOf(fieldMapVal))
			if rvMapVal.Kind() == reflect.Slice || rvMapVal.Kind() == reflect.Array {
				// nothing to do
			} else if rvMapVal.Kind() == reflect.String {
				// 支持string，逗号分割成list
				strVal := fieldMapVal.(string)

				var strSlice []string
				for _, str := range strings.Split(strVal, ",") {
					str = strings.TrimSpace(str)
					if str != "" {
						strSlice = append(strSlice, str)
					}
				}
				rvMapVal = reflect.ValueOf(strSlice)
			} else {
				if !ignoreInvalidFields {
					panic(errors.New(fmt.Sprintf("easyconfig value is not convert to field type for key [%s.%s]",
						viperKey, mapKey)))
				}
				continue
			}

			newSlice := reflect.MakeSlice(rtField, 0, rvMapVal.Len())
			for j := 0; j < rvMapVal.Len(); j++ {
				rvMapValEle := rvMapVal.Index(j)
				// 将配置值map中的value转换成属性map的value类型的 反射值
				rvFieldArrValue, err := tools.ReflectConvert4Any(rtFieldArrValue, reflect.Indirect(rvMapValEle).Interface())
				if err != nil {
					if !ignoreInvalidFields {
						panic(errors.New(fmt.Sprintf("easyconfig value is not convert to field type for key [%s.%s]: err: %s",
							viperKey, mapKey, err.Error())))
					}
					break
				}
				// 切片的元素类型是否指针类型
				if rtFieldMapValueIsPtr {
					newSlice = reflect.Append(newSlice, rvFieldArrValue)
				} else {
					newSlice = reflect.Append(newSlice, rvFieldArrValue.Elem())
				}
			}
			// 切片是否为指针切片
			if isPtr {
				// 创建一个切片的指针对象
				rvFieldSlicePtr := reflect.New(rtField)
				// 将获得的值赋值到这个指针切片对象上
				rvFieldSlicePtr.Elem().Set(newSlice)
				// 最后将这个切片指针对象赋值到rvField上
				rvField.Set(rvFieldSlicePtr)
			} else {
				rvField.Set(newSlice)
			}
		case reflect.Int8, reflect.Int16, reflect.Int, reflect.Int32, reflect.Int64,
			reflect.Uint8, reflect.Uint16, reflect.Uint, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64, reflect.String, reflect.Bool:
			// 配置值的反射值
			var rvMapVal reflect.Value
			if valStr, ok := fieldMapVal.(string); ok {
				// 如果取出来的值是string类型的，进一步监测一下值里面是否有类似${server.port}表达式，有的话也需要转化
				rvMapVal = reflect.Indirect(reflect.ValueOf(getValFromViper(my.mergedViper, "", valStr)))
			} else {
				rvMapVal = reflect.Indirect(reflect.ValueOf(fieldMapVal))
			}

			rvFieldConv, err := tools.ReflectConvert4Any(rtField, rvMapVal.Interface())
			if err != nil {
				if !ignoreInvalidFields {
					panic(errors.New(fmt.Sprintf("easyconfig value is not convert to field type for key [%s.%s]: err: %s",
						viperKey, mapKey, err.Error())))
				}
				continue
			}
			if isPtr {
				rvField.Set(rvFieldConv)
			} else {
				rvField.Set(rvFieldConv.Elem())
			}
		case reflect.Struct:
			nextMapValue, ok := fieldMapVal.(map[string]any)
			if !ok {
				if !ignoreInvalidFields {
					panic(errors.New(fmt.Sprintf("easyconfig value type is not map[string]any, key: [%s.%s]",
						viperKey, mapKey)))
				}
				continue
			}
			if isPtr {
				// 如果是指针类型的结构体，则先创建指针结构体对象，再递归，然后再把指针对象赋值到当前rvField上即可
				// 否则，会报错，因为此时的属性值是nil，也就是rvField的值是无效的反射值
				newRvField := reflect.New(rtField)
				rvField.Set(newRvField)
			}
			my.doSetConfigBean(rvField, definition, nextMapValue, nextViperKey, ignoreEmpty)
		default:
			// nothing to do
		}
	}
}

func (my *properties) registryWatcher() {
	for _, vp := range my.file2Viper {
		// 添加 viper.Viper 监视器
		vp.OnConfigChange(my.onConfigChange)
		// 监视
		vp.WatchConfig()
	}
}

// onConfigChange viper.Viper监视器实现
func (my *properties) onConfigChange(in fsnotify.Event) {
	defer func() {
		if err := recover(); err != nil {
			my.log.Error("auto refresh easyconfig failed: ", err)
		}
	}()
	filePath := in.Name
	my.log.Debug("watched file [%s] changed", filePath)

	// 1、发现一个viper的bug：如果文件使用notepad++打开，更新内容，会发生viper.Viper.AllSettings()的值变为空，
	// 导致my.refreshMergedViper()无作用，配置得不到刷新，而且监视器会调用2次；
	// 但是在viper.WatchConfig()源码中调用fun(in fsnotify.Event)函数前打上断点，则又能正常刷新，监视器也是调用2次。
	// 2、基于这个bug，这里做一下特殊处理：主动去重新读取一下文件，最多重试3次。
	vp := my.file2Viper[filePath]
	for i := 0; i < 3; i++ {
		if len(vp.AllSettings()) <= 0 {
			if err := vp.ReadInConfig(); err != nil {
				msg := fmt.Sprintf("viper.ReadInConfig() err: %v\n", err)
				my.log.Error("viper reload config err: %s", msg)
				time.Sleep(3 * time.Millisecond)
			}
			continue
		}
		break
	}

	oldVp := my.refreshMergedViper()

	if len(my.conf2DefinitionMap) <= 0 {
		return
	}

	if len(my.watchHandlers) <= 0 {
		return
	}

	if in.Has(fsnotify.Write) {
		// 找出配置改变的key，只找key不变，值发生改变的配置
		// 如果删除key（配了key但未配值视为同类），则不会刷新，因为刷新可能给程序带来崩溃
		// 如果增加key，也不会刷新，因为此时你的结构体可能没有定义该字段
		params := make([]WatcherParam, 0)
		for _, k := range vp.AllKeys() {
			if oldVp.IsSet(k) {
				oldVal := oldVp.Get(k)
				oldValStr, ok := oldVal.(string)
				if ok {
					oldVal = getValFromViper(oldVp, k, oldValStr)
				}

				newVal := vp.Get(k)
				newValStr, ok := newVal.(string)
				if ok {
					newVal = getValFromViper(vp, k, newValStr)
				}

				if oldVal == newVal {
					continue
				}
				my.log.Debug("watched config changed, key: [%s], oldVal: [%v], newVal: [%v]", k, oldVal, newVal)
				params = append(params, WatcherParam{
					Key:    k,
					OldVal: oldVal,
					NewVal: newVal,
				})
			}
		}

		// 有的时候某一个配置文件发生变化，会间接引起其他配置文件发生变化，这个时候 viper 就监测不到了
		// 例如：a 文件有一个 [port] 的配置，b文件的 [copy.port] 配置引用了 [port] 的配置
		// 所以这里查找下引用到变更key的其他key
		for path, otherVp := range my.file2Viper {
			if filePath == path {
				// 变更配置的文件跳过
				continue
			}
			for _, k := range otherVp.AllKeys() {
				if v, ok := otherVp.Get(k).(string); ok {
					for _, param := range params {
						if !strings.Contains(v, fmt.Sprintf("${%s}", param.Key)) {
							continue
						}
						oldVal := os.Expand(v, func(embeddedKey string) string {
							if embeddedKey == param.Key {
								old, err := tools.ConvertStr4Any(param.OldVal)
								if err != nil {
									my.log.Error("convert old value err: ", err)
								}
								return old
							}
							return my.mergedViper.GetString(embeddedKey)
						})
						newVal := os.Expand(v, func(embeddedKey string) string {
							if embeddedKey == param.Key {
								newV, err := tools.ConvertStr4Any(param.NewVal)
								if err != nil {
									my.log.Error("convert old value err: ", err)
								}
								return newV
							}
							return my.mergedViper.GetString(embeddedKey)
						})

						if oldVal == newVal {
							continue
						}

						params = append(params, WatcherParam{
							Key:    k,
							OldVal: oldVal,
							NewVal: newVal,
						})
					}
				}

			}
		}

		// 循环执行配置 easyconfig 自定义 监测处理器
		if len(params) > 0 {
			for _, watcher := range my.watchHandlers {
				watcher(params, my.conf2DefinitionMap)
			}
		}
	}
}

func (my *properties) defaultWatchHandler(params []WatcherParam, conf2DefinitionMap map[ConfigurationBean]ConfigurationBeanDefinition) {
	// 找出需要刷新的配置类
	refreshConfigBeans := make([]ConfigurationBean, 0, len(conf2DefinitionMap))
	for configBean, definition := range conf2DefinitionMap {
		if definition.AutoRefresh {
			refreshConfigBeans = append(refreshConfigBeans, configBean)
		}
	}
	// 没有需要刷新的配置类，结束
	if len(refreshConfigBeans) <= 0 {
		return
	}

	// 将监视参数转化为map[string]any（key中不包含"."的map）
	rootMap := watcherParams(params).convMap()

	// 开始刷新配置
	for _, configBean := range refreshConfigBeans {
		definition := conf2DefinitionMap[configBean]
		prefixKey := definition.Prefix

		mapValue, ok := my.findMapByKey(rootMap, strings.Split(prefixKey, ".")).(map[string]any)
		if !ok {
			continue
		}
		my.doSetConfigBean(configBean, definition, mapValue, prefixKey, true)
		my.log.Debug("config bean [%s] auto refresh success", definition.CompleteName)
	}
}

func (my *properties) findMapByKey(sourceMap map[string]any, keys []string) any {
	if len(sourceMap) <= 0 {
		return nil
	}

	if len(keys) <= 0 {
		return sourceMap
	}

	k := keys[0]
	next, ok := sourceMap[k]
	if !ok {
		return nil
	}
	if len(keys) == 1 {
		return next
	}

	nextMap, ok := next.(map[string]any)
	if !ok {
		return nil
	}
	return my.findMapByKey(nextMap, keys[1:])
}

// =========================== 自定义easyconfig配置监测处理器 ============================

// WatchHandler 定义easyconfig配置监测处理器
type WatchHandler func(params []WatcherParam, conf2DefinitionMap map[ConfigurationBean]ConfigurationBeanDefinition)

// WatcherParam 定义easyconfig监视器参数
type WatcherParam struct {
	Key    string // 配置发生改变的key
	OldVal any    // 配置改变前的值
	NewVal any    // 配置改变后的值
}

// =========================== 定义监视参数切片类型并实现map转化方法 ============================

// 定义监视参数列表，实现map转化
type watcherParams []WatcherParam

func (my watcherParams) convMap() map[string]any {
	allSetting := make(map[string]any)
	for _, param := range my {
		allSetting[param.Key] = param.NewVal
	}

	rootMap := make(map[string]any)
	for k, v := range allSetting {
		keys := strings.Split(k, ".")

		lastKey := keys[len(keys)-1]
		tmpMap := rootMap
		for _, key := range keys {
			if key == lastKey {
				tmpMap[key] = v
				continue
			}

			if subMap, ok := tmpMap[key]; ok {
				tmpMap = subMap.(map[string]any)
			} else {
				newMap := make(map[string]any)
				tmpMap[key] = newMap
				tmpMap = newMap
			}
		}
	}
	return rootMap
}

// 如果配置中引用了其他配置的值，需要转化一下
func getValFromViper(vp *viper.Viper, key, val string) string {
	if key == "" && val == "" {
		return ""
	}

	if key != "" {
		val = vp.GetString(key)
	}
	return os.Expand(val, func(embeddedKey string) string {
		return getValFromViper(vp, embeddedKey, "")
	})
}
