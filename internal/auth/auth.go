package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type contextKey string

// Claims structure describing the JWT Claims
type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

// Errors when working with the JWT Token.
var (
	// ErrorTokenContextMissing токен не был передан.
	ErrorTokenContextMissing = errors.New("token up for parsing was not passed through the context")

	// ErrorTokenInvalid означает, что токен не удалось проверить.
	ErrorTokenInvalid = errors.New("JWT was invalid")

	// ErrorUnexpectedSigningMethod означает, что токен был подписан с использованием неожиданного метода подписи.
	ErrorUnexpectedSigningMethod = errors.New("unexpected signing method")

	// ErrorTokenMalformed токен не был отформатирован как JWT.
	ErrorTokenMalformed = errors.New("JWT is malformed")

	// ErrorTokenExpired заголовок срока действия токена прошел.
	ErrorTokenExpired = errors.New("JWT is expired")

	// ErrorTokenNotActive Токен еще не действителен.
	ErrorTokenNotActive = errors.New("token is not valid yet")
)

// JWT Token settings.
const (
	TokenExp                          = time.Hour * 3
	SecretKey                         = "supersecretkey"
	JwtContextKey          contextKey = "JWTToken"
	JwtUserIDContextKey    contextKey = "JWTUserID"
	CookieAccessTokenKey   string     = "access_token"
	HeaderAuthorizationKey string     = "Authorization"
)

// GetAuthIdentifier get the uuid of the jwt token header Authorization user.
func GetAuthIdentifier(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(JwtUserIDContextKey).(string)
	if !ok {
		return "", ErrorTokenContextMissing
	}
	return userID, nil
}
