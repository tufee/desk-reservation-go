package middleware

import (
	"net/http"

	"github.com/tufee/desk-reservation-go/internal/utils"
	pkg "github.com/tufee/desk-reservation-go/pkg/utils"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := pkg.ExtractToken(w, r)
		if token == nil {
			return
		}

		ctx := r.Context()
		ctx = utils.SetContextValue(ctx, utils.AuthUserKey, token.UserId)
		ctx = utils.SetContextValue(ctx, utils.AuthEmailKey, token.Email)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
