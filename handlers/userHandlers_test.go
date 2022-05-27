package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
	"tranee_service/internal/logging"
	"tranee_service/models"
	"tranee_service/services"
	mockservice "tranee_service/services/mocks"
)

func TestHandler_createUser(t *testing.T) {
	type mockBehavior func(s *mockservice.MockAppUsers, country *models.User)

	testTable := []struct {
		name               string
		inputBody          string
		inputUser          *models.User
		mockBehavior       mockBehavior
		expectedStatusCode int
	}{
		{
			name:      "OK",
			inputBody: `{"name":"testName","email":"test@test.ru","description":"test desc","country_id":1,"hobbies":[1,2,3]}`,
			inputUser: &models.User{
				Name:        "testName",
				Email:       "test@test.ru",
				Description: "test desc",
				CountryId:   1,
				Hobbies:     []int{1, 2, 3},
			},
			mockBehavior: func(s *mockservice.MockAppUsers, country *models.User) {
				s.EXPECT().CreateUser(country).Return(1, nil)
			},
			expectedStatusCode: 201,
		},
		{
			name:      "Incorrect data came from the request",
			inputBody: `{"name":"testName","email":"test.test.ru","description":"test desc","country_id":1,"hobbies":[1,2,3]}`,
			inputUser: &models.User{
				Name:        "testName",
				Email:       "test@test.ru",
				Description: "test desc",
				CountryId:   1,
				Hobbies:     []int{1, 2, 3},
			},
			mockBehavior:       func(s *mockservice.MockAppUsers, country *models.User) {},
			expectedStatusCode: 400,
		},
		{
			name:      "Server error",
			inputBody: `{"name":"testName","email":"test@test.ru","description":"test desc","country_id":1,"hobbies":[1,2,3]}`,
			inputUser: &models.User{
				Name:        "testName",
				Email:       "test@test.ru",
				Description: "test desc",
				CountryId:   1,
				Hobbies:     []int{1, 2, 3},
			},
			mockBehavior: func(s *mockservice.MockAppUsers, country *models.User) {
				s.EXPECT().CreateUser(country).Return(0, errors.New("server error"))
			},
			expectedStatusCode: 500,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init dependencies
			c := gomock.NewController(t)
			defer c.Finish()
			appService := mockservice.NewMockAppUsers(c)
			testCase.mockBehavior(appService, testCase.inputUser)
			logger := logging.GetLoggerLogrus()
			serv := &services.Service{AppUsers: appService}
			handler := NewHandler(serv, logger)

			//Init server
			r := handler.InitRoutes()

			//Test request
			w := httptest.NewRecorder()

			req := httptest.NewRequest("POST", "/users", bytes.NewBufferString(testCase.inputBody))

			//Execute the request
			r.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
		})
	}
}

