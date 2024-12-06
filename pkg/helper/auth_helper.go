package helper

import (
	"context"
	"errors"
	"github.com/TrinityKnights/Backend/pkg/jwt"
)

func (h *ContextHelper) GetJWTClaims(ctx context.Context) (*jwt.JWTClaims, error) {
	claims, ok := ctx.Value("claims").(*jwt.JWTClaims)
	if !ok || claims == nil {
		return nil, errors.New("unauthorized: invalid or missing JWT claims")
	}
	return claims, nil
}

func (h *ContextHelper) VerifyOwnership(ctx context.Context, resourceOwnerID string) error {
	claims, err := h.GetJWTClaims(ctx)
	if err != nil {
		return err
	}

	if claims.Role == RoleAdmin || claims.UserID == resourceOwnerID {
		return nil
	}

	return errors.New("forbidden: user does not have permission to access this resource")
}

func (h *ContextHelper) IsAdmin(ctx context.Context) bool {
	claims, err := h.GetJWTClaims(ctx)
	if err != nil {
		return false
	}
	return claims.Role == RoleAdmin
}
