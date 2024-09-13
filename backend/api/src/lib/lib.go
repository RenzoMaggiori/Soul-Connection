package lib

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const ApiBaseUri string = "https://soul-connection.fr"

func GetIdFromRequest(req *http.Request, key string) (int, error) {
	vars := mux.Vars(req)
	stringValue, ok := vars[key]

	if !ok {
		return 0, errors.New("could not find ID")
	}

	value, err := strconv.Atoi(stringValue)
	if err != nil {
		ServerLog("ERROR", err)
		return 0, err
	}

	return value, nil
}
