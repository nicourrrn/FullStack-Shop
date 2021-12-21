package databases

import (
	"fullstack-shop/pkg/login/model/repo"
	"fullstack-shop/pkg/login/model/subj/user"
	"os"
)

type UserDBRepo os.File

func (u UserDBRepo) GetByKey(key repo.BDKey) (*user.User, error) {
	return nil, nil
}

func (u UserDBRepo) Post(user *user.User) error {
	return nil
}
