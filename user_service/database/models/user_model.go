package models

type User struct {
	ID       int    `db:"id" json:"id"`
	Email    string `db:"email" json:"email"`
	Name     string `db:"name" json:"name" binding:"required"`
	Password string `db:"password" json:"password" binding:"required"`
}
