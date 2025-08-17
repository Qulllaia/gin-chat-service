package models

type User struct {
	ID       int    `db:"id" json:"id"`
	Name     string `db:"email" json:"email" binding:"required"`
	Password string `db:"password" json:"password" binding:"required"`
}