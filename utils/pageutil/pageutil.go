package pageutil

import (
	"google.golang.org/protobuf/types/known/anypb"
	"gorm.io/gorm"
	"math"
)

//type Query struct {
//	PageNum  int64 `protobuf:"varint,1,opt,name=pageNum,proto3" json:"pageNum,omitempty"`
//	PageSize int64 `protobuf:"varint,2,opt,name=pageSize,proto3" json:"pageSize,omitempty"`
//}

//type Info struct {
//	Total    int64 `protobuf:"varint,1,opt,name=total,proto3" json:"total,omitempty"`
//	PageNum  int64 `protobuf:"varint,2,opt,name=pageNum,proto3" json:"pageNum,omitempty"`
//	PageSize int64 `protobuf:"varint,3,opt,name=pageSize,proto3" json:"pageSize,omitempty"`
//	Pages    int64 `protobuf:"varint,4,opt,name=pages,proto3" json:"pages,omitempty"`
//}

// 标准分页结构体，接收查询结果的结构体
type Page[T any] struct {
	Info    PageInfo
	List    []*T         `json:"data"`
	AnyList []*anypb.Any `json:"anyList"`
}

// 各种查询条件先在query设置好后再放进来
func GetResult[T any](db *gorm.DB, p *PageQuery) (res *Page[T], e error) {
	page := &Page[T]{}
	page.Info.PageNum = p.GetPageNum()
	page.Info.PageSize = p.GetPageSize()
	db.Count(&page.Info.Total)
	if page.Info.Total == 0 {
		page.List = []*T{}
		return
	}
	e = db.Scopes(Paginate(page)).Find(&page.List).Error
	return page, nil
}

// 各种查询条件先在query设置好后再放进来
func GetResultWithModel[T any](db *gorm.DB, p *PageQuery, t interface{}) (e error) {
	page := &Page[T]{}
	page.Info.PageNum = p.GetPageNum()
	page.Info.PageSize = p.GetPageSize()
	db.Count(&page.Info.Total)
	if page.Info.Total == 0 {
		page.List = []*T{}
		return
	}
	e = db.Model(&t).Scopes(Paginate(page)).Find(&page.List).Error
	return
}

func Paginate[T any](page *Page[T]) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page.Info.PageNum <= 0 {
			page.Info.PageNum = 1
		}
		switch {
		case page.Info.PageSize > 10000:
			page.Info.PageSize = 10000 // 限制一下分页大小
		case page.Info.PageSize <= 0:
			page.Info.PageSize = 10
		}
		// 使用浮点数除法，然后向上取整
		page.Info.Pages = int64(math.Ceil(float64(page.Info.Total) / float64(page.Info.PageSize)))
		p := page.Info.PageNum
		if page.Info.PageNum > page.Info.Pages {
			p = page.Info.Pages
		}
		size := page.Info.PageSize
		offset := int((p - 1) * size)
		return db.Offset(offset).Limit(int(size))
	}
}

func (page *Page[T]) ToAny() {
	// 使用空接口接收泛型参数
	var genericValue interface{} = page.List
	switch v := genericValue.(type) {
	case []*anypb.Any:
		page.AnyList = v
	}
}
