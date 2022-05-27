package repositories

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"strconv"
	"strings"
	"tranee_service/internal/logging"
	"tranee_service/models"
)

type UserRepository struct {
	db     *sql.DB
	logger logging.Logger
}

func NewUserRepository(db *sql.DB, logger logging.Logger) *UserRepository {
	return &UserRepository{db: db, logger: logger}
}

func (u *UserRepository) CreateUser(user *models.User) (int, error) {
	var userId int
	transaction, err := u.db.Begin()
	if err != nil {
		u.logger.Errorf("CreateUser: can not starts transaction:%s", err)
		return 0, fmt.Errorf("CreateUser: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()

	query := "INSERT INTO users (name, email, description, country_id) values (?, ?, ?, ?)"
	result, err := transaction.Exec(query, user.Name, user.Email, user.Description, user.CountryId)
	if err != nil {
		u.logger.Errorf("CreateUser: error while insert user:%s", err)
		return 0, fmt.Errorf("CreateUser: error while insert user:%w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		u.logger.Errorf("CreateUser: error while getting insertId:%s", err)
		return 0, fmt.Errorf("CreateUser: error while getting insertId:%w", err)
	}
	userId = int(id)
	query = "INSERT INTO users_hobbies (user_id, hobby_id) values "
	var values []interface{}
	for _, s := range user.Hobbies {
		values = append(values, userId, s)
		query += `(?,?),`
	}
	query = query[:len(query)-1]
	_, err = transaction.Exec(query, values...)
	if err != nil {
		u.logger.Errorf("CreateUser: error while insert users_hobbies:%s", err)
		return 0, fmt.Errorf("CreateUser: error while insert users_hobbies:%w", err)
	}
	return userId, transaction.Commit()
}

func (u *UserRepository) GetUserById(userId int) (*models.ResponseUser, error) {
	var user models.ResponseUser
	s := squirrel.Select("users.id, users.name, users.email, users.description, users.country_id, GROUP_CONCAT(users_hobbies.hobby_id) AS list").From("users").
		Join("users_hobbies on users.id = users_hobbies.user_id").GroupBy("users.id").Where("users.id = ?", userId)
	query, args, err := s.ToSql()
	if err != nil {
		u.logger.Errorf("GetUserById: can not builds the query into a SQL:%s", err)
		return nil, fmt.Errorf("GetUserById: can not builds the query into a SQL:%s", err)
	}
	row := u.db.QueryRow(query, args...)
	var bytesHobby []byte
	if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Description, &user.CountryId, &bytesHobby); err != nil {
		u.logger.Errorf("Error while scanning for user:%s", err)
		return nil, fmt.Errorf("GetUserById: repository error:%w", err)
	}
	strHobby := string(bytesHobby[:])
	sliceHobby := strings.Split(strHobby, ",")
	for _, n := range sliceHobby {
		number, err := strconv.Atoi(n)
		if err != nil {
			u.logger.Errorf("Error while converting hobby`s id:%s", err)
			return nil, fmt.Errorf("GetUserById: Error while converting hobby`s id:%w", err)
		}
		user.Hobbies = append(user.Hobbies, number)
	}
	return &user, nil
}

func (u *UserRepository) GetUsers(filter *models.Options) ([]models.ResponseUser, int, error) {
	var users []models.ResponseUser
	var sel squirrel.SelectBuilder
	var pages int
	s := squirrel.Select("users.id, users.name, users.email, users.description, users.country_id, GROUP_CONCAT(users_hobbies.hobby_id) AS list").From("users").
		Join("users_hobbies on users.id = users_hobbies.user_id").GroupBy("users.id")
	if filter.Page != 0 && filter.Limit != 0 {
		sel = s.Limit(filter.Limit).Offset((filter.Page - 1) * filter.Limit).OrderBy("users.id")
	} else {
		sel = s
		pages = 1
	}
	query, args, err := sel.ToSql()
	if err != nil {
		u.logger.Errorf("GetUsers: can not builds the query into a SQL:%s", err)
		return nil, 0, fmt.Errorf("GetUsers: can not builds the query into a SQL:%s", err)
	}
	rows, err := u.db.Query(query, args...)
	if err != nil {
		u.logger.Errorf("GetUsers: can not executes a query:%s", err)
		return nil, 0, fmt.Errorf("GetUsers: can not executes a query:%s", err)
	}
	defer rows.Close()
	for rows.Next() {
		var bytesHobby []byte
		var user models.ResponseUser
		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Description, &user.CountryId, &bytesHobby); err != nil {
			u.logger.Errorf("Error while scanning for user:%s", err)
			return nil, 0, fmt.Errorf("GetUsers:repository error:%w", err)
		}
		strHobby := string(bytesHobby[:])
		sliceHobby := strings.Split(strHobby, ",")
		for _, n := range sliceHobby {
			number, err := strconv.Atoi(n)
			if err != nil {
				u.logger.Errorf("Error while converting hobby`s id:%s", err)
				return nil, 0, fmt.Errorf("GetUsers: Error while converting hobby`s id:%w", err)
			}
			user.Hobbies = append(user.Hobbies, number)
		}
		users = append(users, user)
	}

	if pages != 1 {
		query = "SELECT CEILING(COUNT(*)/?) FROM users"
		row := u.db.QueryRow(query, filter.Limit)
		if err := row.Scan(&pages); err != nil {
			u.logger.Errorf("Error while scanning for pages:%s", err)
			return nil, 0, fmt.Errorf("error while scanning for pages:%s", err)
		}
	}
	return users, pages, nil
}

func (u *UserRepository) ChangeUser(user *models.User, userId int) error {
	transaction, err := u.db.Begin()
	if err != nil {
		u.logger.Errorf("ChangeUser: can not starts transaction:%s", err)
		return fmt.Errorf("ChangeUser: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()

	query := "UPDATE users SET name = ?, email = ?, description = ?, country_id = ? WHERE id = ?"
	result, err := transaction.Exec(query, user.Name, user.Email, user.Description, user.CountryId, userId)
	if err != nil {
		u.logger.Errorf("ChangeUser: error while updating user:%s", err)
		return fmt.Errorf("ChangeUser: error while updating user:%w", err)
	}
	numberRows, err := result.RowsAffected()
	if err != nil {
		u.logger.Errorf("Error while getting number affected rows:%s", err)
		return fmt.Errorf("ChangeUser: error while getting number affected rows:%s", err)
	}
	if numberRows == 0 {
		u.logger.Errorf("User with such Id does not exist", err)
		return fmt.Errorf("user with such Id does not exist")
	}

	query = "DELETE FROM users_hobbies WHERE user_id = ?"
	_, err = transaction.Exec(query, userId)
	if err != nil {
		u.logger.Errorf("ChangeUser: error while deleting bound relations:%s", err)
		return fmt.Errorf("ChangeUser: error whiledeleting bound relations:%w", err)
	}

	query = "INSERT INTO users_hobbies (user_id, hobby_id) values "
	var values []interface{}
	for _, s := range user.Hobbies {
		values = append(values, userId, s)
		query += `(?,?),`
	}
	query = query[:len(query)-1]
	_, err = transaction.Exec(query, values...)
	if err != nil {
		u.logger.Errorf("ChangeUser: error while insert users_hobbies:%s", err)
		return fmt.Errorf("ChangeUser: error while insert users_hobbies:%w", err)
	}
	return transaction.Commit()
}

func (u *UserRepository) DeleteUser(userId int) error {
	query := "DELETE from users WHERE id = ?"
	result, err := u.db.Exec(query, userId)
	if err != nil {
		u.logger.Errorf("DeleteUser: can not executes a query:%s", err)
		return fmt.Errorf("DeleteUser: can not executes a query:%s", err)
	}
	numberRows, err := result.RowsAffected()
	if err != nil {
		u.logger.Errorf("Error while getting number affected rows:%s", err)
		return fmt.Errorf("DeleteUser: error while getting number affected rows:%s", err)
	}
	if numberRows == 0 {
		u.logger.Errorf("User with such Id does not exist", err)
		return fmt.Errorf("user with such Id does not exist")
	}
	return nil
}
