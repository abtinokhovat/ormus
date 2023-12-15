package service

import (
	"time"

	"github.com/ormushq/ormus/manager/entity"
	"github.com/ormushq/ormus/param"
	"github.com/ormushq/ormus/pkg/errmsg"
	"github.com/ormushq/ormus/pkg/password"
	"github.com/ormushq/ormus/pkg/richerror"
)

func (s Service) Register(req param.RegisterRequest) (param.RegisterResponse, error) {
	// fetch user to check if exists before user creation
	existing, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return param.RegisterResponse{}, richerror.New("register.repo").WhitWarpError(err)
	}
	if existing != nil {
		return param.RegisterResponse{}, richerror.New("register").WhitMessage(errmsg.ErrAuthUserExisting)
	}

	hashedPassword, err := password.HashPassword(req.Password)
	if err != nil {
		return param.RegisterResponse{}, richerror.New("register.hash").WhitWarpError(err)
	}

	user := entity.User{
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		DeletedAt: nil,
		Email:     req.Email,
		Password:  hashedPassword,
		IsActive:  false,
	}

	createdUser, err := s.repo.Register(user)
	if err != nil {
		return param.RegisterResponse{}, richerror.New("register.repo").WhitWarpError(err)
	}

	// return create new user
	return param.RegisterResponse{
		Email: createdUser.Email,
	}, nil
}
