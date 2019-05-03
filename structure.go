package main

type User struct {
	ID       int
	Username string
	Type     string
	Login    bool
}

func newUser(id int, name, utype string) User {
	return User{id, name, utype, true}
}

func nilUser() User {
	var user User
	user.Login = false
	return user
}
