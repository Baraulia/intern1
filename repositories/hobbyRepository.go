package repositories

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"tranee_service/MyErrors"
	"tranee_service/internal/logging"
	"tranee_service/models"
)

type HobbyRepository struct {
	db     *sql.DB
	logger logging.Logger
}

func NewHobbyRepository(db *sql.DB, logger logging.Logger) *HobbyRepository {
	return &HobbyRepository{db: db, logger: logger}
}

func (h *HobbyRepository) GetHobbyByUserId(userId int) ([]int, error) {
	var ids []int
	query := "SELECT hobby_id from users_hobbies WHERE user_id = ?"
	rows, err := h.db.Query(query, userId)
	if err != nil {
		h.logger.Errorf("GetHobbyByUserId: can not executes a query:%s", err)
		return nil, fmt.Errorf("getHobbyByUserId: can not executes a query:%s", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			h.logger.Errorf("Error while scanning for hobby id:%s", err)
			return nil, fmt.Errorf("getHobbyByUserId:repository error:%w", err)
		}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		h.logger.Errorf("GetHobbyByUserId:object with this id does not exist")
		return nil, errors.Wrap(MyErrors.DoesNotExist, "getHobbyByUserId")
	}
	return ids, nil
}

func (h *HobbyRepository) CreateHobby(hobby *models.Hobby) (int, error) {
	var id int
	query := "INSERT INTO hobbies (name) VALUES (?)"
	result, err := h.db.Exec(query, hobby.Name)
	if err != nil {
		h.logger.Errorf("CreateHobby: can not adding new hobby:%s", err)
		return 0, fmt.Errorf("createHobby: can not adding new hobby:%w", err)
	}
	insertId, err := result.LastInsertId()
	if err != nil {
		h.logger.Errorf("CreateHobby: error while getting insertId:%s", err)
		return 0, fmt.Errorf("createHobby: error while getting insertId:%w", err)
	}
	id = int(insertId)
	return id, nil
}

func (h *HobbyRepository) GetHobbies() ([]models.ResponseHobby, error) {
	var hobbies []models.ResponseHobby
	query := "SELECT id, name FROM hobbies"
	rows, err := h.db.Query(query)
	if err != nil {
		h.logger.Errorf("GetHobbies: can not executes a query:%s", err)
		return nil, fmt.Errorf("getHobbies: can not executes a query:%s", err)
	}
	defer rows.Close()
	for rows.Next() {
		var hobby models.ResponseHobby
		if err := rows.Scan(&hobby.Id, &hobby.Name); err != nil {
			h.logger.Errorf("Error while scanning for hobby:%s", err)
			return nil, fmt.Errorf("getHobbies:repository error:%w", err)
		}
		hobbies = append(hobbies, hobby)
	}
	return hobbies, nil
}
