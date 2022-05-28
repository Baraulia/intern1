package repositories

import (
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"tranee_service/internal/logging"
	"tranee_service/models"
)

func TestRepository_SaveInitialCountries(t *testing.T) {
	logger := logging.GetLoggerLogrus()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewRepository(db, logger)

	testTable := []struct {
		name           string
		mock           func(countries []models.Country)
		inputCountries []models.Country
		expectedError  bool
	}{
		{
			name: "OK",
			inputCountries: []models.Country{
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
			},
			mock: func(countries []models.Country) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"numberRows"}).AddRow(0)
				mock.ExpectQuery("SELECT COUNT").WillReturnRows(rows)
				mock.ExpectExec("INSERT INTO countries").WillReturnResult(driver.ResultNoRows)
				mock.ExpectCommit()
			},
			expectedError: false,
		},
		{
			name: "OK data already exists",
			inputCountries: []models.Country{
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
			},
			mock: func(countries []models.Country) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"numberRows"}).AddRow(1)
				mock.ExpectQuery("SELECT COUNT").WillReturnRows(rows)
				mock.ExpectCommit()
			},
			expectedError: false,
		},
		{
			name: "Data base error",
			inputCountries: []models.Country{
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
			},
			mock: func(countries []models.Country) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"numberRows"}).AddRow(0)
				mock.ExpectQuery("SELECT COUNT").WillReturnRows(rows)
				mock.ExpectExec("INSERT INTO countries").WillReturnError(errors.New("data base error"))
			},
			expectedError: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.inputCountries)
			err := r.SaveInitialCountries(tt.inputCountries)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_GetOneCountry(t *testing.T) {
	logger := logging.GetLoggerLogrus()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewRepository(db, logger)

	testTable := []struct {
		name           string
		mock           func(countryId string)
		inputId        string
		expectedResult *models.Country
		expectedError  bool
	}{
		{
			name:    "OK",
			inputId: "TT",
			mock: func(countryId string) {
				rows := sqlmock.NewRows([]string{"name", "full_name", "english_name", "alpha_2", "alpha_3", "iso", "location", "location_precise", "url"}).
					AddRow("test name", "test full name", "test ennglish name", "tt", "ttt", 1000, "test location", "test location precise", "")
				mock.ExpectQuery("SELECT name, full_name,").WillReturnRows(rows)
			},
			expectedResult: &models.Country{
				Name:            "test name",
				FullName:        "test full name",
				EnglishName:     "test ennglish name",
				Alpha2:          "tt",
				Alpha3:          "ttt",
				Iso:             1000,
				Location:        "test location",
				LocationPrecise: "test location precise",
				Url:             "",
			},
			expectedError: false,
		},
		{
			name:    "Data base error",
			inputId: "TT",
			mock: func(countryId string) {
				mock.ExpectQuery("SELECT name, full_name,").WillReturnError(errors.New("data base error"))
			},
			expectedError: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.inputId)
			country, err := r.GetOneCountry(tt.inputId)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, country)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_GetCountries(t *testing.T) {
	logger := logging.GetLoggerLogrus()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewRepository(db, logger)

	testTable := []struct {
		name           string
		mock           func(filter *models.Filters)
		inputFilter    *models.Filters
		expectedResult []models.Country
		expectedError  bool
	}{
		{
			name: "OK",
			inputFilter: &models.Filters{
				Page:  1,
				Limit: 2,
				Flag:  false,
			},
			mock: func(filter *models.Filters) {
				rows := sqlmock.NewRows([]string{"name", "full_name", "english_name", "alpha_2", "alpha_3", "iso", "location", "location_precise", "url"}).
					AddRow("test name", "test full name", "test english name", "tt", "ttt", 1000, "test location", "test location precise", "").
					AddRow("test name2", "test full name2", "test english name2", "tp", "tpt", 1001, "test location", "test location precise", "")
				mock.ExpectQuery("SELECT name, full_name,").WillReturnRows(rows)
				rows = sqlmock.NewRows([]string{"pages"}).AddRow(1)
				mock.ExpectQuery("SELECT CEILING").WillReturnRows(rows)
			},
			expectedResult: []models.Country{
				{
					Name:            "test name",
					FullName:        "test full name",
					EnglishName:     "test english name",
					Alpha2:          "tt",
					Alpha3:          "ttt",
					Iso:             1000,
					Location:        "test location",
					LocationPrecise: "test location precise",
					Url:             "",
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
					Url:             "",
				},
			},
			expectedError: false,
		},
		{
			name: "OK without pagination",
			inputFilter: &models.Filters{
				Page:  0,
				Limit: 0,
				Flag:  false,
			},
			mock: func(filter *models.Filters) {
				rows := sqlmock.NewRows([]string{"name", "full_name", "english_name", "alpha_2", "alpha_3", "iso", "location", "location_precise", "url"}).
					AddRow("test name", "test full name", "test english name", "tt", "ttt", 1000, "test location", "test location precise", "").
					AddRow("test name2", "test full name2", "test english name2", "tp", "tpt", 1001, "test location", "test location precise", "")
				mock.ExpectQuery("SELECT name, full_name,").WillReturnRows(rows)
			},
			expectedResult: []models.Country{
				{
					Name:            "test name",
					FullName:        "test full name",
					EnglishName:     "test english name",
					Alpha2:          "tt",
					Alpha3:          "ttt",
					Iso:             1000,
					Location:        "test location",
					LocationPrecise: "test location precise",
					Url:             "",
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
					Url:             "",
				},
			},
			expectedError: false,
		},
		{
			name: "OK with flag = true",
			inputFilter: &models.Filters{
				Page:  1,
				Limit: 2,
				Flag:  true,
			},
			mock: func(filter *models.Filters) {
				rows := sqlmock.NewRows([]string{"name", "full_name", "english_name", "alpha_2", "alpha_3", "iso", "location", "location_precise", "url"}).
					AddRow("test name", "test full name", "test english name", "tt", "ttt", 1000, "test location", "test location precise", "").
					AddRow("test name2", "test full name2", "test english name2", "tp", "tpt", 1001, "test location", "test location precise", "")
				mock.ExpectQuery("SELECT name, full_name,").WillReturnRows(rows)
				rows = sqlmock.NewRows([]string{"pages"}).AddRow(1)
				mock.ExpectQuery("SELECT CEILING").WillReturnRows(rows)
			},
			expectedResult: []models.Country{
				{
					Name:            "test name",
					FullName:        "test full name",
					EnglishName:     "test english name",
					Alpha2:          "tt",
					Alpha3:          "ttt",
					Iso:             1000,
					Location:        "test location",
					LocationPrecise: "test location precise",
					Url:             "",
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
					Url:             "",
				},
			},
			expectedError: false,
		},
		{
			name: "Data base error",
			inputFilter: &models.Filters{
				Page:  1,
				Limit: 2,
				Flag:  true,
			},
			mock: func(filter *models.Filters) {
				mock.ExpectQuery("SELECT name, full_name,").WillReturnError(errors.New("data base error"))
			},
			expectedError: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.inputFilter)
			countries, _, err := r.GetCountries(tt.inputFilter)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, countries)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_CreateCountry(t *testing.T) {
	logger := logging.GetLoggerLogrus()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewRepository(db, logger)

	testTable := []struct {
		name           string
		mock           func(country *models.ResponseCountry)
		inputCountry   *models.ResponseCountry
		expectedResult string
		expectedError  bool
	}{
		{
			name: "OK",
			inputCountry: &models.ResponseCountry{
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
			mock: func(country *models.ResponseCountry) {
				result := sqlmock.NewResult(1, 1)
				mock.ExpectExec("INSERT INTO countries").
					WithArgs(country.Name, country.FullName, country.EnglishName, country.Alpha2, country.Alpha3, country.Iso, country.Location, country.LocationPrecise, country.Url).
					WillReturnResult(result)
				rows := sqlmock.NewRows([]string{"insertId"}).AddRow("tt")
				mock.ExpectQuery("SELECT alpha_2 FROM").WithArgs(1).WillReturnRows(rows)
			},
			expectedResult: "tt",
			expectedError:  false,
		},
		{
			name: "Data base error",
			inputCountry: &models.ResponseCountry{
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
			mock: func(country *models.ResponseCountry) {
				mock.ExpectExec("INSERT INTO countries").
					WithArgs(country.Name, country.FullName, country.EnglishName, country.Alpha2, country.Alpha3, country.Iso, country.Location, country.LocationPrecise, country.Url).
					WillReturnError(errors.New("data base error"))
			},
			expectedResult: "",
			expectedError:  true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.inputCountry)
			id, err := r.CreateCountry(tt.inputCountry)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, id)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_ChangeCountry(t *testing.T) {
	logger := logging.GetLoggerLogrus()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewRepository(db, logger)

	testTable := []struct {
		name          string
		mock          func(country *models.ResponseCountry, countryId string)
		inputCountry  *models.ResponseCountry
		inputId       string
		expectedError bool
	}{
		{
			name: "OK",
			inputCountry: &models.ResponseCountry{
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
			inputId: "TT",
			mock: func(country *models.ResponseCountry, countryId string) {
				result := sqlmock.NewResult(1, 1)
				mock.ExpectExec("UPDATE IGNORE countries").
					WithArgs(country.Name, country.FullName, country.EnglishName, country.Alpha2, country.Alpha3, country.Iso, country.Location, country.LocationPrecise, country.Url, countryId, countryId).
					WillReturnResult(result)
			},
			expectedError: false,
		},
		{
			name: "Data base error",
			inputCountry: &models.ResponseCountry{
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
			inputId: "TT",
			mock: func(country *models.ResponseCountry, countryId string) {
				mock.ExpectExec("UPDATE IGNORE countries").
					WithArgs(country.Name, country.FullName, country.EnglishName, country.Alpha2, country.Alpha3, country.Iso, country.Location, country.LocationPrecise, country.Url, countryId, countryId).
					WillReturnError(errors.New("data base error"))
			},
			expectedError: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.inputCountry, tt.inputId)
			err := r.ChangeCountry(tt.inputCountry, tt.inputId)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_DeleteCountry(t *testing.T) {
	logger := logging.GetLoggerLogrus()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewRepository(db, logger)

	testTable := []struct {
		name          string
		mock          func(countryId string)
		inputId       string
		expectedError bool
	}{
		{
			name:    "OK",
			inputId: "TT",
			mock: func(countryId string) {
				result := sqlmock.NewResult(1, 1)
				mock.ExpectExec("DELETE FROM countries").WithArgs(countryId, countryId).WillReturnResult(result)
			},
			expectedError: false,
		},
		{
			name:    "Data base error",
			inputId: "TT",
			mock: func(countryId string) {
				mock.ExpectExec("DELETE FROM countries").WithArgs(countryId, countryId).WillReturnError(errors.New("data base error"))
			},
			expectedError: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.inputId)
			err := r.DeleteCountry(tt.inputId)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_CheckCountryId(t *testing.T) {
	logger := logging.GetLoggerLogrus()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewRepository(db, logger)

	testTable := []struct {
		name          string
		mock          func(countryId string)
		inputId       string
		expectedError bool
	}{
		{
			name:    "OK",
			inputId: "TT",
			mock: func(countryId string) {
				rows := sqlmock.NewRows([]string{"exist"}).AddRow(true)
				mock.ExpectQuery("SELECT EXISTS ").WithArgs(countryId, countryId).WillReturnRows(rows)
			},
			expectedError: false,
		},
		{
			name:    "Data base error",
			inputId: "TT",
			mock: func(countryId string) {
				mock.ExpectQuery("SELECT EXISTS ").WithArgs(countryId, countryId).WillReturnError(errors.New("data base error"))
			},
			expectedError: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.inputId)
			err := r.CheckCountryId(tt.inputId)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_LoadImages(t *testing.T) {
	logger := logging.GetLoggerLogrus()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewRepository(db, logger)

	testTable := []struct {
		name           string
		mock           func(countries []models.Country)
		inputCountries []models.Country
		expectedError  bool
	}{
		{
			name: "OK",
			inputCountries: []models.Country{
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
			},
			mock: func(countries []models.Country) {
				mock.ExpectExec("UPDATE countries").
					WithArgs(countries[0].EnglishName, countries[0].Url, countries[1].EnglishName, countries[1].Url, countries[0].EnglishName, countries[1].EnglishName).
					WillReturnResult(driver.ResultNoRows)
			},
			expectedError: false,
		},
		{
			name: "Data base error",
			inputCountries: []models.Country{
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
			},
			mock: func(countries []models.Country) {
				mock.ExpectExec("UPDATE countries").
					WithArgs(countries[0].EnglishName, countries[0].Url, countries[1].EnglishName, countries[1].Url, countries[0].EnglishName, countries[1].EnglishName).
					WillReturnError(errors.New("data base error"))
			},
			expectedError: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.inputCountries)
			err := r.LoadImages(tt.inputCountries)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