func TestHandler_changeUser(t *testing.T) {
	type mockBehavior func(s *mockservice.MockAppUsers, user *models.User, userId int)

	testTable := []struct {
		name               string
		inputBody          string
		pathId             string
		inputId            int
		inputUser          *models.User
		mockBehavior       mockBehavior
		expectedStatusCode int
	}{
		{
			name:      "OK",
			pathId:    "1",
			inputId:   1,
			inputBody: `{"name":"testName","email":"test@test.ru","description":"test desc","country_id":1,"hobbies":[1,2,3]}`,
			inputUser: &models.User{
				Name:        "testName",
				Email:       "test@test.ru",
				Description: "test desc",
				CountryId:   1,
				Hobbies:     []int{1, 2, 3},
			},
			mockBehavior: func(s *mockservice.MockAppUsers, user *models.User, userId int) {
				s.EXPECT().ChangeUser(user, userId).Return(nil)
			},
			expectedStatusCode: 204,
		},
		{
			name:               "Incorrect data came from the request",
			pathId:             "1",
			inputId:            1,
			inputBody:          `{"name":"testName","email":"test.test.ru","description":"test desc","country_id":1,"hobbies":[1,2,3]}`,
			inputUser:          &models.User{},
			mockBehavior:       func(s *mockservice.MockAppUsers, user *models.User, userId int) {},
			expectedStatusCode: 400,
		},
		{
			name:      "Server error",
			pathId:    "1",
			inputId:   1,
			inputBody: `{"name":"testName","email":"test@test.ru","description":"test desc","country_id":1,"hobbies":[1,2,3]}`,
			inputUser: &models.User{
				Name:        "testName",
				Email:       "test@test.ru",
				Description: "test desc",
				CountryId:   1,
				Hobbies:     []int{1, 2, 3},
			},
			mockBehavior: func(s *mockservice.MockAppUsers, user *models.User, userId int) {
				s.EXPECT().ChangeUser(user, userId).Return(errors.New("server error"))
			},
			expectedStatusCode: 500,
		},
		{
			name:      "User does not exist",
			pathId:    "1",
			inputId:   1,
			inputBody: `{"name":"testName","email":"test@test.ru","description":"test desc","country_id":1,"hobbies":[1,2,3]}`,
			inputUser: &models.User{
				Name:        "testName",
				Email:       "test@test.ru",
				Description: "test desc",
				CountryId:   1,
				Hobbies:     []int{1, 2, 3},
			},
			mockBehavior: func(s *mockservice.MockAppUsers, user *models.User, userId int) {
				s.EXPECT().ChangeUser(user, userId).Return(errors.New("such a user does not exist"))
			},
			expectedStatusCode: 404,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init dependencies
			c := gomock.NewController(t)
			defer c.Finish()
			appService := mockservice.NewMockAppUsers(c)
			testCase.mockBehavior(appService, testCase.inputUser, testCase.inputId)
			logger := logging.GetLoggerLogrus()
			serv := &services.Service{AppUsers: appService}
			handler := NewHandler(serv, logger)

			//Init server
			r := handler.InitRoutes()

			//Test request
			w := httptest.NewRecorder()

			req := httptest.NewRequest("PUT", fmt.Sprintf("/users/%s", testCase.pathId), bytes.NewBufferString(testCase.inputBody))

			//Execute the request
			r.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
		})
	}
}

