package page

import (
	//"github.com/samsaralc/proto_repo/pb/github.com/jiushengTech/common/page"
	"google.golang.org/protobuf/types/known/anypb"
	"gorm.io/gorm"
	"math"
)

type Query struct {
	PageNum  int64 `protobuf:"varint,1,opt,name=pageNum,proto3" json:"pageNum,omitempty"`
	PageSize int64 `protobuf:"varint,2,opt,name=pageSize,proto3" json:"pageSize,omitempty"`
}

type Info struct {
	Total    int64 `protobuf:"varint,1,opt,name=total,proto3" json:"total,omitempty"`
	PageNum  int64 `protobuf:"varint,2,opt,name=pageNum,proto3" json:"pageNum,omitempty"`
	PageSize int64 `protobuf:"varint,3,opt,name=pageSize,proto3" json:"pageSize,omitempty"`
	Pages    int64 `protobuf:"varint,4,opt,name=pages,proto3" json:"pages,omitempty"`
}

// Page 标准分页结构体，接收查询结果的结构体
type Page[T any] struct {
	info    Info
	List    []*T         `json:"data"`
	AnyList []*anypb.Any `json:"anyList"`
}

// 各种查询条件先在query设置好后再放进来
func GetResult[T any](db *gorm.DB, pageNum int64, pageSize int64) (res *Page[T], e error) {
	page := &Page[T]{}
	page.info.PageNum = pageNum
	page.info.PageSize = pageSize
	db.Count(&page.info.Total)
	if page.info.Total == 0 {
		page.List = []*T{}
		return
	}
	e = db.Scopes(Paginate(page)).Find(&page.List).Error
	return page, nil
}

// 各种查询条件先在query设置好后再放进来
func GetResultWithModel[T any](db *gorm.DB, pageNum int64, pageSize int64, t interface{}) (e error) {
	page := &Page[T]{}
	page.info.PageNum = pageNum
	page.info.PageSize = pageSize
	db.Count(&page.info.Total)
	if page.info.Total == 0 {
		page.List = []*T{}
		return
	}
	e = db.Model(&t).Scopes(Paginate(page)).Find(&page.List).Error
	return
}

func Paginate[T any](page *Page[T]) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page.info.PageNum <= 0 {
			page.info.PageNum = 1
		}
		switch {
		case page.info.PageSize > 10000:
			page.info.PageSize = 10000 // 限制一下分页大小
		case page.info.PageSize <= 0:
			page.info.PageSize = 10
		}
		// 使用浮点数除法，然后向上取整
		page.info.Pages = int64(math.Ceil(float64(page.info.Total) / float64(page.info.PageSize)))
		p := page.info.PageNum
		if page.info.PageNum > page.info.Pages {
			p = page.info.Pages
		}
		size := page.info.PageSize
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
