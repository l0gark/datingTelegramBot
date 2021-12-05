package models

// User model
type User struct {
	Id          string `db:"id"`
	Name        string `db:"name"`
	Sex         string `db:"sex"`
	Age         int    `db:"age"`
	Description string `db:"description"`
	City        string `db:"city"`
	Image       string `db:"image"`
}
