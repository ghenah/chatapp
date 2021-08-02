package dsgorm

import "time"

type User struct {
	ID       uint      `gorm:"not null;primaryKey;autoIncrement"`
	Username string    `gorm:"not null;size:16;unique"`
	Picture  string    `gorm:"not null"`
	Email    string    `gorm:"not null;size:255;unique"`
	Password string    `gorm:"not null"`
	RegDate  time.Time `gorm:"not null;autoCreateTime"`
	Friends  []*User   `gorm:"many2many:friends"`
	Ignored  []*User   `gorm:"many2many:ignored"`
}
