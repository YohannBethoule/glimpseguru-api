package analytics

import (
	"glimpseguru-api/authent"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s Service) getCustomEvents(website authent.Website) (map[string]int, error) {
	customEvents := make(map[string]int)

	matchStage := bson.D{{"$match", s.getPageViewMatchFilter(website)}}

	// Group by source type and count occurrences
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$label"},
			{"count", bson.D{{"$sum", 1}}},
		}},
	}

	// Perform the aggregation
	cursor, err := s.customEventsCollection.Aggregate(s.Context, mongo.Pipeline{matchStage, groupStage})
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
		customEvents[result.ID] = result.Count
	}

	return customEvents, nil
}
