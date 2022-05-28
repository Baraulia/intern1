package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
	"tranee_service/MyErrors"
	"tranee_service/internal/logging"
	"tranee_service/models"
	"tranee_service/services"
	mockservice "tranee_service/services/mocks"
)

func TestGetAllCountries(t *testing.T) {
	type mockBehavior func(s *mockservice.MockAppCountries, filter *models.Filters)

	testTable := []struct {
		name                string
		pathQuery           string
		inputFilter         *models.Filters
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:        "OK",
			pathQuery:   "?page=1&limit=2",
			inputFilter: &models.Filters{Page: 1, Limit: 2},
			mockBehavior: func(s *mockservice.MockAppCountries, filter *models.Filters) {
				s.EXPECT().GetCountries(filter).Return([]models.Country{
					{
						Name:            "test name",
						FullName:        "test full name",
						EnglishName:     "test english name",
						Alpha2:          "tt",
						Alpha3:          "ttt",
						Iso:             1000,
						Location:        "test location",
						LocationPrecise: "test location precise",
						Url:             "test url",
					},
					{
						Name:            "test name2",
						FullName:        "test full name2",
						EnglishName:     "test english name2",
						Alpha2:          "tp",
						Alpha3:          "tpt",
						Iso:             1001,
						Location:        "test location",
						LocationPrecise: "test location precise",
						Url:             "test url2",
					},
				}, 1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `[{"name":"test name","full_name":"test full name","english_name":"test english name","alpha_2":"tt","alpha_3":"ttt","iso":1000,"location":"test location","location_precise":"test location precise","url":"test url"},{"name":"test name2","full_name":"test full name2","english_name":"test english name2","alpha_2":"tp","alpha_3":"tpt","iso":1001,"location":"test location","location_precise":"test location precise","url":"test url2"}]`,
		},
		{
			name:        "OK without pagination",
			pathQuery:   "",
			inputFilter: &models.Filters{},
			mockBehavior: func(s *mockservice.MockAppCountries, filter *models.Filters) {
				s.EXPECT().GetCountries(filter).Return([]models.Country{
					{
						Name:            "test name",
						FullName:        "test full name",
						EnglishName:     "test english name",
						Alpha2:          "tt",
						Alpha3:          "ttt",
						Iso:             1000,
						Location:        "test location",
						LocationPrecise: "test location precise",
						Url:             "test url",
					},
					{
						Name:            "test name2",
						FullName:        "test full name2",
						EnglishName:     "test english name2",
						Alpha2:          "tp",
						Alpha3:          "tpt",
						Iso:             1001,
						Location:        "test location",
						LocationPrecise: "test location precise",
						Url:             "test url2",
					},
				}, 1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `[{"name":"test name","full_name":"test full name","english_name":"test english name","alpha_2":"tt","alpha_3":"ttt","iso":1000,"location":"test location","location_precise":"test location precise","url":"test url"},{"name":"test name2","full_name":"test full name2","english_name":"test english name2","alpha_2":"tp","alpha_3":"tpt","iso":1001,"location":"test location","location_precise":"test location precise","url":"test url2"}]`,
		},
		{
			name:                "Invalid query",
			pathQuery:           "?page=-1&limit=2",
			inputFilter:         &models.Filters{},
			mockBehavior:        func(s *mockservice.MockAppCountries, filter *models.Filters) {},
			expectedStatusCode:  400,
			expectedRequestBody: "invalid url request\n",
		},
		{
			name:                "Invalid query2",
			pathQuery:           "?page=a&limit=-2",
			inputFilter:         &models.Filters{},
			mockBehavior:        func(s *mockservice.MockAppCountries, filter *models.Filters) {},
			expectedStatusCode:  400,
			expectedRequestBody: "invalid url request\n",
		},
		{
			name:        "Server error",
			pathQuery:   "?page=1&limit=2",
			inputFilter: &models.Filters{Page: 1, Limit: 2},
			mockBehavior: func(s *mockservice.MockAppCountries, filter *models.Filters) {
				s.EXPECT().GetCountries(filter).Return(nil, 0, errors.New("server error"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: "server error\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			appService := mockservice.NewMockAppCountries(c)
			testCase.mockBehavior(appService, testCase.inputFilter)
			logger := logging.GetLoggerLogrus()
			serv := &services.Service{AppCountries: appService}
			handler := NewHandler(serv, logger)

			r := handler.InitRoutes()

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", fmt.Sprintf("/countries%s", testCase.pathQuery), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestGetOneCountry(t *testing.T) {
	type mockBehavior func(s *mockservice.MockAppCountries, inputId string)

	testTable := []struct {
		name                string
		pathId              string
		inputId             string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:    "OK",
			pathId:  "tt",
			inputId: "TT",
			mockBehavior: func(s *mockservice.MockAppCountries, inputId string) {
				s.EXPECT().GetOneCountry(inputId).Return(&models.Country{
					Name:            "test name",
					FullName:        "test full name",
					EnglishName:     "test english name",
					Alpha2:          "tt",
					Alpha3:          "ttt",
					Iso:             1000,
					Location:        "test location",
					LocationPrecise: "test location precise",
					Url:             "test url",
				}, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"name":"test name","full_name":"test full name","english_name":"test english name","alpha_2":"tt","alpha_3":"ttt","iso":1000,"location":"test location","location_precise":"test location precise","url":"test url"}`,
		},
		{
			name:                "Invalid id",
			pathId:              "1",
			inputId:             "",
			mockBehavior:        func(s *mockservice.MockAppCountries, inputId string) {},
			expectedStatusCode:  400,
			expectedRequestBody: "invalid url parameter\n",
		},
		{
			name:    "Such a news does not exist",
			pathId:  "tt",
			inputId: "TT",
			mockBehavior: func(s *mockservice.MockAppCountries, inputId string) {
				s.EXPECT().GetOneCountry(inputId).Return(nil, MyErrors.DoesNotExist)
			},
			expectedStatusCode:  404,
			expectedRequestBody: "object with this id does not exist\n",
		},
		{
			name:    "Server error",
			pathId:  "tt",
			inputId: "TT",
			mockBehavior: func(s *mockservice.MockAppCountries, inputId string) {
				s.EXPECT().GetOneCountry(inputId).Return(nil, errors.New("server error"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: "server error\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			appService := mockservice.NewMockAppCountries(c)
			testCase.mockBehavior(appService, testCase.inputId)
			logger := logging.GetLoggerLogrus()
			serv := &services.Service{AppCountries: appService}
			handler := NewHandler(serv, logger)

			r := handler.InitRoutes()

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", fmt.Sprintf("/countries/%s", testCase.pathId), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestCreateCountry(t *testing.T) {
	type mockBehavior func(s *mockservice.MockAppCountries, country *models.ResponseCountry)

	testTable := []struct {
		name               string
		inputBody          string
		inputCountry       *models.ResponseCountry
		mockBehavior       mockBehavior
		expectedStatusCode int
	}{
		{
			name:      "OK",
			inputBody: `{"name":"test name","full_name":"test full name","english_name":"test english name","alpha_2":"tt","alpha_3":"ttt","iso":1000,"location":"test location","location_precise":"test location precise"}`,
			inputCountry: &models.ResponseCountry{
				Name:            "test name",
				FullName:        "test full name",
				EnglishName:     "test english name",
				Alpha2:          "tt",
				Alpha3:          "ttt",
				Iso:             1000,
				Location:        "test location",
				LocationPrecise: "test location precise",
			},
			mockBehavior: func(s *mockservice.MockAppCountries, country *models.ResponseCountry) {
				s.EXPECT().CreateCountry(country).Return("tt", nil)
			},
			expectedStatusCode: 201,
		},
		{
			name:               "Incorrect data came from the request",
			inputBody:          `{"name":"test name","full_name":"test full name","english_name":"test english name","alpha_2":"ttqq","alpha_3":"ttt","iso":1000,"location":"test location","location_precise":"test location precise"}`,
			inputCountry:       &models.ResponseCountry{},
			mockBehavior:       func(s *mockservice.MockAppCountries, country *models.ResponseCountry) {},
			expectedStatusCode: 400,
		},
		{
			name:      "Server error",
			inputBody: `{"name":"test name","full_name":"test full name","english_name":"test english name","alpha_2":"tt","alpha_3":"ttt","iso":1000,"location":"test location","location_precise":"test location precise"}`,
			inputCountry: &models.ResponseCountry{
				Name:            "test name",
				FullName:        "test full name",
				EnglishName:     "test english name",
				Alpha2:          "tt",
				Alpha3:          "ttt",
				Iso:             1000,
				Location:        "test location",
				LocationPrecise: "test location precise",
			},
			mockBehavior: func(s *mockservice.MockAppCountries, country *models.ResponseCountry) {
				s.EXPECT().CreateCountry(country).Return("", errors.New("server error"))
			},
			expectedStatusCode: 500,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			appService := mockservice.NewMockAppCountries(c)
			testCase.mockBehavior(appService, testCase.inputCountry)
			logger := logging.GetLoggerLogrus()
			serv := &services.Service{AppCountries: appService}
			handler := NewHandler(serv, logger)

			r := handler.InitRoutes()

			w := httptest.NewRecorder()

			req := httptest.NewRequest("POST", "/countries", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
		})
	}
}

func TestChangeCountry(t *testing.T) {
	type mockBehavior func(s *mockservice.MockAppCountries, country *models.ResponseCountry, countryId string)

	testTable := []struct {
		name               string
		pathId             string
		inputBody          string
		inputCountry       *models.ResponseCountry
		inputId            string
		mockBehavior       mockBehavior
		expectedStatusCode int
	}{
		{
			name:      "OK",
			pathId:    "tt",
			inputBody: `{"name":"test name","full_name":"test full name","english_name":"test english name","alpha_2":"tt","alpha_3":"ttt","iso":1000,"location":"test location","location_precise":"test location precise"}`,
			inputCountry: &models.ResponseCountry{
				Name:            "test name",
				FullName:        "test full name",
				EnglishName:     "test english name",
				Alpha2:          "tt",
				Alpha3:          "ttt",
				Iso:             1000,
				Location:        "test location",
				LocationPrecise: "test location precise",
			},
			inputId: "TT",
			mockBehavior: func(s *mockservice.MockAppCountries, country *models.ResponseCountry, countryId string) {
				s.EXPECT().ChangeCountry(country, countryId).Return(nil)
			},
			expectedStatusCode: 204,
		},
		{
			name:               "Incorrect data came from the request",
			pathId:             "tt",
			inputId:            "TT",
			inputBody:          `{"name":"test name","full_name":"test full name","english_name":"test english name","alpha_2":"ttqq","alpha_3":"ttt","iso":1000,"location":"test location","location_precise":"test location precise"}`,
			inputCountry:       &models.ResponseCountry{},
			mockBehavior:       func(s *mockservice.MockAppCountries, country *models.ResponseCountry, countryId string) {},
			expectedStatusCode: 400,
		},
		{
			name:               "Incorrect country id",
			pathId:             "2",
			inputId:            "",
			inputBody:          `{"name":"test name","full_name":"test full name","english_name":"test english name","alpha_2":"ttqq","alpha_3":"ttt","iso":1000,"location":"test location","location_precise":"test location precise"}`,
			inputCountry:       &models.ResponseCountry{},
			mockBehavior:       func(s *mockservice.MockAppCountries, country *models.ResponseCountry, countryId string) {},
			expectedStatusCode: 400,
		},
		{
			name:      "Server error",
			pathId:    "tt",
			inputId:   "TT",
			inputBody: `{"name":"test name","full_name":"test full name","english_name":"test english name","alpha_2":"tt","alpha_3":"ttt","iso":1000,"location":"test location","location_precise":"test location precise"}`,
			inputCountry: &models.ResponseCountry{
				Name:            "test name",
				FullName:        "test full name",
				EnglishName:     "test english name",
				Alpha2:          "tt",
				Alpha3:          "ttt",
				Iso:             1000,
				Location:        "test location",
				LocationPrecise: "test location precise",
			},
			mockBehavior: func(s *mockservice.MockAppCountries, country *models.ResponseCountry, countryId string) {
				s.EXPECT().ChangeCountry(country, countryId).Return(errors.New("server error"))
			},
			expectedStatusCode: 500,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			appService := mockservice.NewMockAppCountries(c)
			testCase.mockBehavior(appService, testCase.inputCountry, testCase.inputId)
			logger := logging.GetLoggerLogrus()
			serv := &services.Service{AppCountries: appService}
			handler := NewHandler(serv, logger)

			r := handler.InitRoutes()

			w := httptest.NewRecorder()

			req := httptest.NewRequest("PUT", fmt.Sprintf("/countries/%s", testCase.pathId), bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
		})
	}
}

func TestDeleteCountry(t *testing.T) {
	type mockBehavior func(s *mockservice.MockAppCountries, inputId string)

	testTable := []struct {
		name               string
		pathId             string
		inputId            string
		mockBehavior       mockBehavior
		expectedStatusCode int
	}{
		{
			name:    "OK",
			pathId:  "tt",
			inputId: "TT",
			mockBehavior: func(s *mockservice.MockAppCountries, inputId string) {
				s.EXPECT().DeleteCountry(inputId).Return(nil)
			},
			expectedStatusCode: 204,
		},
		{
			name:               "Invalid id",
			pathId:             "1",
			inputId:            "",
			mockBehavior:       func(s *mockservice.MockAppCountries, inputId string) {},
			expectedStatusCode: 400,
		},
		{
			name:    "Such a news does not exist",
			pathId:  "tt",
			inputId: "TT",
			mockBehavior: func(s *mockservice.MockAppCountries, inputId string) {
				s.EXPECT().DeleteCountry(inputId).Return(MyErrors.DoesNotExist)
			},
			expectedStatusCode: 404,
		},
		{
			name:    "Server error",
			pathId:  "tt",
			inputId: "TT",
			mockBehavior: func(s *mockservice.MockAppCountries, inputId string) {
				s.EXPECT().DeleteCountry(inputId).Return(errors.New("server error"))
			},
			expectedStatusCode: 500,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			appService := mockservice.NewMockAppCountries(c)
			testCase.mockBehavior(appService, testCase.inputId)
			logger := logging.GetLoggerLogrus()
			serv := &services.Service{AppCountries: appService}
			handler := NewHandler(serv, logger)

			r := handler.InitRoutes()

			w := httptest.NewRecorder()

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/countries/%s", testCase.pathId), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
		})
	}
}
