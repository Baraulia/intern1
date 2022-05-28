package handlers

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
	"tranee_service/internal/logging"
	"tranee_service/models"
	"tranee_service/services"
	mockservice "tranee_service/services/mocks"
)

func TestHandler_getHobbies(t *testing.T) {
	type mockBehavior func(s *mockservice.MockAppHobbies)

	testTable := []struct {
		name                string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			mockBehavior: func(s *mockservice.MockAppHobbies) {
				s.EXPECT().GetHobbies().Return([]models.ResponseHobby{
					{
						Id:   1,
						Name: "test name",
					},
					{
						Id:   2,
						Name: "test name2",
					},
				}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `[{"id":1,"name":"test name"},{"id":2,"name":"test name2"}]`,
		},
		{
			name: "Server error",
			mockBehavior: func(s *mockservice.MockAppHobbies) {
				s.EXPECT().GetHobbies().Return(nil, errors.New("server error"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: "server error\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			appService := mockservice.NewMockAppHobbies(c)
			testCase.mockBehavior(appService)
			logger := logging.GetLoggerLogrus()
			serv := &services.Service{AppHobbies: appService}
			handler := NewHandler(serv, logger)

			r := handler.InitRoutes()

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", "/hobbies", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_createHobby(t *testing.T) {
	type mockBehavior func(s *mockservice.MockAppHobbies, hobby *models.Hobby)

	testTable := []struct {
		name               string
		inputBody          string
		inputHobby         *models.Hobby
		mockBehavior       mockBehavior
		expectedStatusCode int
	}{
		{
			name:       "OK",
			inputBody:  `{"name":"testName"}`,
			inputHobby: &models.Hobby{Name: "testName"},
			mockBehavior: func(s *mockservice.MockAppHobbies, hobby *models.Hobby) {
				s.EXPECT().CreateHobby(hobby).Return(1, nil)
			},
			expectedStatusCode: 201,
		},
		{
			name:               "Invalid hobby`s name",
			inputBody:          `{"name":"1"}`,
			inputHobby:         &models.Hobby{},
			mockBehavior:       func(s *mockservice.MockAppHobbies, hobby *models.Hobby) {},
			expectedStatusCode: 400,
		},
		{
			name:       "Server error",
			inputBody:  `{"name":"testName"}`,
			inputHobby: &models.Hobby{Name: "testName"},
			mockBehavior: func(s *mockservice.MockAppHobbies, hobby *models.Hobby) {
				s.EXPECT().CreateHobby(hobby).Return(0, errors.New("server error"))
			},
			expectedStatusCode: 500,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			appService := mockservice.NewMockAppHobbies(c)
			testCase.mockBehavior(appService, testCase.inputHobby)
			logger := logging.GetLoggerLogrus()
			serv := &services.Service{AppHobbies: appService}
			handler := NewHandler(serv, logger)

			r := handler.InitRoutes()

			w := httptest.NewRecorder()

			req := httptest.NewRequest("POST", "/hobbies", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
		})
	}
}
