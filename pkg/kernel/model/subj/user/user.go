package user

type User struct {
	ID           int
	OrderIDs     []int
	Login, Email string
	HomeAddress  string
	PassHash     string
}

func (u User) NewOrder() {
	//TODO Написать реализацию после написания кода БД
	//u.OrderIDs = append(u.OrderIDs)
}
