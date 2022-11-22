package user

type User struct {
	Id       int
	Email    string
	Password string
	Nickname string
}

type UserInfo struct {
	Id       int
	Email    string
	Nickname string
	Rank     string
}

type AllUserInfo struct {
	Id         int
	Email      string
	Nickname   string
	Rank       string
	SignUpDate string
}
