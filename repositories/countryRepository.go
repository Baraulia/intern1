package repositories

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"tranee_service/internal/logging"
	"tranee_service/models"
)

type CountryRepository struct {
	db     *sql.DB
	logger logging.Logger
}

func NewCountryRepository(db *sql.DB, logger logging.Logger) *CountryRepository {
	return &CountryRepository{db: db, logger: logger}
}

func (c *CountryRepository) SaveInitialCountries(countries []models.Country) error {
	var numberRows int
	transaction, err := c.db.Begin()
	if err != nil {
		c.logger.Errorf("SaveInitialCountries: can not starts transaction:%s", err)
		return fmt.Errorf("SaveInitialCountries: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()
	query := `SELECT COUNT(*) FROM countries`
	row := transaction.QueryRow(query)
	if err := row.Scan(&numberRows); err != nil {
		c.logger.Errorf("Error while scanning for numberRows:%s", err)
		return fmt.Errorf("error while scanning for numberRows:%s", err)
	}
	if numberRows == 0 {
		query = `INSERT INTO countries (name, full_name, english_name, alpha_2, alpha_3, iso, location, location_precise, url) values `
		var values []interface{}
		for _, s := range countries {
			values = append(values, s.Name, s.FullName, s.EnglishName, s.Alpha2, s.Alpha3, s.Iso, s.Location, s.LocationPrecise, s.Url)
			query += `(?,?,?,?,?,?,?,?,?),`
		}
		query = query[:len(query)-1] // remove the trailing comma
		_, err = transaction.Exec(query, values...)
		if err != nil {
			c.logger.Errorf("SaveInitialCountries: error while insert countries:%s", err)
			return fmt.Errorf("SaveInitialCountries: error while insert countriesr:%w", err)
		}
	}
	return transaction.Commit()
}

func (c *CountryRepository) GetOneCountry(id string) (*models.Country, error) {
	var country models.Country
	query := "SELECT name, full_name, english_name, alpha_2, alpha_3, iso, location, location_precise, url FROM countries WHERE alpha_2 = ? OR alpha_3 = ?"
	row := c.db.QueryRow(query, id, id)
	if err := row.Scan(&country.Name, &country.FullName, &country.EnglishName, &country.Alpha2, &country.Alpha3, &country.Iso, &country.Location, &country.LocationPrecise, &country.Url); err != nil {
		c.logger.Errorf("Error while scanning for country:%s", err)
		return nil, fmt.Errorf("GetOneCountry: repository error:%w", err)
	}
	return &country, nil
}

func (c *CountryRepository) GetCountries(filters *models.Filters) ([]models.Country, int, error) {
	var countries []models.Country
	var pages int
	var sel squirrel.SelectBuilder
	s := squirrel.Select("name", "full_name", "english_name", "alpha_2", "alpha_3", "iso", "location", "location_precise", "url").From("countries")
	if filters.Flag {
		sel = s.Where(squirrel.Eq{"url": ""})
	}
	if filters.Page != 0 && filters.Limit != 0 {
		sel = s.Limit(filters.Limit).Offset((filters.Page - 1) * filters.Limit).OrderBy("alpha_2")
	} else {
		sel = s
		pages = 1
	}
	query, args, err := sel.ToSql()
	if err != nil {
		c.logger.Errorf("GetCountries: can not builds the query into a SQL:%s", err)
		return nil, 0, fmt.Errorf("GetCountries: can not builds the query into a SQL:%s", err)
	}
	rows, err := c.db.Query(query, args...)
	if err != nil {
		c.logger.Errorf("GetCountries: can not executes a query:%s", err)
		return nil, 0, fmt.Errorf("GetCountries: can not executes a query:%s", err)
	}
	defer rows.Close()
	for rows.Next() {
		var country models.Country
		if err := rows.Scan(&country.Name, &country.FullName, &country.EnglishName, &country.Alpha2, &country.Alpha3, &country.Iso, &country.Location, &country.LocationPrecise, &country.Url); err != nil {
			c.logger.Errorf("Error while scanning for country:%s", err)
			return nil, 0, fmt.Errorf("GetCountries:repository error:%w", err)
		}
		countries = append(countries, country)
	}
	if pages != 1 {
		query = "SELECT CEILING(COUNT(*)/?) FROM countries"
		row := c.db.QueryRow(query, filters.Limit)
		if err := row.Scan(&pages); err != nil {
			c.logger.Errorf("Error while scanning for pages:%s", err)
			return nil, 0, fmt.Errorf("error while scanning for pages:%s", err)
		}
	}
	return countries, pages, nil
}

func (c *CountryRepository) CreateCountry(country *models.ResponseCountry) (string, error) {
	var id string
	query := "INSERT INTO countries (name, full_name, english_name, alpha_2, alpha_3, iso, location, location_precise, url) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := c.db.Exec(query, country.Name, country.FullName, country.EnglishName, country.Alpha2, country.Alpha3, country.Iso, country.Location, country.LocationPrecise, country.Url)
	if err != nil {
		c.logger.Errorf("CreateCountry: can not adding new country:%s", err)
		return "", fmt.Errorf("CreateCountry: can not adding new country:%w", err)
	}
	insertId, err := result.LastInsertId()
	if err != nil {
		c.logger.Errorf("CreateCountry: error while getting insertId:%s", err)
		return "", fmt.Errorf("CreateCountry: error while getting insertId:%w", err)
	}
	query = "SELECT alpha_2 FROM countries WHERE id = ?"
	row := c.db.QueryRow(query, insertId)
	if err = row.Scan(&id); err != nil {
		c.logger.Errorf("CreateCountry: error while scanning for countryId:%s", err)
		return "", fmt.Errorf("CreateCountry: error while scanning for countryId:%w", err)
	}
	return id, nil
}

func (c *CountryRepository) ChangeCountry(country *models.ResponseCountry, countryId string) error {
	query := "UPDATE IGNORE countries SET name = ?, full_name = ?, english_name = ?, alpha_2 = ?, alpha_3 = ?, iso = ?, location = ?, location_precise = ?, url = ? WHERE alpha_2 = ? OR alpha_3 = ?"
	_, err := c.db.Exec(query, country.Name, country.FullName, country.EnglishName, country.Alpha2, country.Alpha3, country.Iso, country.Location, country.LocationPrecise, country.Url, countryId, countryId)
	if err != nil {
		c.logger.Errorf("ChangeCountry: error while updating country:%s", err)
		return fmt.Errorf("ChangeCountry: error while updating country:%w", err)
	}
	return nil
}

func (c *CountryRepository) DeleteCountry(countryId string) error {
	query := "DELETE FROM countries WHERE alpha_2 = ? OR alpha_3 = ?"
	result, err := c.db.Exec(query, countryId, countryId)
	if err != nil {
		c.logger.Errorf("Error while scanning for countryId:%s", err)
		return fmt.Errorf("DeleteCountry: error while scanning for countryId:%s", err)
	}
	numberRows, err := result.RowsAffected()
	if err != nil {
		c.logger.Errorf("Error while getting number affected rows:%s", err)
		return fmt.Errorf("DeleteCountry: error while getting number affected rows:%s", err)
	}
	if numberRows == 0 {
		c.logger.Errorf("Country with such Id does not exist", err)
		return fmt.Errorf("country with such Id does not exist")
	}
	return nil
}

func (c *CountryRepository) CheckCountryId(countryId string) error {
	var exist bool
	query := "SELECT EXISTS (select 1 from countries where alpha_2 = ? OR alpha_3 = ?)"
	row := c.db.QueryRow(query, countryId, countryId)
	if err := row.Scan(&exist); err != nil {
		c.logger.Errorf("Error while scanning for existing country:%s", err)
		return err
	}
	if !exist {
		return fmt.Errorf("such a country does not exist")
	}
	return nil
}

func (c *CountryRepository) LoadImages(countries []models.Country) error {
	query := `UPDATE countries SET url = CASE english_name `
	query2 := " "
	var values []interface{}
	var values2 []interface{}
	for _, s := range countries {
		values = append(values, s.EnglishName, s.Url)
		query += `WHEN ? THEN ? `
		query2 += `?,`
		values2 = append(values2, s.EnglishName)
	}
	query += `ELSE url END ` + `WHERE english_name IN (` + query2[:len(query2)-1] + `)`
	for _, value := range values2 {
		values = append(values, value)
	}

	_, err := c.db.Exec(query, values...)
	if err != nil {
		c.logger.Errorf("LoadImages: error while insert flag url:%s", err)
		return fmt.Errorf("LoadImages: error while insert flag url:%w", err)
	}
	return nil
}
