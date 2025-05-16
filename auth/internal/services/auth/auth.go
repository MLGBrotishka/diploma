package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"auth/internal/entity"

	"golang.org/x/crypto/bcrypt"
)

type authService interface {
	Login(ctx context.Context, login, password string) (string, error)
	Register(ctx context.Context, login, password string) (int64, error)
	CheckPermission(ctx context.Context, userId int64, permission entity.Permission) (bool, error)
	Logout(ctx context.Context, token string) error
}

var _ authService = (*Auth)(nil)

type authRepo interface {
	GetUserByLogin(ctx context.Context, login string) (entity.User, error)
	SaveUser(ctx context.Context, login string, passwordHash []byte) (int64, error)
	CheckUserPermission(ctx context.Context, userID int64, permission entity.Permission) (bool, error)
	SetUserInactive(ctx context.Context, userID int64) error
}

type tokenProvider interface {
	NewToken(user entity.User, duration time.Duration) (string, error)
	ParseToken(token string) (entity.User, error)
}

// Auth - сервис аутентификации и авторизации.
type Auth struct {
	authRepo      authRepo
	tokenProvider tokenProvider
	tokenTTL      time.Duration
}

// New - конструктор сервиса аутентификации и авторизации.
func New(
	authRepo authRepo,
	tokenProvider tokenProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		tokenTTL:      tokenTTL,
		authRepo:      authRepo,
		tokenProvider: tokenProvider,
	}
}

// Login проверяет учетные данные пользователя и возвращает токен доступа.
//
// Если пользователь существует, но пароль неверный, возвращает ошибку.
// Если пользователь не существует, возвращает ошибку.
// Аргументы:
//
//	ctx: context.Context - Контекст запроса.
//	login: string - Логин пользователя.
//	password: string - Пароль пользователя.
//
// Возвращает:
//
//	string: Токен доступа.
//	error: Ошибка, если таковая имеется (например, неверные учетные данные).
func (a *Auth) Login(
	ctx context.Context,
	login string,
	password string,
) (string, error) {
	user, err := a.authRepo.GetUserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return "", fmt.Errorf("a.authRepo.GetUserByLogin: %w", entity.ErrInvalidCredentials)
		}
		return "", fmt.Errorf("a.authRepo.GetUserByLogin: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		return "", fmt.Errorf("bcrypt.CompareHashAndPassword: %w", entity.ErrInvalidCredentials)
	}

	token, err := a.tokenProvider.NewToken(user, a.tokenTTL)
	if err != nil {
		return "", fmt.Errorf("a.tokenProvider.NewToken: %w", err)
	}

	return token, nil
}

// Register регистрирует нового пользователя в системе и возвращает его ID.
// Если пользователь с данным логином уже существует, возвращает ошибку ErrLoginAlreadyExists.
// Аргументы:
//
//	ctx: context.Context - Контекст запроса.
//	login: string - Логин нового пользователя.
//	pass: string - Пароль нового пользователя.
//
// Возвращает:
//
//	int64: Уникальный идентификатор созданного пользователя.
//	error: Ошибка, если таковая имеется (например, логин уже существует).
func (a *Auth) Register(ctx context.Context, login string, pass string) (int64, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("bcrypt.GenerateFromPassword: %w", err)
	}

	id, err := a.authRepo.SaveUser(ctx, login, passHash)
	if err != nil {
		if errors.Is(err, entity.ErrLoginAlreadyExists) {
			return 0, fmt.Errorf("a.authRepo.SaveUser: %w", entity.ErrLoginAlreadyExists)
		}
		return 0, fmt.Errorf("a.authRepo.SaveUser: %w", err)
	}

	return id, nil
}

// CheckPermission проверяет, имеет ли пользователь определенное разрешение.
// Аргументы:
//
//	ctx: context.Context - Контекст запроса.
//	userID: int64 - Идентификатор пользователя.
//	permission: entity.Permission - Проверяемое разрешение.
//
// Возвращает:
//
//	bool: true, если пользователь имеет разрешение, false в противном случае.
//	error: Ошибка, если таковая имеется.
func (a *Auth) CheckPermission(ctx context.Context, userID int64, permission entity.Permission) (bool, error) {
	allowed, err := a.authRepo.CheckUserPermission(ctx, userID, permission)
	if err != nil {
		return false, fmt.Errorf("a.authRepo.CheckUserPermission: %w", err)
	}
	return allowed, nil
}

// Logout делает сессию пользователя недействительным.
// Аргументы:
//
//	ctx: context.Context - Контекст запроса.
//	token: string - Токен доступа пользователя для инвалидации.
//
// Возвращает:
//
//	error: Ошибка, если таковая имеется. Не возвращает ошибку, если токен уже недействителен.
func (a *Auth) Logout(ctx context.Context, token string) error {
	user, err := a.tokenProvider.ParseToken(token)
	if err != nil {
		return nil
	}

	err = a.authRepo.SetUserInactive(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("a.authRepo.SetUserInactive: %w", err)
	}

	return nil
}