func TestHandler_getUsers(t *testing.T) {
	type mockBehavior func(s *mockservice.MockAppUsers, filter *models.Options)

	testTable := []struct {
		name                string
		pathQuery           string
		inputFilter         *models.Options
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:        "OK",
			pathQuery:   "?page=1&limit=2",
			inputFilter: &models.Options{Page: 1, Limit: 2},
			mockBehavior: func(s *mockservice.MockAppUsers, filter *models.Options) {
				s.EXPECT().GetUsers(filter).Return([]models.ResponseUser{
					{
						Name:        "test name",
						Email:       "test@email.ru",
						Description: "test",
						CountryId:   1,
						Hobbies:     []int{1, 2, 3},
					},
					{
						Name:        "test name2",
						Email:       "test2@email.ru",
						Description: "test",
						CountryId:   1,
						Hobbies:     []int{1, 2, 3},
					},
				}, 1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `[{"id":0,"name":"test name","email":"test@email.ru","description":"test","country_id":1,"hobbies":[1,2,3]},{"id":0,"name":"test name2","email":"test2@email.ru","description":"test","country_id":1,"hobbies":[1,2,3]}]`,
		},
		{
			name:        "OK without pagination",
			pathQuery:   "",
			inputFilter: &models.Options{},
			mockBehavior: func(s *mockservice.MockAppUsers, filter *models.Options) {
				s.EXPECT().GetUsers(filter).Return([]models.ResponseUser{
					{
						Name:        "test name",
						Email:       "test@email.ru",
						Description: "test",
						CountryId:   1,
						Hobbies:     []int{1, 2, 3},
					},
					{
						Name:        "test name2",
						Email:       "test2@email.ru",
						Description: "test",
						CountryId:   1,
						Hobbies:     []int{1, 2, 3},
					},
				}, 1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `[{"id":0,"name":"test name","email":"test@email.ru","description":"test","country_id":1,"hobbies":[1,2,3]},{"id":0,"name":"test name2","email":"test2@email.ru","description":"test","country_id":1,"hobbies":[1,2,3]}]`,
		},
		{
			name:                "Invalid query",
			pathQuery:           "?page=-1&limit=2",
			inputFilter:         &models.Options{},
			mockBehavior:        func(s *mockservice.MockAppUsers, filter *models.Options) {},
			expectedStatusCode:  400,
			expectedRequestBody: "Invalid url request\n",
		},
		{
			name:                "Invalid query2",
			pathQuery:           "?page=a&limit=-2",
			inputFilter:         &models.Options{},
			mockBehavior:        func(s *mockservice.MockAppUsers, filter *models.Options) {},
			expectedStatusCode:  400,
			expectedRequestBody: "Invalid url request\n",
		},
		{
			name:        "Server error",
			pathQuery:   "?page=1&limit=2",
			inputFilter: &models.Options{Page: 1, Limit: 2},
			mockBehavior: func(s *mockservice.MockAppUsers, filter *models.Options) {
				s.EXPECT().GetUsers(filter).Return(nil, 0, errors.New("server error"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: "server error\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init dependencies
			c := gomock.NewController(t)
			defer c.Finish()
			appService := mockservice.NewMockAppUsers(c)
			testCase.mockBehavior(appService, testCase.inputFilter)
			logger := logging.GetLoggerLogrus()
			serv := &services.Service{AppUsers: appService}
			handler := NewHandler(serv, logger)

			//Init server
			r := handler.InitRoutes()

			//Test request
			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", fmt.Sprintf("/users%s", testCase.pathQuery), nil)

			//Execute the request
			r.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_getUserById(t *testing.T) {
	type mockBehavior func(s *mockservice.MockAppUsers, userId int)

	testTable := []struct {
		name                string
		pathQuery           string
		userId              int
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			pathQuery: "1",
			userId:    1,
			mockBehavior: func(s *mockservice.MockAppUsers, userId int) {
				s.EXPECT().GetUserById(userId).Return(&models.ResponseUser{
					Name:        "test name",
					Email:       "test@email.ru",
					Description: "test",
					CountryId:   1,
					Hobbies:     []int{1, 2, 3},
				}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":0,"name":"test name","email":"test@email.ru","description":"test","country_id":1,"hobbies":[1,2,3]}`,
		},
		{
			name:                "Invalid query",
			pathQuery:           "-1",
			userId:              1,
			mockBehavior:        func(s *mockservice.MockAppUsers, userId int) {},
			expectedStatusCode:  400,
			expectedRequestBody: "Invalid url request\n",
		},
		{
			name:                "Invalid query2",
			pathQuery:           "a",
			userId:              1,
			mockBehavior:        func(s *mockservice.MockAppUsers, userId int) {},
			expectedStatusCode:  400,
			expectedRequestBody: "Invalid url request\n",
		},
		{
			name:      "Server error",
			pathQuery: "1",
			userId:    1,
			mockBehavior: func(s *mockservice.MockAppUsers, userId int) {
				s.EXPECT().GetUserById(userId).Return(nil, errors.New("server error"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: "server error\n",
		},
		{
			name:      "Such user does not exist",
			pathQuery: "1",
			userId:    1,
			mockBehavior: func(s *mockservice.MockAppUsers, userId int) {
				s.EXPECT().GetUserById(userId).Return(nil, errors.New("such a user does not exist"))
			},
			expectedStatusCode:  404,
			expectedRequestBody: "such user does not exist\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init dependencies
			c := gomock.NewController(t)
			defer c.Finish()
			appService := mockservice.NewMockAppUsers(c)
			testCase.mockBehavior(appService, testCase.userId)
			logger := logging.GetLoggerLogrus()
			serv := &services.Service{AppUsers: appService}
			handler := NewHandler(serv, logger)

			//Init server
			r := handler.InitRoutes()

			//Test request
			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", fmt.Sprintf("/users/%s", testCase.pathQuery), nil)

			//Execute the request
			r.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_deleteUser(t *testing.T) {
	type mockBehavior func(s *mockservice.MockAppUsers, userId int)

	testTable := []struct {
		name               string
		pathQuery          string
		userId             int
		mockBehavior       mockBehavior
		expectedStatusCode int
	}{
		{
			name:      "OK",
			pathQuery: "1",
			userId:    1,
			mockBehavior: func(s *mockservice.MockAppUsers, userId int) {
				s.EXPECT().DeleteUser(userId).Return(nil)
			},
			expectedStatusCode: 204,
		},
		{
			name:               "Invalid query",
			pathQuery:          "-1",
			userId:             1,
			mockBehavior:       func(s *mockservice.MockAppUsers, userId int) {},
			expectedStatusCode: 400,
		},
		{
			name:               "Invalid query2",
			pathQuery:          "a",
			userId:             1,
			mockBehavior:       func(s *mockservice.MockAppUsers, userId int) {},
			expectedStatusCode: 400,
		},
		{
			name:      "Server error",
			pathQuery: "1",
			userId:    1,
			mockBehavior: func(s *mockservice.MockAppUsers, userId int) {
				s.EXPECT().DeleteUser(userId).Return(errors.New("server error"))
			},
			expectedStatusCode: 500,
		},
		{
			name:      "Such user does not exist",
			pathQuery: "1",
			userId:    1,
			mockBehavior: func(s *mockservice.MockAppUsers, userId int) {
				s.EXPECT().DeleteUser(userId).Return(errors.New("user with such Id does not exist"))
			},
			expectedStatusCode: 404,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init dependencies
			c := gomock.NewController(t)
			defer c.Finish()
			appService := mockservice.NewMockAppUsers(c)
			testCase.mockBehavior(appService, testCase.userId)
			logger := logging.GetLoggerLogrus()
			serv := &services.Service{AppUsers: appService}
			handler := NewHandler(serv, logger)

			//Init server
			r := handler.InitRoutes()

			//Test request
			w := httptest.NewRecorder()

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/users/%s", testCase.pathQuery), nil)

			//Execute the request
			r.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
		})
	}
}

func TestHandler_getHobbyByUserId(t *testing.T) {
	type mockBehavior func(s *mockservice.MockAppUsers, id int)

	testTable := []struct {
		name                string
		pathQuery           string
		inputId             int
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			pathQuery: "1",
			inputId:   1,
			mockBehavior: func(s *mockservice.MockAppUsers, id int) {
				s.EXPECT().GetHobbyByUserId(id).Return([]int{1, 2}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `[1,2]`,
		},
		{
			name:                "Invalid query",
			pathQuery:           "a",
			inputId:             0,
			mockBehavior:        func(s *mockservice.MockAppUsers, id int) {},
			expectedStatusCode:  400,
			expectedRequestBody: "Invalid url request\n",
		},
		{
			name:                "Invalid query2",
			pathQuery:           "-1",
			inputId:             0,
			mockBehavior:        func(s *mockservice.MockAppUsers, id int) {},
			expectedStatusCode:  400,
			expectedRequestBody: "Invalid url request\n",
		},
		{
			name:      "Such user does not exist",
			pathQuery: "1",
			inputId:   1,
			mockBehavior: func(s *mockservice.MockAppUsers, id int) {
				s.EXPECT().GetHobbyByUserId(id).Return(nil, errors.New("user with such Id does not exist"))
			},
			expectedStatusCode:  404,
			expectedRequestBody: "such user does not exist\n",
		},
		{
			name:      "Server error",
			pathQuery: "1",
			inputId:   1,
			mockBehavior: func(s *mockservice.MockAppUsers, id int) {
				s.EXPECT().GetHobbyByUserId(id).Return(nil, errors.New("server error"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: "server error\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init dependencies
			c := gomock.NewController(t)
			defer c.Finish()
			appService := mockservice.NewMockAppUsers(c)
			testCase.mockBehavior(appService, testCase.inputId)
			logger := logging.GetLoggerLogrus()
			serv := &services.Service{AppUsers: appService}
			handler := NewHandler(serv, logger)

			//Init server
			r := handler.InitRoutes()

			//Test request
			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", fmt.Sprintf("/users/%s/hobbies", testCase.pathQuery), nil)

			//Execute the request
			r.ServeHTTP(w, req)

			//Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
