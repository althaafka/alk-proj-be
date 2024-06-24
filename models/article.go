package models

type Article struct {
	ID       uint   `gorm:"primaryKey; autoIncrement" json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	UserID   uint   `json:"user_id"`
	Likes    uint   `gorm:"default:0" json:"likes"`
	PhotoURL string `json:"photo_url"`
}
