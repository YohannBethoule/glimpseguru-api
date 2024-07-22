package analytics

import (
	"context"
	"glimpseguru-api/authent"
	"glimpseguru-api/db"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Service struct {
	pageViewsCollection    *mongo.Collection
	sessionsCollection     *mongo.Collection
	customEventsCollection *mongo.Collection
	context.Context
	start time.Time
	end   time.Time
}

func NewService(start int64, end int64) *Service {
	return &Service{
		Context:                context.Background(),
		start:                  time.Unix(start, 0),
		end:                    time.Unix(end, 0),
		pageViewsCollection:    db.MongoDb.Collection("page_views"),
		sessionsCollection:     db.MongoDb.Collection("sessions"),
		customEventsCollection: db.MongoDb.Collection("custom_events"),
	}
}

func (s Service) GetUsersAnalytics(user authent.User) ([]AnalyticsData, []error) {
	var userAnalytics []AnalyticsData
	var errors []error

	for _, website := range user.Websites {
		websiteAnalytics, errsWebsiteAnalytics := s.getWebsiteAnalytics(website)
		if errsWebsiteAnalytics != nil {
			errors = append(errors, errsWebsiteAnalytics...)
		}
		userAnalytics = append(userAnalytics, websiteAnalytics)
	}
	return userAnalytics, errors
}

func (s Service) getWebsiteAnalytics(website authent.Website) (AnalyticsData, []error) {
	var errors []error
	websiteAnalytics := AnalyticsData{
		Website: website,
	}

	viewCount, errViewCount := s.getViewCount(website)
	if errViewCount != nil {
		errors = append(errors, errViewCount)
	} else {
		websiteAnalytics.ViewCount = viewCount
	}

	sessionCount, errSessionCount := s.getSessionCount(website)
	if errSessionCount != nil {
		errors = append(errors, errSessionCount)
	} else {
		websiteAnalytics.SessionCount = sessionCount
	}

	sessionTime, errSessionTime := s.getSessionTime(website)
	if errSessionTime != nil {
		errors = append(errors, errSessionTime)
	} else {
		websiteAnalytics.SessionTime = sessionTime
	}

	if deviceUsage, errDeviceUsage := s.getDeviceUsage(website); errDeviceUsage != nil {
		errors = append(errors, errDeviceUsage)
	} else {
		websiteAnalytics.DeviceType = deviceUsage
	}

	if sourceType, errSourceType := s.getSourceType(website); errSourceType != nil {
		errors = append(errors, errSourceType)
	} else {
		websiteAnalytics.SourceType = sourceType
	}

	if pageViews, errPageViews := s.getPageViews(website); errPageViews != nil {
		errors = append(errors, errPageViews)
	} else {
		websiteAnalytics.PagesViews = pageViews
	}

	if customEvents, errCustomEvents := s.getCustomEvents(website); errCustomEvents != nil {
		errors = append(errors, errCustomEvents)
	} else {
		websiteAnalytics.CustomEvents = customEvents
	}

	if viewsHistogram, errViewsHistogram := s.getViewsHistogram(website); errViewsHistogram != nil {
		errors = append(errors, errViewsHistogram)
	} else {
		websiteAnalytics.ViewsHistogram = viewsHistogram
	}

	return websiteAnalytics, errors
}
