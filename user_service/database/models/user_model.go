package models

type User struct {
	ID       int    `db:"id" json:"id"`
	Name     string `db:"name" json:"name" binding:"required"`
	Password string `db:"password" json:"password" binding:"required"`
}