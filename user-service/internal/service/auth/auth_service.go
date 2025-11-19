package auth

import (
	"errors"
	"fmt"
	"log/slog"
	"wog-server/domain"
	"wog-server/internal/db"
	"wog-server/internal/lib/hash"
	"wog-server/internal/lib/jwt"
	"wog-server/internal/logger"
)

type AuthService struct {
	db db.IUserRepo
	jwt *jwt.JWTSrv
}

func NewAuthService(dbuser db.IUserRepo, jwtSrv *jwt.JWTSrv) *AuthService {
	return &AuthService{db: dbuser, jwt: jwtSrv}
}

func (r *AuthService) checkFieldUniqueness(field, value, errMsg string) error {
	exists, err := r.db.IsFieldValueExists(field, value)
	if err != nil {
		return fmt.Errorf("checking %s exists error: %w", field, err)
	}
	if exists {
		return fmt.Errorf("%s", errMsg)
	}
	return nil
}

func (r *AuthService) Register(user domain.User) (userID string, err error) {
	if err := r.checkFieldUniqueness("email", user.Email, "this email is already taken"); err != nil {
		logger.Log.Error(err.Error())
		return "", err
	}
	if err := r.checkFieldUniqueness("name", user.Name, "this name is already taken"); err != nil {
		logger.Log.Error(err.Error())
		return "", err
	}
	
	hashPswrd, err := hash.HashPassword(user.Password)
	if err != nil {
		logger.Log.Error(err.Error())
	}
	user.Password = hashPswrd

	userID, err = r.db.Insert(user)
	if err != nil {
		logger.Log.Error("Error creating user in database", slog.String("err", err.Error()))
		return "", db.ErrCreatingUser
	}
	return userID, nil
}

func (r *AuthService) Login(input domain.User) (*domain.User, error) {
	var (
		user *domain.DBUser
		err  error
	)
	
	if input.Email != "" {
		user, err = r.db.Find("email", input.Email)
	} else if input.Name != "" {
		user, err = r.db.Find("name", input.Name)
	}

	if err != nil {
		logger.Log.Error("Ошибка получения юзера из бд", slog.String("err", err.Error()))
		return nil, errors.New("Ошибка получения юзера из бд")
	}
	if user == nil {
		return nil, errors.New("Юзер не найден")
	}
	
	logger.Log.Debug(user.Password)
	logger.Log.Debug(user.Email)

	if isCorrect := hash.CheckPassword(input.Password, user.Password); !isCorrect {
		return nil, fmt.Errorf("Password is incorrect")
	}

	token, err := r.jwt.Generate(user.GetStrId())
	if err != nil {
		logger.Log.Error("Ошибка генерации jwt токена", slog.String("err", err.Error()))
		return nil, errors.New("Ошибка генерации jwt токена")
	}
	res := domain.User{
		Token: token,
		Id: user.GetStrId(),
		Name: user.Name,
		Email: user.Email,
	}

	return &res, nil
}

