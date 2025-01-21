package models

type UserDetails struct {
	ID         int64 `gorm:"primaryKey"`
	FirstName  string
	LastName   string
	Email      string
	Age        string
	Gender     string
	Department string
	Company    string
	Salary     string
	DateJoined string
	IsActive   string
}
