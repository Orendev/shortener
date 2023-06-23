package middlewares

import (
	"context"
	"errors"
	"fmt"
	"github.com/Orendev/shortener/internal/auth"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

const (
	bearer       string = "bearer"
	bearerFormat string = "Bearer %s"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// по умолчанию устанавливаем оригинальный http.ResponseWriter, http.Request как тот,
		// который будем передавать следующей функции
		ow := w
		or := r
		ctx, err := HTTPToContext(or)
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie) || errors.Is(err, auth.ErrorTokenExpired):

				if ctx != nil {
					or = or.WithContext(ctx)
				}

				ctx, err = NewSigner(or.Context())
				if err != nil {
					http.Error(ow, "server error", http.StatusInternalServerError)
					return
				}

				or, err = contextToHTTP(ow, or.WithContext(ctx))
				if err != nil {
					http.Error(ow, "server error", http.StatusInternalServerError)
					return
				}

			default:
				http.Error(ow, "server error", http.StatusInternalServerError)
				return
			}
		} else {
			or = or.WithContext(ctx)
		}

		ctx, err = newParse(or.Context())

		if err != nil {
			switch {
			case errors.Is(err, auth.ErrorTokenContextMissing):
				http.Error(ow, "Unauthorized", http.StatusUnauthorized)
			case errors.Is(err, auth.ErrorTokenExpired):
				http.Error(ow, "Unauthorized", http.StatusUnauthorized)
			default:
				http.Error(ow, "server error", http.StatusInternalServerError)
			}
			return
		}

		next.ServeHTTP(ow, or.WithContext(ctx))
	})
}

// NewSigner создаёт JWT, указывая идентификатор ключа,
func NewSigner(ctx context.Context) (context.Context, error) {

	userID, err := auth.GetAuthIdentifier(ctx)
	if err != nil {
		userID = uuid.New().String()
	}
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(auth.TokenExp)),
		},
		// собственное утверждение
		UserID: userID,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(auth.SecretKey))
	if err != nil {
		return nil, err
	}

	return context.WithValue(ctx, auth.JwtContextKey, tokenString), nil
}

func HTTPToContext(r *http.Request) (context.Context, error) {
	token, ok := extractTokenFromAuthHeader(r.Header.Get(auth.HeaderAuthorizationKey))
	if !ok {
		tokenValue, err := r.Cookie(auth.CookieAccessTokenKey)
		if err != nil {
			return nil, err
		}
		token = tokenValue.Value
	}

	return newParse(context.WithValue(r.Context(), auth.JwtContextKey, token))
}

func contextToHTTP(w http.ResponseWriter, r *http.Request) (*http.Request, error) {
	ctx := r.Context()
	tokenString, ok := ctx.Value(auth.JwtContextKey).(string)

	if !ok {
		return nil, auth.ErrorTokenContextMissing
	}

	w.Header().Add(auth.HeaderAuthorizationKey, generateAuthHeaderFromToken(tokenString))
	http.SetCookie(w, &http.Cookie{
		Name:     auth.CookieAccessTokenKey,
		Value:    tokenString,
		Path:     "/",
		MaxAge:   int(auth.TokenExp.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	return r, nil
}

func newParse(ctx context.Context) (context.Context, error) {
	claims := &auth.Claims{}
	tokenString, ok := ctx.Value(auth.JwtContextKey).(string)
	if !ok {
		return nil, auth.ErrorTokenContextMissing
	}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, auth.ErrorUnexpectedSigningMethod
			}
			return []byte(auth.SecretKey), nil
		})

	if err != nil {
		if e, ok := err.(*jwt.ValidationError); ok {
			switch {
			case e.Errors&jwt.ValidationErrorMalformed != 0:
				// Token is malformed
				return nil, auth.ErrorTokenMalformed
			case e.Errors&jwt.ValidationErrorExpired != 0:
				// Token is expired
				return context.WithValue(ctx, auth.JwtUserIDContextKey, claims.UserID), auth.ErrorTokenExpired
			case e.Errors&jwt.ValidationErrorNotValidYet != 0:
				// Token is not active yet
				return nil, auth.ErrorTokenNotActive
			case e.Inner != nil:
				// report e.Inner
				return nil, e.Inner
			}
		}
		return nil, err
	}

	if !token.Valid {
		return nil, auth.ErrorTokenInvalid
	}

	return context.WithValue(ctx, auth.JwtUserIDContextKey, claims.UserID), nil
}

func extractTokenFromAuthHeader(val string) (token string, ok bool) {
	authHeaderParts := strings.Split(val, " ")
	if len(authHeaderParts) != 2 || !strings.EqualFold(authHeaderParts[0], bearer) {
		return "", false
	}

	return authHeaderParts[1], true
}

func generateAuthHeaderFromToken(token string) string {
	return fmt.Sprintf(bearerFormat, token)
}
