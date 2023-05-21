package middlewares

import (
	"context"
	"errors"
	"github.com/Orendev/shortener/internal/auth"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// по умолчанию устанавливаем оригинальный http.ResponseWriter, http.Request как тот,
		// который будем передавать следующей функции
		ow := w
		or := r
		ctx, err := cookieToContext(or)

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

				or = or.WithContext(ctx)
				cookie, err := contextToCookie(or.Context())
				if err != nil {
					http.Error(ow, "server error", http.StatusInternalServerError)
					return
				}

				http.SetCookie(ow, cookie)
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

	userId, err := auth.GetAuthIdentifier(ctx)
	if err != nil {
		userId = uuid.New().String()
	}
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(auth.TokenExp)),
		},
		// собственное утверждение
		UserID: userId,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(auth.SecretKey))
	if err != nil {
		return nil, err
	}

	return context.WithValue(ctx, auth.JwtContextKey, tokenString), nil
}

func cookieToContext(r *http.Request) (context.Context, error) {
	token, err := r.Cookie(auth.CookieAccessTokenKey)
	if err != nil {
		return nil, err
	}
	return newParse(context.WithValue(r.Context(), auth.JwtContextKey, token.Value))
}

func contextToCookie(ctx context.Context) (*http.Cookie, error) {
	tokenString, ok := ctx.Value(auth.JwtContextKey).(string)

	if !ok {
		return nil, auth.ErrorTokenContextMissing
	}

	return &http.Cookie{
		Name:     auth.CookieAccessTokenKey,
		Value:    tokenString,
		Path:     "/",
		MaxAge:   int(auth.TokenExp.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}, nil
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
