package analytics

import (
	"glimpseguru-api/authent"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func (s Service) getPageViewMatchFilter(website authent.Website) bson.D {
	return bson.D{
		{"websiteID", website.ID},
		{"timestamp", bson.D{{"$gte", s.start}, {"$lte", s.end}}},
	}
}

func (s Service) getViewCount(website authent.Website) (int, error) {
	filter := s.getPageViewMatchFilter(website)
	viewCount, err := s.pageViewsCollection.CountDocuments(s.Context, filter)
	if err != nil {
		return 0, err
	}
	return int(viewCount), nil
}

func (s Service) getDeviceUsage(website authent.Website) (map[string]float64, error) {
	deviceUsage := make(map[string]float64)

	matchStage := bson.D{{"$match", s.getPageViewMatchFilter(website)}}

	// Group by device type and count occurrences
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$deviceType"},
			{"count", bson.D{{"$sum", 1}}},
		}},
	}

	// Perform the aggregation
	cursor, err := s.pageViewsCollection.Aggregate(s.Context, mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		return nil, err
	}

	// Iterate over the cursor and populate the result map
	for cursor.Next(s.Context) {
		var result struct {
			ID    string `bson:"_id"`
			Count int    `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		deviceUsage[result.ID] = float64(result.Count)
	}

	// Calculate total views to determine percentages
	totalViews := 0
	for _, count := range deviceUsage {
		totalViews += int(count)
	}

	// Convert counts to percentages
	for deviceType, count := range deviceUsage {
		deviceUsage[deviceType] = (count / float64(totalViews)) * 100
	}

	return deviceUsage, nil
}

func (s Service) getSourceType(website authent.Website) (map[string]float64, error) {
	sourceTypeUsage := make(map[string]float64)

	matchStage := bson.D{{"$match", s.getPageViewMatchFilter(website)}}

	// Group by source type and count occurrences
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$sourceType"},
			{"count", bson.D{{"$sum", 1}}},
		}},
	}

	// Perform the aggregation
	cursor, err := s.pageViewsCollection.Aggregate(s.Context, mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		return nil, err
	}

	// Iterate over the cursor and populate the result map
	for cursor.Next(s.Context) {
		var result struct {
			ID    string `bson:"_id"`
			Count int    `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		sourceTypeUsage[result.ID] = float64(result.Count)
	}

	// Calculate total views to determine percentages
	totalViews := 0
	for _, count := range sourceTypeUsage {
		totalViews += int(count)
	}

	// Convert counts to percentages
	for sourceType, count := range sourceTypeUsage {
		sourceTypeUsage[sourceType] = (count / float64(totalViews)) * 100
	}

	return sourceTypeUsage, nil
}

func (s Service) getPageViews(website authent.Website) (map[string]int, error) {
	pageViews := make(map[string]int)

	matchStage := bson.D{{"$match", s.getPageViewMatchFilter(website)}}

	// Group by source type and count occurrences
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$pathname"},
			{"count", bson.D{{"$sum", 1}}},
		}},
	}

	// Perform the aggregation
	cursor, err := s.pageViewsCollection.Aggregate(s.Context, mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		return nil, err
	}

	// Iterate over the cursor and populate the result map
	for cursor.Next(s.Context) {
		var result struct {
			ID    string `bson:"_id"`
			Count int    `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		pageViews[result.ID] = result.Count
	}

	return pageViews, nil
}

func (s Service) getViewsHistogram(website authent.Website) (map[string]int, error) {
	viewsHistogram := make(map[string]int)
	timeFormat := "2006-01-02"
	mongoTimeFormat := "%Y-%m-%d"
	timeIncrement := 24 * time.Hour

	if s.end.Sub(s.start) < 48*time.Hour {
		timeFormat = "2006-01-02 15:00"
		mongoTimeFormat = "%Y-%m-%d %H:00"
		timeIncrement = time.Hour
	}

	// Generate all time slots from start to end
	for t := s.start; !t.After(s.end); t = t.Add(timeIncrement) {
		viewsHistogram[t.Format(timeFormat)] = 0
	}

	matchStage := bson.D{{"$match", s.getPageViewMatchFilter(website)}}

	// Group by date or hour and count occurrences
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", bson.D{
				{"$dateToString", bson.D{
					{"format", mongoTimeFormat},
					{"date", "$timestamp"},
				}},
			}},
			{"count", bson.D{{"$sum", 1}}},
		}},
	}

	// Perform the aggregation
	cursor, err := s.pageViewsCollection.Aggregate(s.Context, mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		return nil, err
	}

	// Iterate over the cursor and populate the result map
	for cursor.Next(s.Context) {
		var result struct {
			ID    string `bson:"_id"`
			Count int    `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		viewsHistogram[result.ID] = result.Count
	}

	return viewsHistogram, nil
}
