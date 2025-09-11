package services

import "samsungvoicebe/repo"

type UserService struct {
	userRepo *repo.UserRepo
}

func NewUserService(userRepo *repo.UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(userID string) error {
	err := s.userRepo.CreateUser(userID)
	if err != nil {
		return err
	}
	return nil
}
