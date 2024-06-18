package server

import (
	"cognix.ch/api/v2/core/repository"
	"cognix.ch/api/v2/core/security"
	"cognix.ch/api/v2/core/utils"
	"context"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

const ContextParamUser = "CONTEXT_USER"

type AuthMiddleware struct {
	jwtService security.JWTService
	userRepo   repository.UserRepository
}

func NewAuthMiddleware(jwtService security.JWTService,
	userRepo repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{jwtService: jwtService,
		userRepo: userRepo}
}

func (m *AuthMiddleware) RequireAuth(c *gin.Context) {

	//Get the  bearer Token
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		handleError(c, utils.ErrorUnauthorized.New("Authorization token is required"))
		c.Abort()
		return
	}

	extractedToken := strings.Split(tokenString, "Bearer ")

	if len(extractedToken) != 2 {
		handleError(c, utils.ErrorBadRequest.New("Incorrect format of authorization token"))
		c.Abort()
		return
	}

	claims, err := m.jwtService.ParseAndValidate(strings.TrimSpace(extractedToken[1]))
	if err != nil {
		handleError(c, utils.ErrorBadRequest.New("Token is not valid"))
		c.Abort()
		return
	}

	if claims.ExpiresAt != 0 && time.Now().Unix() > claims.ExpiresAt {
		handleError(c, utils.ErrorUnauthorized.New("token expired"))
		c.Abort()
		return
	}

	if claims.User, err = m.userRepo.GetByIDAndTenantID(c.Request.Context(), claims.User.ID, claims.User.TenantID); err != nil {
		handleError(c, utils.ErrorUnauthorized.Wrap(err, "wrong user"))
		c.Abort()
		return
	}
	c.Request = c.Request.WithContext(context.WithValue(
		c.Request.Context(), ContextParamUser, claims))
	c.Next()
}

func GetContextIdentity(c *gin.Context) (*security.Identity, error) {
	claims, ok := c.Request.Context().Value(ContextParamUser).(*security.Identity)
	if !ok {
		return nil, utils.ErrorPermission.New("broken session")
	}
	return claims, nil
}
