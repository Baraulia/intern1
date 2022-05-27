// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_services is a generated GoMock package.
package mock_services

import (
	reflect "reflect"
	models "tranee_service/models"

	gomock "github.com/golang/mock/gomock"
)

// MockAppCountries is a mock of AppCountries interface.
type MockAppCountries struct {
	ctrl     *gomock.Controller
	recorder *MockAppCountriesMockRecorder
}

// MockAppCountriesMockRecorder is the mock recorder for MockAppCountries.
type MockAppCountriesMockRecorder struct {
	mock *MockAppCountries
}

// NewMockAppCountries creates a new mock instance.
func NewMockAppCountries(ctrl *gomock.Controller) *MockAppCountries {
	mock := &MockAppCountries{ctrl: ctrl}
	mock.recorder = &MockAppCountriesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAppCountries) EXPECT() *MockAppCountriesMockRecorder {
	return m.recorder
}

// ChangeCountry mocks base method.
func (m *MockAppCountries) ChangeCountry(country *models.ResponseCountry, countryId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeCountry", country, countryId)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeCountry indicates an expected call of ChangeCountry.
func (mr *MockAppCountriesMockRecorder) ChangeCountry(country, countryId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeCountry", reflect.TypeOf((*MockAppCountries)(nil).ChangeCountry), country, countryId)
}

// CreateCountry mocks base method.
func (m *MockAppCountries) CreateCountry(country *models.ResponseCountry) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCountry", country)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCountry indicates an expected call of CreateCountry.
func (mr *MockAppCountriesMockRecorder) CreateCountry(country interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCountry", reflect.TypeOf((*MockAppCountries)(nil).CreateCountry), country)
}

// DeleteCountry mocks base method.
func (m *MockAppCountries) DeleteCountry(countryId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCountry", countryId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCountry indicates an expected call of DeleteCountry.
func (mr *MockAppCountriesMockRecorder) DeleteCountry(countryId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCountry", reflect.TypeOf((*MockAppCountries)(nil).DeleteCountry), countryId)
}

// GetCountries mocks base method.
func (m *MockAppCountries) GetCountries(filters *models.Filters) ([]models.Country, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCountries", filters)
	ret0, _ := ret[0].([]models.Country)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetCountries indicates an expected call of GetCountries.
func (mr *MockAppCountriesMockRecorder) GetCountries(filters interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCountries", reflect.TypeOf((*MockAppCountries)(nil).GetCountries), filters)
}

// GetOneCountry mocks base method.
func (m *MockAppCountries) GetOneCountry(id string) (*models.Country, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOneCountry", id)
	ret0, _ := ret[0].(*models.Country)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOneCountry indicates an expected call of GetOneCountry.
func (mr *MockAppCountriesMockRecorder) GetOneCountry(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOneCountry", reflect.TypeOf((*MockAppCountries)(nil).GetOneCountry), id)
}

// LoadImages mocks base method.
func (m *MockAppCountries) LoadImages() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "LoadImages")
}

// LoadImages indicates an expected call of LoadImages.
func (mr *MockAppCountriesMockRecorder) LoadImages() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadImages", reflect.TypeOf((*MockAppCountries)(nil).LoadImages))
}

// MockAppUsers is a mock of AppUsers interface.
type MockAppUsers struct {
	ctrl     *gomock.Controller
	recorder *MockAppUsersMockRecorder
}

// MockAppUsersMockRecorder is the mock recorder for MockAppUsers.
type MockAppUsersMockRecorder struct {
	mock *MockAppUsers
}

// NewMockAppUsers creates a new mock instance.
func NewMockAppUsers(ctrl *gomock.Controller) *MockAppUsers {
	mock := &MockAppUsers{ctrl: ctrl}
	mock.recorder = &MockAppUsersMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAppUsers) EXPECT() *MockAppUsersMockRecorder {
	return m.recorder
}

// ChangeUser mocks base method.
func (m *MockAppUsers) ChangeUser(user *models.User, userId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeUser", user, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeUser indicates an expected call of ChangeUser.
func (mr *MockAppUsersMockRecorder) ChangeUser(user, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeUser", reflect.TypeOf((*MockAppUsers)(nil).ChangeUser), user, userId)
}

// CreateUser mocks base method.
func (m *MockAppUsers) CreateUser(user *models.User) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", user)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockAppUsersMockRecorder) CreateUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockAppUsers)(nil).CreateUser), user)
}

// DeleteUser mocks base method.
func (m *MockAppUsers) DeleteUser(userId int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockAppUsersMockRecorder) DeleteUser(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockAppUsers)(nil).DeleteUser), userId)
}

// GetHobbyByUserId mocks base method.
func (m *MockAppUsers) GetHobbyByUserId(userId int) ([]int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHobbyByUserId", userId)
	ret0, _ := ret[0].([]int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHobbyByUserId indicates an expected call of GetHobbyByUserId.
func (mr *MockAppUsersMockRecorder) GetHobbyByUserId(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHobbyByUserId", reflect.TypeOf((*MockAppUsers)(nil).GetHobbyByUserId), userId)
}

// GetUserById mocks base method.
func (m *MockAppUsers) GetUserById(userId int) (*models.ResponseUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserById", userId)
	ret0, _ := ret[0].(*models.ResponseUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserById indicates an expected call of GetUserById.
func (mr *MockAppUsersMockRecorder) GetUserById(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserById", reflect.TypeOf((*MockAppUsers)(nil).GetUserById), userId)
}

// GetUsers mocks base method.
func (m *MockAppUsers) GetUsers(filter *models.Options) ([]models.ResponseUser, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsers", filter)
	ret0, _ := ret[0].([]models.ResponseUser)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetUsers indicates an expected call of GetUsers.
func (mr *MockAppUsersMockRecorder) GetUsers(filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockAppUsers)(nil).GetUsers), filter)
}

// MockAppHobbies is a mock of AppHobbies interface.
type MockAppHobbies struct {
	ctrl     *gomock.Controller
	recorder *MockAppHobbiesMockRecorder
}

// MockAppHobbiesMockRecorder is the mock recorder for MockAppHobbies.
type MockAppHobbiesMockRecorder struct {
	mock *MockAppHobbies
}

// NewMockAppHobbies creates a new mock instance.
func NewMockAppHobbies(ctrl *gomock.Controller) *MockAppHobbies {
	mock := &MockAppHobbies{ctrl: ctrl}
	mock.recorder = &MockAppHobbiesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAppHobbies) EXPECT() *MockAppHobbiesMockRecorder {
	return m.recorder
}

// CreateHobby mocks base method.
func (m *MockAppHobbies) CreateHobby(hobby *models.Hobby) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateHobby", hobby)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateHobby indicates an expected call of CreateHobby.
func (mr *MockAppHobbiesMockRecorder) CreateHobby(hobby interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateHobby", reflect.TypeOf((*MockAppHobbies)(nil).CreateHobby), hobby)
}

// GetHobbies mocks base method.
func (m *MockAppHobbies) GetHobbies() ([]models.ResponseHobby, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHobbies")
	ret0, _ := ret[0].([]models.ResponseHobby)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHobbies indicates an expected call of GetHobbies.
func (mr *MockAppHobbiesMockRecorder) GetHobbies() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHobbies", reflect.TypeOf((*MockAppHobbies)(nil).GetHobbies))
}
