package model

type Password struct {
	Id       int64  `gorm:"primaryKey;autoIncrement;column:id"`
	Link     string `gorm:"column:link;unique"`
	Password string `gorm:"column:password;unique"`
}

func NewPassword(link string, password string) *Password {
	return &Password{
		Link:     link,
		Password: password,
	}
}

func (Password) TableName() string {
	return "tbl_passwords"
}
