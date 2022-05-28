package repositories

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"tranee_service/internal/logging"
	"tranee_service/models"
)

func TestGetHobbies(t *testing.T) {
	logger := logging.GetLoggerLogrus()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewRepository(db, logger)

	testTable := []struct {
		name           string
		mock           func()
		expectedResult []models.ResponseHobby
		expectedError  bool
	}{
		{
			name: "OK",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "test name").AddRow(2, "test name2")
				mock.ExpectQuery("SELECT id, name ").WillReturnRows(rows)
			},
			expectedResult: []models.ResponseHobby{
				{
					Id:   1,
					Name: "test name",
				},
				{
					Id:   2,
					Name: "test name2",
				},
			},
			expectedError: false,
		},
		{
			name: "Data base error",
			mock: func() {
				mock.ExpectQuery("SELECT id, name ").WillReturnError(errors.New("data base error"))
			},
			expectedError: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			hobbies, err := r.GetHobbies()
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, hobbies)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetHobbyByUserId(t *testing.T) {
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
		expectedResult []int
		expectedError  bool
	}{
		{
			name:    "OK",
			inputId: 1,
			mock: func(userId int) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow(1).AddRow(2)
				mock.ExpectQuery("SELECT hobby_id ").WillReturnRows(rows)
			},
			expectedResult: []int{1, 2},
			expectedError:  false,
		},
		{
			name:    "Data base error",
			inputId: 1,
			mock: func(userId int) {
				mock.ExpectQuery("SELECT hobby_id ").WillReturnError(errors.New("data base error"))
			},
			expectedError: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.inputId)
			hobbies, err := r.GetHobbyByUserId(tt.inputId)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, hobbies)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreateHobby(t *testing.T) {
	logger := logging.GetLoggerLogrus()
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	r := NewRepository(db, logger)

	testTable := []struct {
		name           string
		mock           func(hobby *models.Hobby)
		inputHobby     *models.Hobby
		expectedResult int
		expectedError  bool
	}{
		{
			name:       "OK",
			inputHobby: &models.Hobby{Name: "testName"},
			mock: func(hobby *models.Hobby) {
				result := sqlmock.NewResult(1, 1)
				mock.ExpectExec("INSERT INTO hobbies ").WillReturnResult(result)
			},
			expectedResult: 1,
			expectedError:  false,
		},
		{
			name:       "Data base error",
			inputHobby: &models.Hobby{Name: "testName"},
			mock: func(hobby *models.Hobby) {
				mock.ExpectExec("INSERT INTO hobbies ").WillReturnError(errors.New("data base error"))
			},
			expectedError: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.inputHobby)
			hobbies, err := r.CreateHobby(tt.inputHobby)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, hobbies)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
