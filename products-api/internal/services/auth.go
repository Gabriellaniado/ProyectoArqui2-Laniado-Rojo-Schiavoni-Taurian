package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type AuthServiceImpl struct {
	usersAPIURL string
	httpClient  *http.Client
}

// NewAuthService crea una nueva instancia del servicio de autenticación
func NewAuthService(usersAPIURL string) *AuthServiceImpl {
	return &AuthServiceImpl{
		usersAPIURL: usersAPIURL,
		httpClient:  &http.Client{},
	}
}

// VerifyToken verifica un token llamando al users-api
func (s *AuthServiceImpl) VerifyToken(ctx context.Context, token string) error {
	if strings.TrimSpace(token) == "" {
		return errors.New("token is required")
	}

	url := fmt.Sprintf("%s/auth/verify-token", s.usersAPIURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error calling users-api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("invalid or expired token")
	}

	return nil
}

// VerifyAdminToken verifica un token de admin llamando al users-api
func (s *AuthServiceImpl) VerifyAdminToken(ctx context.Context, token string) error {
	if strings.TrimSpace(token) == "" {
		return errors.New("token is required")
	}

	url := fmt.Sprintf("%s/auth/verify-admin-token", s.usersAPIURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error calling users-api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("invalid admin token or insufficient permissions")
	}

	return nil
}
