package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/vabshere/vernacular-auth/models"
)

type mockSignUpReq struct {
	name,
	email,
	password string
	code,
	httpCode int
}

type ResponseStruct struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func mockSaveUser(u *models.User) error {
	u.Id = 100000
	return nil
}

const defaultPass = "pass"

func TestSignUp(t *testing.T) {
	oldSaveUser := saveUser
	defer func() {
		saveUser = oldSaveUser
	}()

	saveUser = mockSaveUser
	mockSignUpRequests := []mockSignUpReq{
		// {"", "abc@adb.ab", defaultPass, 1, http.StatusBadRequest},                  // empty name
		// {"foo", "", defaultPass, 1, http.StatusBadRequest},                         // empty email
		// {"foo", "abc@adb.c", "", 1, http.StatusBadRequest},                         // empty password
		// {"foo", "abcadb.c", defaultPass, 1, http.StatusBadRequest},                 // invalid email
		// {"foo", "abcadb@.c", defaultPass, 1, http.StatusBadRequest},                // invalid email
		// {"foo", "abcadb@c", defaultPass, 1, http.StatusBadRequest},                 // invalid email
		// {"foo", "abc ab@ab.c", defaultPass, 1, http.StatusBadRequest},              // invalid email
		{"foo", "abc.hj@dgd.dd", defaultPass, 0, http.StatusOK}, // valid input
	}

	for _, mockRequest := range mockSignUpRequests {
		reqBody := url.Values{}
		reqBody.Set("name", mockRequest.name)
		reqBody.Add("password", mockRequest.password)
		reqBody.Add("email", mockRequest.email)
		req, err := http.NewRequest(http.MethodPost, "/reg", strings.NewReader(reqBody.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(reqBody.Encode())))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		var user *models.User
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user = SignUp(w, r)
		})
		handler.ServeHTTP(rr, req)
		correctStatusFlag := true
		if status := rr.Code; status != mockRequest.httpCode {
			correctStatusFlag = false
			t.Errorf("handler returned wrong status code: got %d want %d",
				status, mockRequest.httpCode)
		}

		if status := rr.Code; status == http.StatusCreated && correctStatusFlag {
			if user == nil {
				var response ResponseStruct
				decoder := json.NewDecoder(rr.Body)
				err := decoder.Decode(&response)
				if err != nil {
					t.Error(err)
				}

				u, ok := response.Data.(models.User)
				if ok && (response.Code != 0 || u.Name != mockRequest.name || u.Email != mockRequest.email) {
					t.Errorf("Wrong JSON value returned")
				}
			} else if user.Email != mockRequest.email {
				t.Errorf("Wrong JSON value returned")
			}

		}
	}
}

type mockSignInReq struct {
	email,
	password string
	code,
	httpCode int
}

var mockSignInRequests = []mockSignInReq{
	// {"", defaultPass, 1, http.StatusBadRequest},             // empty email
	// {"abc@adb.c", "", 1, http.StatusBadRequest},             // empty password
	// {"abcadb.c", defaultPass, 1, http.StatusBadRequest},     // invalid email
	// {"abcadb@.c", defaultPass, 1, http.StatusBadRequest},    // invalid email
	// {"abcadb@c", defaultPass, 1, http.StatusBadRequest},     // invalid email
	// {"abc ab@ab.c", defaultPass, 1, http.StatusBadRequest},  // invalid email
	// {"abc@adb.abc", defaultPass + "cash", 1, http.StatusOK}, // password mismatch
	{"abc@adb.abc", defaultPass, 0, http.StatusOK}, // valid input
}

func mockGetUserByEmail(e string) (*models.User, error) {
	var u *models.User
	for _, mockRequest := range mockSignInRequests {
		if mockRequest.email == e {
			return &models.User{Name: "abc", Email: mockRequest.email, Password: []byte("$2a$10$4yRvhHitq43PspRFg1wDEewAcurn1tA3H/Njo067YP9yRKjVU5sae"), Id: 100000}, nil // $2a$10$4yRvhHitq43PspRFg1wDEewAcurn1tA3H/Njo067YP9yRKjVU5sae is bcrypt hash for "pass"
		}
	}
	return u, nil
}

func TestSignIn(t *testing.T) {
	oldGetUserByEmail := getUserByEmail
	defer func() {
		getUserByEmail = oldGetUserByEmail
	}()
	getUserByEmail = mockGetUserByEmail

	for _, mockRequest := range mockSignInRequests {
		reqBody := url.Values{}
		reqBody.Set("email", mockRequest.email)
		reqBody.Add("password", mockRequest.password)
		req, err := http.NewRequest(http.MethodPost, "/oauth", strings.NewReader(reqBody.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(reqBody.Encode())))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		var user *models.User
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user = SignIn(w, r)
		})
		handler.ServeHTTP(rr, req)

		correctStatusFlag := true
		if status := rr.Code; status != mockRequest.httpCode {
			correctStatusFlag = false
			t.Errorf("handler returned wrong status code: got %d want %d",
				status, mockRequest.httpCode)
		}

		if status := rr.Code; status == http.StatusOK && correctStatusFlag {
			if user == nil {
				var response ResponseStruct
				decoder := json.NewDecoder(rr.Body)
				err := decoder.Decode(&response)
				if err != nil {
					t.Error(err)
				}

				u, ok := response.Data.(models.User)
				if ok && (response.Code != 0 || u.Email != mockRequest.email) {
					t.Errorf("Wrong JSON value returned")
				} else if !ok && response.Code == 0 {
					t.Errorf("Wrong JSON value returned")
				}
			} else if user.Email != mockRequest.email {
				t.Errorf("Wrong JSON value returned")
			}
		}
	}
}
