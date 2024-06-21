package models

type Comment struct {
	ID       uint   `gorm:"primaryKey; autoIncrement" json:"id"`
	Content  string `json:"content"`
	UserID   uint   `json:"user_id"`
	ArticleID uint  `json:"article_id"`
}