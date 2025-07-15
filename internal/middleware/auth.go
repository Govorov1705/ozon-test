package middleware

import (
	"context"
	"strings"

	"github.com/Govorov1705/ozon-test/internal/jwt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type contextKey string

const userIDKey contextKey = "userID"

func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	return userID, ok
}

func Auth(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.Next()
		return
	}

	authHeaderParts := strings.Fields(authHeader)
	if len(authHeaderParts) != 2 {
		c.Next()
		return
	}

	token := authHeaderParts[1]
	claims, err := jwt.ValidateJWT(token)
	if err != nil {
		c.Next()
		return
	}

	userIDStr := claims["sub"].(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.Next()
		return
	}

	ctx := WithUserID(c.Request.Context(), userID)
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}
