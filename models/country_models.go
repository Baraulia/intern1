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
}
