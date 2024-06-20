package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"myRestAPIWithPagination/data"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/rs/zerolog/log"
)

// func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
// 	var requestPayload struct {
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 	}

// 	err := app.readJSON(w, r, &requestPayload)
// 	if err != nil {
// 		app.errorJSON(w, err, http.StatusBadRequest)
// 		return
// 	}

// 	// validate the user against the database
// 	user, err := app.Models.User.GetByEmail(requestPayload.Email)
// 	if err != nil {
// 		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
// 		return
// 	}

// 	valid, err := user.PasswordMatches(requestPayload.Password)
// 	if err != nil || !valid {
// 		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
// 		return
// 	}

// 	// log authentication
// 	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
// 	if err != nil {
// 		app.errorJSON(w, err)
// 		return
// 	}

// 	payload := jsonResponse{
// 		Error:   false,
// 		Message: fmt.Sprintf("Logged in user %s", user.Email),
// 		Data:    user,
// 	}

// 	app.writeJSON(w, http.StatusAccepted, payload)
// }

// func (app *Config) logRequest(name, data string) error {
// 	var entry struct {
// 		Name string `json:"name"`
// 		Data string `json:"data"`
// 	}

// 	entry.Name = name
// 	entry.Data = data

// 	jsonData, _ := json.Marshal(entry)
// 	logServiceURL := "http://logger-service/log"

// 	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		return err
// 	}

// 	client := &http.Client{}
// 	_, err = client.Do(request)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (app *Config) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var user data.User

	data, err := io.ReadAll(r.Body)
	if err != nil {
		app.errorJSON(w, err)
	}

	err = json.Unmarshal(data, &user)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
	}

	// log.Info().Msgf("Unmarshalling user json data: %v", user)

	_, err = app.Models.User.Insert(user)
	if err != nil {
		// log.Info().Msgf("Internal error: can't store to db:%v", err)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// fmt.Println(pgErr.Message)
			// fmt.Println(pgErr.Code)
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				http.Error(w, "Provided value already exists", http.StatusBadRequest) // the provided value already exists in the database
			default:
				http.Error(w, "Internal error: can't store to db", http.StatusInternalServerError)
			}

		}

	}
}

func (app *Config) GetEmployeeByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	// v := r.URL.Query().Get("v")
	fmt.Println("Url Parameters:", id)
	// fmt.Println("Query:", v)

	// userId, err := strconv.Atoi(id)
	// if err != nil {
	// 	http.Error(w, "Provided id is of unsupported format", http.StatusBadRequest)
	// }

	// check if the id does exist or not. If the id doesn't exist then it would crash the application
	// err = app.Models.User.CheckId(userId)
	err := app.Models.User.CheckId(id)
	if err == nil {
		// Now, get the user as the id does exist
		// user, err := app.Models.User.GetOne(userId)
		user, err := app.Models.User.GetOne(id)
		if err != nil {
			http.Error(w, "could not fetch record from db", http.StatusInternalServerError)
		}
		user.Password = "-"
		app.writeJSON(w, http.StatusAccepted, user)
	} else {
		// log.Info().Msgf("user id check error:- %v", err)
		http.Error(w, "Provided id doesn't exist", http.StatusBadRequest)
	}

}

func (app *Config) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	// v := r.URL.Query().Get("v")
	fmt.Println("Url Parameters:", id)
	// fmt.Println("Query:", v)

	// userId, err := strconv.Atoi(id)
	// if err != nil {
	// 	http.Error(w, "Provided id is of unsupported format", http.StatusBadRequest)
	// }

	// check if the id does exist or not. If the id doesn't exist then it would crash the application
	// err = app.Models.User.CheckId(userId)
	err := app.Models.User.CheckId(id)
	if err == nil {
		// Now, update the user as the id does exist
		var user data.User

		data, err := io.ReadAll(r.Body)
		if err != nil {
			app.errorJSON(w, err)
		}

		err = json.Unmarshal(data, &user)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
		}
		user.Password = "-"
		// user.ID = userId
		user.ID = id
		// User's password can't/shouldn't be changed through this method
		app.Models.User = user
		err = app.Models.User.Update()
		if err != nil {
			http.Error(w, "could not update the record in db", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "Provided user doesn't exist", http.StatusBadRequest)
	}
}

func (app *Config) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	// v := r.URL.Query().Get("v")
	fmt.Println("Url Parameters:", id)
	// fmt.Println("Query:", v)

	userId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Provided id is of unsupported format", http.StatusBadRequest)
	}

	// Now, check if the user is trying to delete other users or itself
	// If the user is a admin only then we'll let it delete other users
	// otherwise it would not be able to delete users other than itself
	//

	// check if the user is admin or not and also get the user id of the actual user who is making the request
	// requestMakingUserID, admin := app.GetIDOfRequestMakingUser(w, r)
	requestMakingUser, admin := app.GetIDOfRequestMakingUser(w, r)
	if admin {
		err := app.Models.User.DeleteByID(userId)
		if err != nil {
			http.Error(w, "could not delete the record from db", http.StatusInternalServerError)
		}
	} else {
		app.Models.User = *requestMakingUser
		err = app.Models.User.Delete()
		if err != nil {
			http.Error(w, "could not delete the record from db", http.StatusInternalServerError)
		}
	}
}

