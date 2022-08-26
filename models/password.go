package models

type Password struct {
	Id       uint `gorm:"primary_key"`
	Link     string
	Password string
}

func NewPassword(link string, password string) *Password {
	return &Password{
		Link:     link,
		Password: password,
	}
}
