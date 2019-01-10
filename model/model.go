package model

// Data acts as a container for generic data schemas that will need to be downloaded

type Data struct {
	id    string
	Name  string
	Data  []interface{}
	Found bool
}

//FoodData holds attributes of a single query to fda food data api

type FoodData struct {
	Meta struct {
		Disclaimer  string `json:"disclaimer"`
		Terms       string `json:"terms"`
		License     string `json:"license"`
		LastUpdated string `json:"last_updated"`
		Results     struct {
			Skip  int `json:"skip"`
			Limit int `json:"limit"`
			Total int `json:"total"`
		} `json:"results"`
	} `json:"meta"`
	Results []struct {
		ReportNumber string        `json:"report_number"`
		Outcomes     []interface{} `json:"outcomes"`
		DateCreated  string        `json:"date_created"`
		Reactions    []string      `json:"reactions"`
		DateStarted  interface{}   `json:"date_started"`
		Consumer     struct{}      `json:"consumer"`
		Products     []struct {
			IndustryName string `json:"industry_name"`
			Role         string `json:"role"`
			IndustryCode string `json:"industry_code"`
			NameBrand    string `json:"name_brand"`
		} `json:"products"`
	} `json:"results"`
}