func (app *Config) GetIDOfRequestMakingUser(w http.ResponseWriter, r *http.Request) (*data.User, bool) {
	isAdmin := false
	var newUser *data.User
	username, _, ok := r.BasicAuth()
	if ok {
		// validate the user against the database
		user, err := app.Models.User.GetByEmail(username)
		if err != nil {
			app.errorJSON(w, errors.New("couldn't fetch record from db"), http.StatusInternalServerError)
			return &data.User{}, isAdmin
		}
		newUser = user
		if user.Email == "admin.example.com" {
			isAdmin = true
		}
	}

	return newUser, isAdmin
}

func (app *Config) GetAllEmployee(w http.ResponseWriter, r *http.Request) {
	var AllUsers []*data.User

	limitNumber := chi.URLParam(r, "limit")
	// cursorAsTimeStamp := r.URL.Query().Get("cursor")
	cursorAsTimeStamp := chi.URLParam(r, "cursor")
	fmt.Println("Url Parameters:", limitNumber)
	fmt.Println("length of Cursor:", len(cursorAsTimeStamp))
	fmt.Println("Cursor:", cursorAsTimeStamp)

	var decodeCursorStringTime time.Time
	var decodeCursorStringUUID string
	var err error
	var isFirstQuery bool

	if len(cursorAsTimeStamp) > 0 && cursorAsTimeStamp == "first" {
		// decodeCursorStringTime = time.Now()
		isFirstQuery = true
		decodeCursorStringTime = time.Date(1970, time.Month(1), 0, 0, 0, 0, 0, time.UTC) // epoch time: 1970, 1 January 00:00:00 UTC
		decodeCursorStringUUID = "aaaaaaa"
	} else {
		isFirstQuery = false
		decodeCursorStringTime, decodeCursorStringUUID, err = decodeCursor(cursorAsTimeStamp)
		if err != nil {
			app.errorJSON(w, errors.New("provided cursor is invalid"), http.StatusInternalServerError)
			log.Info().Msgf("Provided cursor is invalid:%v", err)
			return
		}
	}

	log.Info().Msgf("decodeCursorStringTime:%v", decodeCursorStringTime)
	log.Info().Msgf("decodeCursorStringUUID:%v", decodeCursorStringUUID)

	// decodeCursorStringTime:2024-06-02 05:00:36.357147 +0000 UTC" but we have to remove " +0000 UTC" part
	// from the timestamp otherwise sql query crashes and as a result returns empty slice of all users

	AllUsers, err = app.Models.User.GetAllForPagination(decodeCursorStringTime, decodeCursorStringUUID, isFirstQuery)
	if err != nil {
		log.Info().Msgf("couldn't fetch record from db: %v", err)
		app.errorJSON(w, errors.New("couldn't fetch record from db"), http.StatusInternalServerError)
		return
	}

	log.Info().Msgf("AllUsers:%v", AllUsers)

	// LastElementTimeForThisPage := AllUsers[len(AllUsers)-1].CreatedAt
	// LastElementUUIDForThisPage := AllUsers[len(AllUsers)-1].ID

	if len(AllUsers) > 0 {
		// LastElementTimeForThisPage := AllUsers[0].CreatedAt
		// LastElementUUIDForThisPage := AllUsers[0].ID
		LastElementTimeForThisPage := AllUsers[len(AllUsers)-1].CreatedAt
		LastElementUUIDForThisPage := AllUsers[len(AllUsers)-1].ID
		log.Info().Msgf("encodeCursor before encoding: LastElementTimeForThisPage %v", LastElementTimeForThisPage)
		log.Info().Msgf("encodeCursor before encoding: LastElementUUIDForThisPage %v", LastElementUUIDForThisPage)

		EncodedCursorString := encodeCursor(LastElementTimeForThisPage, LastElementUUIDForThisPage)

		type Result struct {
			Employees []*data.User
			TotalItem int
			Cursor    string
		}

		NewResult := Result{
			Employees: AllUsers,
			TotalItem: len(AllUsers),
			Cursor:    EncodedCursorString,
		}

		app.writeJSON(w, http.StatusAccepted, NewResult)
	} else {
		app.errorJSON(w, errors.New("there is no more record"), http.StatusBadRequest)
		return
	}

	// log.Info().Msgf("all employees: %v", allUsers)

}

func decodeCursor(encodedCursor string) (res time.Time, uuid string, err error) {
	// encodedCursor = strings.TrimSpace(encodedCursor)
	byt, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return
	}

	arrStr := strings.Split(string(byt), ",")
	if len(arrStr) != 2 {
		err = errors.New("cursor is invalid")
		return
	}

	res, err = time.Parse(time.RFC3339Nano, arrStr[0])
	// res, err = time.Parse(time.UTC.String(), arrStr[0]) ISO8601
	if err != nil {
		return
	}
	uuid = arrStr[1]

	return
}

func encodeCursor(t time.Time, uuid string) string {
	key := fmt.Sprintf("%s,%s", t.Format(time.RFC3339Nano), uuid)
	return base64.StdEncoding.EncodeToString([]byte(key))
}
