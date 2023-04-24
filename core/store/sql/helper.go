package sql

import (
	"reflect"
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Option func(*gorm.DB) *gorm.DB

func OptionDB(db *gorm.DB, opts ...Option) *gorm.DB {
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

// Paginate 分页
func paginate(page, pageSize *int) func(db *gorm.DB) *gorm.DB {
	if *page <= 0 {
		*page = 1
	}
	if *pageSize <= 0 {
		*pageSize = 10
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset((*page - 1) * *pageSize).Limit(*pageSize)
	}
}

// WithPage 分页
func WithPage(page, pageSize *int) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(paginate(page, pageSize))
	}
}

// WithKey 指定关键字进行查询 0值不进行查询
func WithKey(cond map[string]interface{}) Option {
	for k, v := range cond {
		if reflect.ValueOf(v).IsZero() {
			delete(cond, k)
		}
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(cond)
	}
}

// WithKeyNil 指定关键字进行查询 零值同样查询
func WithKeyNil(cond map[string]interface{}) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(cond)
	}
}
