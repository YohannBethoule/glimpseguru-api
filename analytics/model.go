package analytics

import "glimpseguru-api/authent"

type AnalyticsData struct {
	Website        authent.Website    `json:"website"`
	ViewCount      int                `json:"viewCount"`
	SessionCount   int                `json:"sessionCount"`
	SessionTime    SessionTime        `json:"sessionTime"`
	DeviceType     map[string]float64 `json:"deviceType"`
	SourceType     map[string]float64 `json:"sourceType"`
	PagesViews     map[string]int     `json:"pagesViews"`
	CustomEvents   map[string]int     `json:"customEvents"`
	ViewsHistogram map[string]int     `json:"viewsHistogram"`
}

type SessionTime struct {
	Min  int `json:"min"`
	Max  int `json:"max"`
	Mean int `json:"mean"`
}
