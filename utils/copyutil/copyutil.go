package copyutil

import (
	ctime "github.com/jiushengTech/common/time"
	"reflect"
	"time"
)

// localTimeToString 将 LocalTime 类型转为 string
func localTimeToString(lt ctime.LocalTime) string {
	return time.Time(lt).Format(time.DateTime)
}

// ConvertLocalTimeToString 将源结构体中 common/time 包中的 LocalTime 类型转为字符串，并复制到目标结构体
func ConvertLocalTimeToString(target any, source any) {
	// 获取目标结构体和源结构体的反射值
	targetValue := reflect.ValueOf(target).Elem()
	sourceValue := reflect.ValueOf(source)

	// 遍历源结构体的字段
	for i := 0; i < sourceValue.NumField(); i++ {
		// 获取源结构体字段和目标结构体对应字段的反射值
		sourceField := sourceValue.Field(i)
		targetField := targetValue.FieldByName(sourceValue.Type().Field(i).Name)

		// 检查目标结构体字段是否存在
		if targetField.IsValid() {
			// 如果是 common/time 包中的 LocalTime 类型，进行转换
			if sourceField.Type() == reflect.TypeOf(ctime.LocalTime{}) {
				targetField.SetString(localTimeToString(sourceField.Interface().(ctime.LocalTime)))
			} else if targetField.Type() == sourceField.Type() {
				// 如果字段类型相同，直接赋值
				targetField.Set(sourceField)
			}
		}
	}
}

// CopyStruct 通过反射复制源结构体的字段到目标结构体。
// 注意：该函数要求目标结构体字段名称和类型必须与源结构体一一对应且可导出。
// 参数 source 是源结构体的实例，参数 target 是目标结构体的指针。
func CopyStruct(target any, source any) {
	targetValue := reflect.ValueOf(target).Elem()
	sourceValue := reflect.ValueOf(source)

	for i := 0; i < sourceValue.NumField(); i++ {
		sourceField := sourceValue.Field(i)
		targetField := targetValue.FieldByName(sourceValue.Type().Field(i).Name)

		if targetField.IsValid() && targetField.Type() == sourceField.Type() {
			targetField.Set(sourceField)
		}
	}
}
