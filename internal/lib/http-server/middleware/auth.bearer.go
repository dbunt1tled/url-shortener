package middleware

import (
	"go_first/internal/lib/common/jwt"
	"go_first/internal/lib/common/model/apierror"
	"go_first/internal/lib/common/model/user"
	"go_first/storage"
	"go_first/storage/mysql"
	"net/http"
	"strings"
)

const BearerSchema = "Bearer"

func AuthBearer(storage storage.Storage) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				apierror.ErrorUnauthorizedString(w, "Unauthorized", 40100001)
				return
			}
			authToken := strings.TrimSpace(strings.Split(authHeader, BearerSchema)[1])
			if len(authToken) < 2 {
				apierror.ErrorUnauthorizedString(w, "Unauthorized", 40100002)
			}
			token, err := jwt.JWToken{}.Decode(authToken, true)
			if err != nil {
				apierror.ErrorUnauthorizedString(w, "Unauthorized", 40100003)
			}
			u, err := storage.GetUser(mysql.UserFilter{ID: int64(token["iss"].(float64))})
			if err != nil {
				apierror.ErrorUnauthorizedString(w, "Unauthorized", 40100004)
			}

			if u.Status != user.StatusActive {
				apierror.ErrorUnauthorizedString(w, "Unauthorized", 40100005)
			}
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}

}
