package users

import domain "github.com/myself/golang-social-media/apps/gateway/internal/domain/user"

type Service interface {
	SampleUser() domain.User
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) SampleUser() domain.User {
	return domain.User{ID: "1", Username: "demo", FullName: "Gateway Demo"}
}
