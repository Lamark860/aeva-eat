package service

import (
	"errors"
	"time"

	"github.com/aeva-eat/backend/internal/model"
	"github.com/aeva-eat/backend/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserExists         = errors.New("user with this username already exists")
	ErrInvalidInvite      = errors.New("invalid or expired invite code")
	ErrInviteUsed         = errors.New("invite code has already been used")
)

type AuthService struct {
	userRepo   *repository.UserRepo
	inviteRepo *repository.InviteRepo
	jwtSecret  []byte
}

func NewAuthService(userRepo *repository.UserRepo, inviteRepo *repository.InviteRepo, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		inviteRepo: inviteRepo,
		jwtSecret:  []byte(jwtSecret),
	}
}

type TokenPair struct {
	AccessToken string `json:"access_token"`
}

func (s *AuthService) Register(username, displayName, password, inviteCode string) (*model.User, *TokenPair, error) {
	// Validate invite code
	invite, err := s.inviteRepo.GetByCode(inviteCode)
	if err != nil {
		return nil, nil, ErrInvalidInvite
	}
	if invite.UsedBy != nil {
		return nil, nil, ErrInviteUsed
	}
	if invite.ExpiresAt != nil && invite.ExpiresAt.Before(time.Now()) {
		return nil, nil, ErrInvalidInvite
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}

	user, err := s.userRepo.Create(username, displayName, string(hash))
	if err != nil {
		return nil, nil, ErrUserExists
	}

	// Mark invite as used
	_ = s.inviteRepo.MarkUsed(inviteCode, user.ID)

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, &TokenPair{AccessToken: token}, nil
}

func (s *AuthService) Login(username, password string) (*model.User, *TokenPair, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, &TokenPair{AccessToken: token}, nil
}

func (s *AuthService) GetUserByID(id int) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *AuthService) ChangePassword(userID int, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByIDWithPassword(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return ErrInvalidCredentials
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(userID, string(hash))
}

func (s *AuthService) UpdateAvatarURL(userID int, avatarURL string) error {
	return s.userRepo.UpdateAvatarURL(userID, avatarURL)
}

func (s *AuthService) generateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
