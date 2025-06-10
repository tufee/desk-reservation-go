package api

import (
	"net/http"

	"github.com/tufee/desk-reservation-go/internal/domain"
	"github.com/tufee/desk-reservation-go/internal/service"
	"github.com/tufee/desk-reservation-go/internal/utils"
	pkg "github.com/tufee/desk-reservation-go/pkg/utils"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials domain.Credentials

	if err := pkg.ParseAndValidateRequest(r, &credentials, w); err != nil {
		return
	}

	credentials = buildCredentialsFromRequest(credentials)
	ctx := utils.SetContextValue(r.Context(), utils.LoginKey, credentials)

	token, err := service.LoginService(ctx)
	if err != nil {
		pkg.HandleHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"token": "` + token.Token + `"}`))
}

func buildCredentialsFromRequest(data domain.Credentials) domain.Credentials {
	return domain.Credentials{
		Email:    data.Email,
		Password: data.Password,
	}
}
