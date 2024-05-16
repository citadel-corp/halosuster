package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/citadel-corp/halosuster/internal/common/id"
	"github.com/citadel-corp/halosuster/internal/common/jwt"
	"github.com/citadel-corp/halosuster/internal/common/password"
)

type Service interface {
	CreateITUser(ctx context.Context, req CreateITUserPayload) (*UserResponse, error)
	CreateNurseUser(ctx context.Context, req CreateNurseUserPayload) (*UserResponse, error)
	LoginITUser(ctx context.Context, req ITUserLoginPayload) (*UserResponse, error)
	LoginNurseUser(ctx context.Context, req NurseUserLoginPayload) (*UserResponse, error)
	UpdateNurse(ctx context.Context, userID string, req UpdateNursePayload) error
	DeleteNurse(ctx context.Context, userID string) error
	GrantNurseAccess(ctx context.Context, userID string, req GrantNurseAccessPayload) error
}

type userService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &userService{repository: repository}
}

func (s *userService) CreateITUser(ctx context.Context, req CreateITUserPayload) (*UserResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		return nil, err
	}
	user := &User{
		ID:             id.GenerateStringID(16),
		NIP:            req.NIP,
		Name:           req.Name,
		UserType:       IT,
		HashedPassword: &hashedPassword,
	}
	err = s.repository.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	// create access token with signed jwt
	accessToken, err := jwt.Sign(time.Hour*2, string(IT), fmt.Sprint(user.ID))
	if err != nil {
		return nil, err
	}
	return &UserResponse{
		UserID:      user.ID,
		NIP:         user.NIP,
		Name:        req.Name,
		AccessToken: &accessToken,
	}, nil
}

func (s *userService) CreateNurseUser(ctx context.Context, req CreateNurseUserPayload) (*UserResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	user := &User{
		ID:              id.GenerateStringID(16),
		NIP:             req.NIP,
		Name:            req.Name,
		UserType:        Nurse,
		IdentityCardURL: &req.IdentityCardScanImg,
	}
	err = s.repository.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return &UserResponse{
		UserID: user.ID,
		NIP:    user.NIP,
		Name:   req.Name,
	}, nil
}

func (s *userService) LoginITUser(ctx context.Context, req ITUserLoginPayload) (*UserResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	user, err := s.repository.GetByNIP(ctx, req.NIP)
	if err != nil {
		return nil, err
	}
	match, err := password.Matches(req.Password, *user.HashedPassword)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, ErrWrongPassword
	}
	// create access token with signed jwt
	accessToken, err := jwt.Sign(time.Hour*2, string(IT), fmt.Sprint(user.ID))
	if err != nil {
		return nil, err
	}
	return &UserResponse{
		UserID:      user.ID,
		NIP:         user.NIP,
		Name:        user.Name,
		AccessToken: &accessToken,
	}, nil
}

// LoginNurseUser implements Service.
func (s *userService) LoginNurseUser(ctx context.Context, req NurseUserLoginPayload) (*UserResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	user, err := s.repository.GetByNIP(ctx, req.NIP)
	if err != nil {
		return nil, err
	}
	if user.HashedPassword == nil {
		return nil, ErrPasswordNotCreated
	}
	match, err := password.Matches(req.Password, *user.HashedPassword)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, ErrWrongPassword
	}
	// create access token with signed jwt
	accessToken, err := jwt.Sign(time.Hour*2, string(Nurse), fmt.Sprint(user.ID))
	if err != nil {
		return nil, err
	}
	return &UserResponse{
		UserID:      user.ID,
		NIP:         user.NIP,
		Name:        user.Name,
		AccessToken: &accessToken,
	}, nil
}

// UpdateNurse implements Service.
func (s *userService) UpdateNurse(ctx context.Context, userID string, req UpdateNursePayload) error {
	err := req.Validate()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	user, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user.UserType != Nurse {
		return ErrUserNotFound
	}
	_, err = s.repository.GetByNIP(ctx, req.NIP)
	if errors.Is(err, ErrUserNotFound) {
		user.NIP = req.NIP
		user.Name = req.Name
		return s.repository.Update(ctx, user)
	}
	if err != nil {
		return err
	}
	return ErrNIPAlreadyExists
}

// DeleteNurse implements Service.
func (s *userService) DeleteNurse(ctx context.Context, userID string) error {
	user, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user.UserType != Nurse {
		return ErrUserNotFound
	}
	return s.repository.DeleteByID(ctx, userID)
}

// GrantNurseAccess implements Service.
func (s *userService) GrantNurseAccess(ctx context.Context, userID string, req GrantNurseAccessPayload) error {
	err := req.Validate()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	user, err := s.repository.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user.UserType != Nurse {
		return ErrUserNotFound
	}
	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		return err
	}
	user.HashedPassword = &hashedPassword
	return s.repository.Update(ctx, user)
}
