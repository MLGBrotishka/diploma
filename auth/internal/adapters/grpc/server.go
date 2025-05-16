package grpc_server

import (
	"context"
	"errors"

	"auth/internal/entity"
	desc "auth/pkg/api/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, login, password string) (string, error)
	Register(ctx context.Context, login, password string) (int64, error)
	CheckPermission(ctx context.Context, userId int64, permission entity.Permission) (bool, error)
	Logout(ctx context.Context, token string) error
}

type Service struct {
	desc.UnimplementedAuthServer
	auth Auth
}

func New(auth Auth) *Service {
	return &Service{
		auth: auth,
	}
}

func (s *Service) Login(
	ctx context.Context,
	in *desc.LoginRequest,
) (*desc.LoginResponse, error) {
	if in.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	token, err := s.auth.Login(ctx, in.GetLogin(), in.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &desc.LoginResponse{Token: token}, nil
}

func (s *Service) Register(
	ctx context.Context,
	in *desc.RegisterRequest,
) (*desc.RegisterResponse, error) {
	if in.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	uid, err := s.auth.Register(ctx, in.GetLogin(), in.GetPassword())
	if err != nil {
		if errors.Is(err, entity.ErrLoginAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &desc.RegisterResponse{UserId: uid}, nil
}

func (s *Service) CheckPermission(
	ctx context.Context,
	in *desc.PermissionRequest,
) (*desc.PermissionResponse, error) {
	if in.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	if in.Permission == 0 {
		return nil, status.Error(codes.InvalidArgument, "permission is required")
	}

	permission := convertToEntityPermission(in.GetPermission())
	if permission == entity.PermissionNone {
		return nil, status.Error(codes.InvalidArgument, "invalid permission")
	}

	havePermission, err := s.auth.CheckPermission(ctx, in.GetUserId(), permission)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to check permission")
	}

	return &desc.PermissionResponse{HavePermission: havePermission}, nil
}

func convertToEntityPermission(perm desc.Permission) entity.Permission {
	switch perm {
	case desc.Permission_PERMISSION_NONE:
		return entity.PermissionNone
	case desc.Permission_PERMISSION_CREATE:
		return entity.PermissionCreate
	case desc.Permission_PERMISSION_APPLY:
		return entity.PermissionApply
	case desc.Permission_PERMISSION_ROLLBACK:
		return entity.PermissionRollback
	case desc.Permission_PERMISSION_LIST:
		return entity.PermissionList
	case desc.Permission_PERMISSION_GET:
		return entity.PermissionGet
	case desc.Permission_PERMISSION_APPLY_OTHER:
		return entity.PermissionApplyOther
	case desc.Permission_PERMISSION_ROLLBACK_OTHER:
		return entity.PermissionRollbackOther
	}
	return entity.PermissionNone
}

func (s *Service) Logout(
	ctx context.Context,
	in *desc.LogoutRequest,
) (*desc.LogoutResponse, error) {
	if in.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	err := s.auth.Logout(ctx, in.GetToken())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to logout")
	}

	return &desc.LogoutResponse{Success: true}, nil
}
