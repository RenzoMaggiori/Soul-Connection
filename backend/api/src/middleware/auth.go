package middleware

import (
	"fmt"
	"net/http"

	"soul-connection.com/api/src/lib"
)

type AuthProvider struct {
	ApiKey string
}

func (p *AuthProvider) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		bearer, ok := req.Header["Authorization"]
		if !ok {
			http.Error(res, "Missing Authorization header", http.StatusUnauthorized)
			return
		}
		resp, err := lib.Fetch(&http.Client{}, lib.FetchRequest{
			Method:  "GET",
			Url:     fmt.Sprintf("%s/api/employees/me", lib.ApiBaseUri),
			Body:    nil,
			Headers: map[string]string{"Authorization": bearer[0], "X-Group-Authorization": p.ApiKey},
		})
		if err != nil {
			http.Error(res, "Internal server error", http.StatusInternalServerError)
			lib.ServerLog("ERROR", err)
			return
		}
		if resp.StatusCode != 200 {
			http.Error(res, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(res, req)
	})
}
