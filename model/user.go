package model

type User struct {
	Id string `xorm:"pk" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c User) TableName() string {
	return "users"
}
