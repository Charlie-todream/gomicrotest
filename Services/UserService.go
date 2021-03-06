package Services

import "errors"

type IUserService interface {
	GetName(userid int) string
	DelUser(userid int) error
}

type UserService struct {
}

func (this UserService) GetName(userid int) string {
	if userid == 101 {
		return "charlie66"
	}
	return "guest"
}

func (this UserService) DelUser(userid int) error {
	if userid == 101 {
		return errors.New("无权限")
	}
	return nil
}
