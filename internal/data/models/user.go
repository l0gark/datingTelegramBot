package models

// User model
type User struct {
	Id          string `db:"id"`
	Name        string `db:"name"`
	Sex         bool   `db:"sex"` // True if Sex is MALE, False if Sex is FEMALE
	Age         int    `db:"age"`
	Description string `db:"description"`
	City        string `db:"city"`
	Image       string `db:"image"`
}
