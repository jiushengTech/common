// excel.go
package excel

import (
	"errors"
	"fmt"
	"github.com/tealeg/xlsx"
	"os"
	"reflect"
	"time"
)

// ExportStructToExcel exports data from a struct to an Excel file.
// support maxsize is 5m,if exceed ,it will be discarded
func ExportStructToExcel(data interface{}, filePath string, sheetName string, maxSize int64) error {
	if maxSize == 0 {
		maxSize = 20 * 1024 * 1024
	}
	// 检查文件是否存在
	fileInfo, err := os.Stat(filePath)
	if err == nil && fileInfo.Size() > maxSize { // 大小超过5MB
		return nil
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}

	file := xlsx.NewFile()
	var sheet *xlsx.Sheet

	// 如果文件存在，打开现有文件并获取 Sheet
	if _, err := xlsx.OpenFile(filePath); err == nil {
		file, err = xlsx.OpenFile(filePath)
		if err != nil {
			return err
		}

		// 查找指定 Sheet
		for _, existingSheet := range file.Sheets {
			if existingSheet.Name == sheetName {
				sheet = existingSheet
				break
			}
		}
	}

	// 如果 Sheet 不存在，创建新 Sheet
	if sheet == nil {
		// 使用外部变量，而不是创建局部变量
		var err error
		sheet, err = file.AddSheet(sheetName)
		if err != nil {
			return err
		}

		// 添加列名
		headerRow := sheet.AddRow()
		val := reflect.Indirect(reflect.ValueOf(data))
		typ := val.Type()
		for i := 0; i < val.NumField(); i++ {
			field := typ.Field(i)
			cell := headerRow.AddCell()
			cell.SetString(field.Name)
		}
	}

	// 添加数据
	dataRow := sheet.AddRow()
	val := reflect.Indirect(reflect.ValueOf(data))
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		cell := dataRow.AddCell()

		switch field.Type.Kind() {
		case reflect.String:
			cell.SetString(fmt.Sprintf("%v", val.Field(i).Interface()))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			cell.SetInt(int(val.Field(i).Int()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			cell.SetInt(int(val.Field(i).Uint()))
		case reflect.Float32, reflect.Float64:
			cell.SetFloat(val.Field(i).Float())
		case reflect.Bool:
			cell.SetBool(val.Field(i).Bool())
		case reflect.Struct:
			if field.Type == reflect.TypeOf(time.Time{}) {
				cell.SetDateTime(val.Field(i).Interface().(time.Time))
			}
		case reflect.Func:
			// 处理函数类型，这里可以根据需要进行定制
			cell.SetString("FunctionType")
		default:
			return errors.New("unsupported type")
		}
	}

	// 保存 Excel 文件
	err = file.Save(filePath)
	if err != nil {
		return err
	}
	return nil
}
