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
)

type AuthService struct {
	userRepo  *repository.UserRepo
	jwtSecret []byte
}

func NewAuthService(userRepo *repository.UserRepo, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

type TokenPair struct {
	AccessToken string `json:"access_token"`
}

func (s *AuthService) Register(username, displayName, password string) (*model.User, *TokenPair, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}

	user, err := s.userRepo.Create(username, displayName, string(hash))
	if err != nil {
		return nil, nil, ErrUserExists
	}

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

func (s *AuthService) generateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
