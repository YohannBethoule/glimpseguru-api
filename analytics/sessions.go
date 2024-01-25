package analytics

import (
	"glimpseguru-api/authent"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s Service) getSessionMatchFilter(website authent.Website) bson.D {
	return bson.D{
		{"websiteID", website.ID},
		{"start", bson.D{{"$gte", s.start}, {"$lte", s.end}}},
	}
}

func (s Service) getSessionCount(website authent.Website) (int, error) {
	filter := s.getSessionMatchFilter(website)
	sessionCount, err := s.sessionsCollection.CountDocuments(s.Context, filter)
	if err != nil {
		return 0, err
	}
	return int(sessionCount), nil
}

func (s Service) getSessionTime(website authent.Website) (SessionTime, error) {
	var sessionTime SessionTime

	matchStage := bson.D{{"$match", s.getSessionMatchFilter(website)}}

	// Group stage
	groupStage := bson.D{{"$group", bson.D{
		{"_id", nil},
		{"minSessionTime", bson.D{{"$min", bson.D{{"$subtract", bson.A{"$end", "$start"}}}}}},
		{"maxSessionTime", bson.D{{"$max", bson.D{{"$subtract", bson.A{"$end", "$start"}}}}}},
		{"avgSessionTime", bson.D{{"$avg", bson.D{{"$subtract", bson.A{"$end", "$start"}}}}}},
	}}}

	// Aggregate pipeline
	cursor, err := s.sessionsCollection.Aggregate(s.Context, mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		return sessionTime, err
	}

	// Fetching results
	if cursor.Next(s.Context) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return sessionTime, err
		}

		// Convert results to milliseconds
		sessionTime.Min = convertMilliSecondsToMinutes(result["minSessionTime"])
		sessionTime.Max = convertMilliSecondsToMinutes(result["maxSessionTime"])
		sessionTime.Mean = convertMilliSecondsToMinutes(result["avgSessionTime"])
	}

	return sessionTime, nil
}

func convertMilliSecondsToMinutes(val interface{}) int {
	if val == nil {
		return 0
	}

	var millis int64

	switch v := val.(type) {
	case int64:
		millis = v
	case float64:
		millis = int64(v)
	default:
		return 0 // Unknown type
	}

	// Convert milliseconds to minutes
	return int(millis / 60000)
}
