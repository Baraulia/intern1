package repositories

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"tranee_service/MyErrors"
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
		return 0, fmt.Errorf("createUser: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()
	hobbiesId, err := CheckUserData(transaction, user)
	if err != nil {
		u.logger.Errorf("CreateUser: error while checking user data:%s", err)
		return 0, fmt.Errorf("createUser: error while checking user data:%w", err)
	}
	if len(hobbiesId) == 0 {
		u.logger.Errorf("CreateUser: there are no correct hobbies for the user")
		return 0, fmt.Errorf("createUser: there are no correct hobbies for the user")
	}
	user.Hobbies = hobbiesId
	query := "INSERT INTO users (name, email, description, country_id) values (?, ?, ?, ?)"
	result, err := transaction.Exec(query, user.Name, user.Email, user.Description, user.CountryId)
	if err != nil {
		u.logger.Errorf("CreateUser: error while insert user:%s", err)
		return 0, fmt.Errorf("createUser: error while insert user:%w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		u.logger.Errorf("CreateUser: error while getting insertId:%s", err)
		return 0, fmt.Errorf("createUser: error while getting insertId:%w", err)
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
		return 0, fmt.Errorf("createUser: error while insert users_hobbies:%w", err)
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
		return nil, fmt.Errorf("getUserById: can not builds the query into a SQL:%s", err)
	}
	row := u.db.QueryRow(query, args...)
	var bytesHobby []byte
	if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Description, &user.CountryId, &bytesHobby); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			u.logger.Errorf("GetUserById:object with this id does not exist")
			return nil, errors.Wrap(MyErrors.DoesNotExist, "getUserById")
		} else {
			u.logger.Errorf("Error while scanning for user:%s", err)
			return nil, fmt.Errorf("getUserById: repository error:%w", err)
		}
	}
	strHobby := string(bytesHobby[:])
	sliceHobby := strings.Split(strHobby, ",")
	for _, n := range sliceHobby {
		number, err := strconv.Atoi(n)
		if err != nil {
			u.logger.Errorf("Error while converting hobby`s id:%s", err)
			return nil, fmt.Errorf("getUserById: Error while converting hobby`s id:%w", err)
		}
		user.Hobbies = append(user.Hobbies, number)
	}
	return &user, nil
}

