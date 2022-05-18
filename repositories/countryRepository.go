package repositories

import (
	"database/sql"
	"fmt"
	"tranee_service/models"
)

type CountryRepository struct {
	db     *sql.DB
	logger models.Logger
}

func NewCountryRepository(db *sql.DB, logger models.Logger) *CountryRepository {
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
	if numberRows > 0 {
		return transaction.Commit()
	} else {
		query = `INSERT INTO countries (name, full_name, english_name, alpha_2, alpha_3, iso, location, location_precise, url) values `
		var values []interface{}
		for _, s := range countries {
			values = append(values, s.Name, s.FullName, s.EnglishName, s.Alpha2, s.Alpha3, s.Iso, s.Location, s.LocationPrecise, s.Url)

			numFields := 9 // the number of fields you are inserting

			query += `(`
			for j := 0; j < numFields; j++ {
				query += `?` + `,`
			}
			query = query[:len(query)-1] + `),`
		}
		query = query[:len(query)-1] // remove the trailing comma
		_, err = transaction.Exec(query, values...)
		if err != nil {
			c.logger.Errorf("SaveInitialCountries: error while insert countries:%s", err)
			return fmt.Errorf("SaveInitialCountries: error while insert countriesr:%w", err)
		}
		return transaction.Commit()
	}
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

func (c *CountryRepository) GetCountries(page int, limit int) ([]models.Country, int, error) {
	var countries []models.Country
	var pages int
	query := "SELECT name, full_name, english_name, alpha_2, alpha_3, iso, location, location_precise, url FROM countries ORDER BY alpha_2 LIMIT ? OFFSET ?"
	rows, err := c.db.Query(query, limit, (page-1)*limit)
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
	query = "SELECT CEILING(COUNT(*)/?) FROM countries"
	row := c.db.QueryRow(query, limit)
	if err := row.Scan(&pages); err != nil {
		c.logger.Errorf("Error while scanning for pages:%s", err)
		return nil, 0, fmt.Errorf("error while scanning for pages:%s", err)
	}
	return countries, pages, nil
}

func (c *CountryRepository) GetCountriesWithoutPagination() ([]models.Country, int, error) {
	var countries []models.Country
	query := "SELECT name, full_name, english_name, alpha_2, alpha_3, iso, location, location_precise, url FROM countries ORDER BY alpha_2"
	rows, err := c.db.Query(query)
	if err != nil {
		c.logger.Errorf("GetCountriesWithoutPagination: can not executes a query:%s", err)
		return nil, 0, fmt.Errorf("GetCountriesWithoutPagination: can not executes a query:%s", err)
	}
	defer rows.Close()
	for rows.Next() {
		var country models.Country
		if err := rows.Scan(&country.Name, &country.FullName, &country.EnglishName, &country.Alpha2, &country.Alpha3, &country.Iso, &country.Location, &country.LocationPrecise, &country.Url); err != nil {
			c.logger.Errorf("Error while scanning for country:%s", err)
			return nil, 0, fmt.Errorf("GetCountriesWithoutPagination:repository error:%w", err)
		}
		countries = append(countries, country)
	}
	return countries, 1, nil
}

func (c *CountryRepository) CreateCountry(country *models.ResponseCountry) (string, error) {
	transaction, err := c.db.Begin()
	if err != nil {
		c.logger.Errorf("CreateCountry: can not starts transaction:%s", err)
		return "", fmt.Errorf("CreateCountry: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()
	var id string
	query := "INSERT INTO countries (name, full_name, english_name, alpha_2, alpha_3, iso, location, location_precise, url) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := transaction.Exec(query, country.Name, country.FullName, country.EnglishName, country.Alpha2, country.Alpha3, country.Iso, country.Location, country.LocationPrecise, country.Url)
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
	row := transaction.QueryRow(query, insertId)
	if err = row.Scan(&id); err != nil {
		c.logger.Errorf("CreateCountry: error while scanning for countryId:%s", err)
		return "", fmt.Errorf("CreateCountry: error while scanning for countryId:%w", err)
	}
	return id, transaction.Commit()
}

func (c *CountryRepository) ChangeCountry(country *models.ResponseCountry, countryId string) error {
	transaction, err := c.db.Begin()
	if err != nil {
		c.logger.Errorf("ChangeCountry: can not starts transaction:%s", err)
		return fmt.Errorf("ChangeCountry: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()
	query := "UPDATE IGNORE countries SET name = ?, full_name = ?, english_name = ?, alpha_2 = ?, alpha_3 = ?, iso = ?, location = ?, location_precise = ?, url = ? WHERE alpha_2 = ? OR alpha_3 = ?"
	_, err = transaction.Exec(query, country.Name, country.FullName, country.EnglishName, country.Alpha2, country.Alpha3, country.Iso, country.Location, country.LocationPrecise, country.Url, countryId, countryId)
	if err != nil {
		c.logger.Errorf("ChangeCountry: error while updating country:%s", err)
		return fmt.Errorf("ChangeCountry: error while updating country:%w", err)
	}
	return transaction.Commit()
}

func (c *CountryRepository) DeleteCountry(countryId string) error {
	transaction, err := c.db.Begin()
	if err != nil {
		c.logger.Errorf("DeleteCountry: can not starts transaction:%s", err)
		return fmt.Errorf("DeleteCountry: can not starts transaction:%w", err)
	}
	defer transaction.Rollback()
	query := "DELETE FROM countries WHERE alpha_2 = ? OR alpha_3 = ?"
	_, err = c.db.Exec(query, countryId, countryId)
	if err != nil {
		c.logger.Errorf("Error while scanning for countryId:%s", err)
		return fmt.Errorf("DeleteCountry: error while scanning for countryId:%s", err)
	}
	return transaction.Commit()
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
