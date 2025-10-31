package controllers

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthService interface {
	VerifyToken(ctx context.Context, token string) error
	VerifyAdminToken(ctx context.Context, token string) error
}

// AuthController maneja la autenticación sin depender de otros servicios
type AuthController struct {
	service AuthService
}

// NewAuthController crea una nueva instancia del controller de autenticación
func NewAuthController(authService AuthService) *AuthController {
	return &AuthController{
		service: authService,
	}
}

func (c *AuthController) VerifyToken(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
		ctx.Abort()
		return
	}

	tokenString := strings.TrimPrefix(token, "Bearer ")
	if tokenString == token {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		ctx.Abort()
		return
	}

	if err := c.service.VerifyToken(ctx.Request.Context(), tokenString); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Invalid token",
			"details": err.Error(),
		})
		ctx.Abort()
		return
	}

	ctx.Next()
}

func (c *AuthController) VerifyAdminToken(ctx *gin.Context) {
	// Recibir el token desde el header de la request
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
		ctx.Abort()
		return
	}

	// Limpiar el prefijo "Bearer " ---
	tokenString := strings.TrimPrefix(token, "Bearer ")
	if tokenString == token {
		// Si no había prefijo "Bearer ", el token está mal formado o no es Bearer
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		ctx.Abort()
		return
	}
	// Llamar al servicio de verify admin token
	err := c.service.VerifyAdminToken(ctx.Request.Context(), tokenString)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid admin token or insufficient permissions"})
		ctx.Abort()
		return
	}
	// Token válido, continuar con la siguiente función
	ctx.Next()
}
