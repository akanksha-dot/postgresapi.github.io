package models

type Country struct {
	CountryID int64  `json:"countryid"`
	Name      string `json:"name"`
	Type      string `json:"type"`
}
