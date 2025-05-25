package users

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// const tokenTTL = 72 * time.Hour
const tokenTTL = 5 * time.Minute

type Service interface {
	Signup(ctx context.Context, req *Signup) error
	Login(ctx context.Context, req *Login) (*Token, error)
}

type service struct {
	repo   Repository
	jwtKey []byte
	logger Logger
}

func NewService(r Repository, key []byte, logger Logger) Service {
	return &service{
		repo:   r,
		jwtKey: key,
		logger: logger,
	}
}

func (s *service) Signup(ctx context.Context, req *Signup) error {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return err
	}

	if user != nil {
		return errors.New("user already exists")
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	now := time.Now().Unix()

	user = &User{
		ID:        uuid.New().String(),
		Email:     req.Email,
		Password:  string(hash),
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
	}
	return s.repo.Create(ctx, user)
}

func (s *service) Login(ctx context.Context, req *Login) (*Token, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil || !checkPassword(req.Password, user.Password) {
		return nil, err
	}

	accessToken, err := generateJWT(user.ID, s.jwtKey)
	if err != nil {
		return nil, err
	}

	token := &Token{
		AccessToken:  accessToken,
		RefreshToken: "todo",
		ExpiresIn:    int64(tokenTTL.Seconds()),
	}
	return token, nil
}

func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateJWT(userID string, jwtKey []byte) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(tokenTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
