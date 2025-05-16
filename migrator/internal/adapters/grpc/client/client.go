package client

import (
	"context"
	"fmt"

	desc "migrator/pkg/api/auth"
)

type authWrapper struct {
	auth desc.AuthClient
}

func New(auth desc.AuthClient) *authWrapper {
	return &authWrapper{auth: auth}
}

func (a *authWrapper) CheckPermissionCreate(ctx context.Context, userID int64) (bool, error) {
	resp, err := a.auth.CheckPermission(ctx, &desc.PermissionRequest{
		UserId:     userID,
		Permission: desc.Permission_PERMISSION_CREATE,
	})
	if err != nil {
		return false, fmt.Errorf("a.auth.CheckPermission: %w", err)
	}
	return resp.GetHavePermission(), nil
}

func (a *authWrapper) CheckPermissionApply(ctx context.Context, userID int64) (bool, error) {
	resp, err := a.auth.CheckPermission(ctx, &desc.PermissionRequest{
		UserId:     userID,
		Permission: desc.Permission_PERMISSION_APPLY,
	})
	if err != nil {
		return false, fmt.Errorf("a.auth.CheckPermission: %w", err)
	}
	return resp.GetHavePermission(), nil
}

func (a *authWrapper) CheckPermissionRollback(ctx context.Context, userID int64) (bool, error) {
	resp, err := a.auth.CheckPermission(ctx, &desc.PermissionRequest{
		UserId:     userID,
		Permission: desc.Permission_PERMISSION_ROLLBACK,
	})
	if err != nil {
		return false, fmt.Errorf("a.auth.CheckPermission: %w", err)
	}
	return resp.GetHavePermission(), nil
}

func (a *authWrapper) CheckPermissionList(ctx context.Context, userID int64) (bool, error) {
	resp, err := a.auth.CheckPermission(ctx, &desc.PermissionRequest{
		UserId:     userID,
		Permission: desc.Permission_PERMISSION_LIST,
	})
	if err != nil {
		return false, fmt.Errorf("a.auth.CheckPermission: %w", err)
	}
	return resp.GetHavePermission(), nil
}

func (a *authWrapper) CheckPermissionGet(ctx context.Context, userID int64) (bool, error) {
	resp, err := a.auth.CheckPermission(ctx, &desc.PermissionRequest{
		UserId:     userID,
		Permission: desc.Permission_PERMISSION_GET,
	})
	if err != nil {
		return false, fmt.Errorf("a.auth.CheckPermission: %w", err)
	}
	return resp.GetHavePermission(), nil
}

func (a *authWrapper) CheckPermissionApplyOther(ctx context.Context, userID int64) (bool, error) {
	resp, err := a.auth.CheckPermission(ctx, &desc.PermissionRequest{
		UserId:     userID,
		Permission: desc.Permission_PERMISSION_APPLY_OTHER,
	})
	if err != nil {
		return false, fmt.Errorf("a.auth.CheckPermission: %w", err)
	}
	return resp.GetHavePermission(), nil
}

func (a *authWrapper) CheckPermissionRollbackOther(ctx context.Context, userID int64) (bool, error) {
	resp, err := a.auth.CheckPermission(ctx, &desc.PermissionRequest{
		UserId:     userID,
		Permission: desc.Permission_PERMISSION_ROLLBACK_OTHER,
	})
	if err != nil {
		return false, fmt.Errorf("a.auth.CheckPermission: %w", err)
	}
	return resp.GetHavePermission(), nil
}
