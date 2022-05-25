package models

type Country struct {
	Name            string `json:"name"`
	FullName        string `json:"full_name"`
	EnglishName     string `json:"english_name"`
	Alpha2          string `json:"alpha_2"`
	Alpha3          string `json:"alpha_3"`
	Iso             int    `json:"iso"`
	Location        string `json:"location"`
	LocationPrecise string `json:"location_precise"`
	Url             string `json:"url"`
}

type ResponseCountry struct {
	Name            string `json:"name" valid:"required"`
	FullName        string `json:"full_name"`
	EnglishName     string `json:"english_name" valid:"required"`
	Alpha2          string `json:"alpha_2"  valid:"stringlength(2|2)"`
	Alpha3          string `json:"alpha_3"  valid:"stringlength(3|3)"`
	Iso             int    `json:"iso" valid:"required"`
	Location        string `json:"location" `
	LocationPrecise string `json:"location_precise"`
	Url             string `json:"url"`
}

type Filters struct {
	Page  uint64
	Limit uint64
	Flag  bool
}

type User struct {
	Name        string   `json:"name" valid:"alpha"`
	Email       string   `json:"email" valid:"email"`
	Description string   `json:"description"`
	CountryId   int      `json:"country_id"`
	Hobbies     []string `json:"hobbies"`
}

type ResponseUser struct {
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Description string   `json:"description"`
	CountryId   int      `json:"country_id"`
	Hobbies     []string `json:"hobbies"`
}

type Options struct {
	Page  uint64
	Limit uint64
}
