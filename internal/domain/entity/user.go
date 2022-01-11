package entity

type User struct {
	FullName  string `json:"fullName" bson:"FullName"`
	UserName  string `json:"userName" bson:"UserName"`
	UserEmail string `json:"userEmail" bson:"UserEmail"`
	UserId    string `json:"userId" bson:"UserId"`
}

func NewUser(fullName, userName, userEmail, userId string) *User {
	return &User{
		FullName:  fullName,
		UserName:  userName,
		UserEmail: userEmail,
		UserId:    userId,
	}
}
