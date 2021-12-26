package subj

type Supplier struct {
	ID           int
	OrderIDs     []int
	Login, Email string
	HomeAddress  string
	PassHash     string
}

func (u User) NewSupplier() {
	//TODO Написать реализацию после написания кода БД
	//u.OrderIDs = append(u.OrderIDs)
}
