package requests

import (
	"backend/src/utils/adapter"
	"backend/src/utils/database"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type UserServer struct {
	Database *database.DBConnect
}

func NewUserServer(database *database.DBConnect) *UserServer {
	return &UserServer{Database: database}
}

const USERS_COLLECTION = "users"

func (server *UserServer) GetUser(w http.ResponseWriter, r *http.Request) {
	setSuccessHeader(w)

	path := r.URL.Path[1:]
	dirs := strings.Split(path, "/")

	if len(dirs) < 2 {
		makeResponse(w, "Bad Path")
		return
	}

	if dirs[0] != USERS_COLLECTION {
		makeResponse(w, "Bad Path")
		return
	}

	val, err := strconv.Atoi(dirs[1])

	if err != nil {
		makeResponse(w, "Bad ID")
		return
	}

	tokenBody := r.Header.Get("token")
	token := adapter.AuthToken{ID: val, Token: tokenBody}

	userDatabaseAdapter := adapter.CreateUserDatabaseAdapter(server.Database)
	ok, err := userDatabaseAdapter.CheckToken(token)

	if err != nil || !ok {
		makeResponse(w, "Bad Auth")
		return
	}

	user, err := userDatabaseAdapter.GetUser(token.ID)

	if err != nil {
		makeResponse(w, "Bad User ID")
		return
	}

	json, err := json.Marshal(user)

	if err != nil {
		makeResponse(w, "Bad JSON")
		return
	}

	response := fmt.Sprintf("{\"status\":\"Success\", \"user\" : %v}", string(json))
	w.Write([]byte(response))
}

func (server *UserServer) CreateUser(w http.ResponseWriter, r *http.Request) {
	setSuccessHeader(w)

	path := r.URL.Path[1:]
	dirs := strings.Split(path, "/")

	if len(dirs) < 2 {
		makeResponse(w, "Bad Path")
		return
	}

	if dirs[0] != USERS_COLLECTION {
		makeResponse(w, "Bad Path")
		return
	}

	if dirs[1] != CREATE_ACTION {
		makeResponse(w, "Bad Path")
		return
	}

	createUser := adapter.User{}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&createUser)

	if err != nil {
		makeResponse(w, "Bad Body")
		return
	}

	userDatabaseAdapter := adapter.CreateUserDatabaseAdapter(server.Database)
	token, err := userDatabaseAdapter.CreateUser(&createUser)

	if err != nil {
		makeResponse(w, "Bad Auth")
		return
	}

	json, err := json.Marshal(token)

	if err != nil {
		makeResponse(w, "Bad JSON")
		return
	}

	response := fmt.Sprintf("{\"status\":\"Success\", \"token\" : %v}", string(json))
	w.Write([]byte(response))
}

func (server *UserServer) GetToken(w http.ResponseWriter, r *http.Request) {
	setSuccessHeader(w)

	path := r.URL.Path[1:]

	if path != "auth" {
		makeResponse(w, "Bad Path")
		return
	}

	authData := adapter.AuthData{}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&authData)

	if err != nil {
		makeResponse(w, "Bad Body")
		return
	}

	userDatabaseAdapter := adapter.CreateUserDatabaseAdapter(server.Database)
	token, err := userDatabaseAdapter.AuthUser(authData)

	if err != nil {
		makeResponse(w, "Bad Auth")
		return
	}

	json, err := json.Marshal(token)

	if err != nil {
		makeResponse(w, "Bad JSON")
		return
	}

	response := fmt.Sprintf("{\"status\":\"Success\", \"token\" : %v}", string(json))
	w.Write([]byte(response))
}
