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
	Name            string `json:"name" validate:"required"`
	FullName        string `json:"full_name"`
	EnglishName     string `json:"english_name" validate:"required"`
	Alpha2          string `json:"alpha_2"  validate:"alpha_2"`
	Alpha3          string `json:"alpha_3"  validate:"alpha_3"`
	Iso             int    `json:"iso" validate:"required"`
	Location        string `json:"location" `
	LocationPrecise string `json:"location_precise"`
	Url             string `json:"url"`
}