func (u *UserRepository) GetUsers(options *models.Options) ([]models.ResponseUser, int, error) {
	var users []models.ResponseUser
	var sel squirrel.SelectBuilder
	var pages int
	s := squirrel.Select("users.id, users.name, users.email, users.description, users.country_id, GROUP_CONCAT(users_hobbies.hobby_id) AS list").From("users").
		Join("users_hobbies on users.id = users_hobbies.user_id").GroupBy("users.id")
	if options.Page != 0 && options.Limit != 0 {
		sel = s.Limit(options.Limit).Offset((options.Page - 1) * options.Limit).OrderBy("users.id")
	} else {
		sel = s
		pages = 1
	}
	query, args, err := sel.ToSql()
	if err != nil {
		u.logger.Errorf("GetUsers: can not builds the query into a SQL:%s", err)
		return nil, 0, fmt.Errorf("getUsers: can not builds the query into a SQL:%s", err)
	}
	rows, err := u.db.Query(query, args...)
	if err != nil {
		u.logger.Errorf("GetUsers: can not executes a query:%s", err)
		return nil, 0, fmt.Errorf("getUsers: can not executes a query:%s", err)
	}
	defer rows.Close()
	for rows.Next() {
		var bytesHobby []byte
		var user models.ResponseUser
		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Description, &user.CountryId, &bytesHobby); err != nil {
			u.logger.Errorf("Error while scanning for user:%s", err)
			return nil, 0, fmt.Errorf("getUsers:repository error:%w", err)
		}
		strHobby := string(bytesHobby[:])
		sliceHobby := strings.Split(strHobby, ",")
		for _, n := range sliceHobby {
			number, err := strconv.Atoi(n)
			if err != nil {
				u.logger.Errorf("Error while converting hobby`s id:%s", err)
				return nil, 0, fmt.Errorf("getUsers: Error while converting hobby`s id:%w", err)
			}
			user.Hobbies = append(user.Hobbies, number)
		}
		users = append(users, user)
	}

	if pages != 1 {
		query = "SELECT CEILING(COUNT(*)/?) FROM users"
		row := u.db.QueryRow(query, options.Limit)
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
		return fmt.Errorf("changeUser: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()

	hobbiesId, err := CheckUserData(transaction, user)
	if err != nil {
		u.logger.Errorf("ChangeUser: error while checking user data:%s", err)
		return fmt.Errorf("сhangeUser: error while checking user data:%w", err)
	}
	if len(hobbiesId) == 0 {
		u.logger.Errorf("ChangeUser: there are no correct hobbies for the user")
		return fmt.Errorf("сhangeUser: there are no correct hobbies for the user")
	}
	user.Hobbies = hobbiesId

	query := "UPDATE users SET name = ?, email = ?, description = ?, country_id = ? WHERE id = ?"
	result, err := transaction.Exec(query, user.Name, user.Email, user.Description, user.CountryId, userId)
	if err != nil {
		u.logger.Errorf("ChangeUser: error while updating user:%s", err)
		return fmt.Errorf("changeUser: error while updating user:%w", err)
	}
	numberRows, err := result.RowsAffected()
	if err != nil {
		u.logger.Errorf("Error while getting number affected rows:%s", err)
		return fmt.Errorf("changeUser: error while getting number affected rows:%s", err)
	}
	if numberRows == 0 {
		u.logger.Errorf("ChangeUser:object with this id does not exist")
		return errors.Wrap(MyErrors.DoesNotExist, "changeUser")
	}

	query = "DELETE FROM users_hobbies WHERE user_id = ?"
	_, err = transaction.Exec(query, userId)
	if err != nil {
		u.logger.Errorf("ChangeUser: error while deleting bound relations:%s", err)
		return fmt.Errorf("changeUser: error whiledeleting bound relations:%w", err)
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
		return fmt.Errorf("changeUser: error while insert users_hobbies:%w", err)
	}
	return transaction.Commit()
}

func (u *UserRepository) DeleteUser(userId int) error {
	query := "DELETE from users WHERE id = ?"
	result, err := u.db.Exec(query, userId)
	if err != nil {
		u.logger.Errorf("DeleteUser: can not executes a query:%s", err)
		return fmt.Errorf("deleteUser: can not executes a query:%s", err)
	}
	numberRows, err := result.RowsAffected()
	if err != nil {
		u.logger.Errorf("Error while getting number affected rows:%s", err)
		return fmt.Errorf("deleteUser: error while getting number affected rows:%s", err)
	}
	if numberRows == 0 {
		u.logger.Errorf("DeleteUser:object with this id does not exist")
		return errors.Wrap(MyErrors.DoesNotExist, "deleteUser")
	}
	return nil
}

func CheckUserData(tr *sql.Tx, user *models.User) ([]int, error) {
	var exist bool
	var hobbiesId []int
	var userHobbies []int
	query := "SELECT EXISTS (SELECT 1 FROM countries WHERE id = ?)"
	row := tr.QueryRow(query, user.CountryId)
	if err := row.Scan(&exist); err != nil {
		return nil, err
	}
	if !exist {
		return nil, MyErrors.DoesNotExist
	}
	query = "SELECT id FROM hobbies"
	rows, err := tr.Query(query)
	if err != nil {
		return nil, fmt.Errorf("checkUserData: can not executes a query:%s", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("checkUserData:repository error:%w", err)
		}
		hobbiesId = append(hobbiesId, id)
	}

	for _, id := range user.Hobbies {
		for _, hobby := range hobbiesId {
			if hobby == id {
				userHobbies = append(userHobbies, id)
			}
		}
	}
	return userHobbies, nil
}
