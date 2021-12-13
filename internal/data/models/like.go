package models

// Like model
type Like struct {
	Id     int64  `db:"id"`
	FromId string `db:"from_id"`
	ToId   string `db:"to_id"`
	Value  bool   `db:"value"`
}
