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

func TestRepository_CreateUser(t *testing.T) {
	logger := logging.GetLoggerLogrus()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewRepository(db, logger)

	testTable := []struct {
		name           string
		mock           func(user *models.User)
		inputUser      *models.User
		expectedResult int
		expectedError  bool
	}{
		{
			name: "OK",
			inputUser: &models.User{
				Name:        "testName",
				Email:       "test@test.ru",
				Description: "test desc",
				CountryId:   1,
				Hobbies:     []int{1, 2},
			},
			mock: func(user *models.User) {
				mock.ExpectBegin()
				result := sqlmock.NewResult(1, 1)
				mock.ExpectExec("INSERT INTO users").
					WithArgs(user.Name, user.Email, user.Description, user.CountryId).
					WillReturnResult(result)
				mock.ExpectExec("INSERT INTO users_hobbies").WithArgs(1, 1, 1, 2).
					WillReturnResult(driver.ResultNoRows)
				mock.ExpectCommit()
			},
			expectedResult: 1,
			expectedError:  false,
		},
		{
			name: "Data base error",
			inputUser: &models.User{
				Name:        "testName",
				Email:       "test@test.ru",
				Description: "test desc",
				CountryId:   1,
				Hobbies:     []int{1, 2},
			},
			mock: func(user *models.User) {
				mock.ExpectBegin()
				result := sqlmock.NewResult(1, 1)
				mock.ExpectExec("INSERT INTO users").
					WithArgs(user.Name, user.Email, user.Description, user.CountryId).
					WillReturnResult(result)
				mock.ExpectExec("INSERT INTO users_hobbies").WithArgs(1, 1, 1, 2).
					WillReturnError(errors.New("data base error"))
				mock.ExpectRollback()
			},
			expectedResult: 0,
			expectedError:  true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.inputUser)
			id, err := r.CreateUser(tt.inputUser)
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

func TestRepository_ChangeUser(t *testing.T) {
	logger := logging.GetLoggerLogrus()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewRepository(db, logger)

	testTable := []struct {
		name          string
		mock          func(user *models.User, userId int)
		inputUser     *models.User
		inputId       int
		expectedError bool
	}{
		{
			name: "OK",
			inputUser: &models.User{
				Name:        "testName",
				Email:       "test@test.ru",
				Description: "test desc",
				CountryId:   1,
				Hobbies:     []int{1, 2},
			},
			inputId: 1,
			mock: func(user *models.User, userId int) {
				mock.ExpectBegin()
				result := sqlmock.NewResult(1, 1)
				mock.ExpectExec("UPDATE users SET").
					WithArgs(user.Name, user.Email, user.Description, user.CountryId, userId).
					WillReturnResult(result)
				mock.ExpectExec("DELETE FROM users_hobbies").WithArgs(1).
					WillReturnResult(driver.ResultNoRows)
				mock.ExpectExec("INSERT INTO users_hobbies").WithArgs(1, 1, 1, 2).
					WillReturnResult(driver.ResultNoRows)
				mock.ExpectCommit()
			},
			expectedError: false,
		},
		{
			name: "User with such Id does not exist",
			inputUser: &models.User{
				Name:        "testName",
				Email:       "test@test.ru",
				Description: "test desc",
				CountryId:   1,
				Hobbies:     []int{1, 2},
			},
			inputId: 1,
			mock: func(user *models.User, userId int) {
				mock.ExpectBegin()
				result := sqlmock.NewResult(0, 0)
				mock.ExpectExec("UPDATE users SET").
					WithArgs(user.Name, user.Email, user.Description, user.CountryId, userId).
					WillReturnResult(result)
				mock.ExpectRollback()
			},
			expectedError: true,
		},
		{
			name: "Data base error",
			inputUser: &models.User{
				Name:        "testName",
				Email:       "test@test.ru",
				Description: "test desc",
				CountryId:   1,
				Hobbies:     []int{1, 2},
			},
			inputId: 1,
			mock: func(user *models.User, userId int) {
				mock.ExpectBegin()
				result := sqlmock.NewResult(1, 1)
				mock.ExpectExec("UPDATE users SET").
					WithArgs(user.Name, user.Email, user.Description, user.CountryId, userId).
					WillReturnResult(result)
				mock.ExpectExec("DELETE FROM users_hobbies").WithArgs(1).
					WillReturnResult(driver.ResultNoRows)
				mock.ExpectExec("INSERT INTO users_hobbies").WithArgs(1, 1, 1, 2).
					WillReturnError(errors.New("data base error"))
				mock.ExpectRollback()
			},
			expectedError: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.inputUser, tt.inputId)
			err := r.ChangeUser(tt.inputUser, tt.inputId)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_DeleteUser(t *testing.T) {
	logger := logging.GetLoggerLogrus()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewRepository(db, logger)

	testTable := []struct {
		name          string
		mock          func(userId int)
		inputId       int
		expectedError bool
	}{
		{
			name:    "OK",
			inputId: 1,
			mock: func(userId int) {
				result := sqlmock.NewResult(1, 1)
				mock.ExpectExec("DELETE from users ").WithArgs(userId).
					WillReturnResult(result)
			},
			expectedError: false,
		},
		{
			name:    "User with such Id does not exist",
			inputId: 1,
			mock: func(userId int) {
				result := sqlmock.NewResult(0, 0)
				mock.ExpectExec("DELETE from users ").WithArgs(userId).
					WillReturnResult(result)
			},
			expectedError: true,
		},
		{
			name:    "Data base error",
			inputId: 1,
			mock: func(userId int) {
				mock.ExpectExec("DELETE from users ").WithArgs(userId).WillReturnError(errors.New("data base error"))
			},
			expectedError: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.inputId)
			err := r.DeleteUser(tt.inputId)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRepository_GetUserById(t *testing.T) {
	logger := logging.GetLoggerLogrus()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewRepository(db, logger)

	testTable := []struct {
		name           string
		mock           func(userId int)
		inputId        int
		expectedResult *models.ResponseUser
		expectedError  bool
	}{
		{
			name:    "OK",
			inputId: 1,
			mock: func(userId int) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "description", "countryId", "list"}).
					AddRow(1, "test name", "test email", "test desc", 1, []byte("1"+","+"2"))
				mock.ExpectQuery("SELECT users.id, ").WithArgs(userId).
					WillReturnRows(rows)
			},
			expectedResult: &models.ResponseUser{
				Id:          1,
				Name:        "test name",
				Email:       "test email",
				Description: "test desc",
				CountryId:   1,
				Hobbies:     []int{1, 2},
			},
			expectedError: false,
		},
		{
			name:    "Data base error",
			inputId: 1,
			mock: func(userId int) {
				mock.ExpectQuery("SELECT users.id ").WithArgs(userId).WillReturnError(errors.New("data base error"))
			},
			expectedError: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.inputId)
			country, err := r.GetUserById(tt.inputId)
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

func TestRepository_GetUsers(t *testing.T) {
	logger := logging.GetLoggerLogrus()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewRepository(db, logger)

	testTable := []struct {
		name           string
		mock           func(options *models.Options)
		inputOptions   *models.Options
		expectedResult []models.ResponseUser
		expectedError  bool
	}{
		{
			name: "OK",
			inputOptions: &models.Options{
				Page:  1,
				Limit: 2,
			},
			mock: func(options *models.Options) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "description", "countryId", "list"}).
					AddRow(1, "test name", "test email", "test desc", 1, []byte("1"+","+"2")).
					AddRow(2, "test name2", "test email2", "test desc2", 1, []byte("1"+","+"2"))
				mock.ExpectQuery("SELECT users.id, ").WillReturnRows(rows)
				rows = sqlmock.NewRows([]string{"pages"}).AddRow(1)
				mock.ExpectQuery("SELECT CEILING").WithArgs(options.Limit).WillReturnRows(rows)
			},
			expectedResult: []models.ResponseUser{
				{
					Id:          1,
					Name:        "test name",
					Email:       "test email",
					Description: "test desc",
					CountryId:   1,
					Hobbies:     []int{1, 2},
				},
				{
					Id:          2,
					Name:        "test name2",
					Email:       "test email2",
					Description: "test desc2",
					CountryId:   1,
					Hobbies:     []int{1, 2},
				},
			},
			expectedError: false,
		},
		{
			name: "OK without pagination",
			inputOptions: &models.Options{
				Page:  0,
				Limit: 0,
			},
			mock: func(options *models.Options) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "description", "countryId", "list"}).
					AddRow(1, "test name", "test email", "test desc", 1, []byte("1"+","+"2")).
					AddRow(2, "test name2", "test email2", "test desc2", 1, []byte("1"+","+"2"))
				mock.ExpectQuery("SELECT users.id, ").WillReturnRows(rows)
			},
			expectedResult: []models.ResponseUser{
				{
					Id:          1,
					Name:        "test name",
					Email:       "test email",
					Description: "test desc",
					CountryId:   1,
					Hobbies:     []int{1, 2},
				},
				{
					Id:          2,
					Name:        "test name2",
					Email:       "test email2",
					Description: "test desc2",
					CountryId:   1,
					Hobbies:     []int{1, 2},
				},
			},
			expectedError: false,
		},
		{
			name: "Data base error",
			inputOptions: &models.Options{
				Page:  1,
				Limit: 2,
			},
			mock: func(filter *models.Options) {
				mock.ExpectQuery("SELECT users.id ").WillReturnError(errors.New("data base error"))
			},
			expectedError: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.inputOptions)
			countries, _, err := r.GetUsers(tt.inputOptions)
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
