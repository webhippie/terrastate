package basicauth

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/webhippie/terrastate/pkg/config"
)

// Basicauth integrates a simple basic authentication.
func Basicauth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Terrastate"`)

		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

		if len(s) != 2 {
			http.Error(w, http.StatusText(http.StatusUnauthorized), 401)
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])

		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		pair := strings.SplitN(string(b), ":", 2)

		if len(pair) != 2 {
			http.Error(w, http.StatusText(http.StatusUnauthorized), 401)
			return
		}

		if pair[0] != config.General.Username {
			http.Error(w, http.StatusText(http.StatusUnauthorized), 401)
			return
		}

		if pair[1] != config.General.Password {
			http.Error(w, http.StatusText(http.StatusUnauthorized), 401)
			return
		}

		next.ServeHTTP(w, r)
	})
}
