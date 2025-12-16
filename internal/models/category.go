package models

type Category struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"uniqueIndex;not null;type:varchar(100)" json:"name"`

	Articles []Article `gorm:"many2many:article_categories;" json:"-"`
}

func (Category) TableName() string {
	return "categories"
}

type ArticleCategory struct {
	ArticleID  uint `gorm:"primaryKey;column:article_id"`
	CategoryID uint `gorm:"primaryKey;column:category_id"`

	Article  Article `gorm:"foreignKey:ArticleID" json:"-"`
	Category Category `gorm:"foreignKey:CategoryID" json:"-"`
}

func (ArticleCategory) TableName() string {
	return "article_categories"
}
