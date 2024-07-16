package utils

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// get a path param from request
func GetPathParam(r *http.Request, name string) (string, error) {
	vars := mux.Vars(r)
	value, ok := vars[name]
	if !ok {
		return "", fmt.Errorf("invalid or missing %s in request param", name)
	}
	return value, nil
}

func ParseQueryParams(r *http.Request, repoSearchParams *RepositorySearchParams) {
	query := r.URL.Query()
	if query.Get("name") != "" {
		repoSearchParams.Name = query.Get("name")
	}
	if query.Get("language") != "" {
		repoSearchParams.Language = query.Get("language")
	}
	if query.Get("top_stars") != "" {
		repoSearchParams.TopStarsCount, _ = strconv.Atoi(query.Get("top_stars"))
	}
}
