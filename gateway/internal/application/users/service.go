package users

import "github.com/myself/golang-social-media/common/domain/user"

type Service interface {
	SampleUser() user.User
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) SampleUser() user.User {
	return user.User{ID: "1", Username: "demo", FullName: "Gateway Demo"}
}
