package fake

import (
	"encoding/json"
	"errors"
	"fullstack-shop/pkg/kernel/model/subj"
	"os"

	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	Users []*subj.User
}

func NewUserRepo() *UserRepo {
	repo := UserRepo{}
	p1, _ := bcrypt.GenerateFromPassword([]byte("200303"), bcrypt.DefaultCost)
	p2, _ := bcrypt.GenerateFromPassword([]byte("20030307"), bcrypt.DefaultCost)
	repo.Users = []*subj.User{
		subj.NewUser("ex@g.com", string(p1)),
		subj.NewUser("e@gma.com", string(p2)),
	}
	return &repo
}

func (u *UserRepo) GetByEmail(email string) (*subj.User, error) {
	for _, user := range u.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("User not found")
}

func (u *UserRepo) GetById(Id int) (*subj.User, error) {
	for _, user := range u.Users {
		if user.ID == Id {
			return user, nil
		}
	}
	return nil, errors.New("User not found")
}

func (u *UserRepo) Save() error {
	file, err := os.Create("userRepo.json")
	if err != nil {
		return err
	}
	err = json.NewEncoder(file).Encode(u.Users)
	if err != nil {
		return err
	}
	return nil
}

func Load() (*UserRepo, error) {
	file, err := os.Open("userRepo.json")
	if err != nil {
		return nil, err
	}
	var userRepo UserRepo
	err = json.NewDecoder(file).Decode(&userRepo.Users)
	if err != nil {
		return nil, err
	}
	return &userRepo, nil
}
