package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	foildb                = "foildb"
	contestantsCollection = "contestants"
	queueCollection       = "queue"
	leaderboardCollection = "leaderboard"
)

type Contestant struct {
	Id    primitive.ObjectID `bson:"_id"`
	Email string             `bson:"email"`
	Name  string             `bson:"name"`
}

type QueueItem struct {
	Id           primitive.ObjectID `bson:"_id"`
	Timestamp    int64              `bson:"timestamp"`
	ContestantId primitive.ObjectID `bson:"contestantId"`
}

type LeaderboardEntry struct {
	Id        primitive.ObjectID `bson:"_id"`
	Email     string             `bson:"emil"`
	SplitTime float64            `bson:"splitTime"`
	EndTime   float64            `bson:"endTime"`
}

type ContestantDTO struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type QueueItemDTO struct {
	Timestamp  int64         `json:"timestamp"`
	Contestant ContestantDTO `json:"contestant"`
}

type GetContestantsType int

type GameResultDTO struct {
	SplitTime float64 `json:"splitTime"`
	EndTime   float64 `json:"endTime"`
}

type LeaderboardEntryDTO struct {
	Contestant ContestantDTO `json:"contestant"`
	GameResult GameResultDTO `json:"result"`
}

const (
	GetContestantsAll GetContestantsType = iota
	GetContestantsEnqueued
	GetContestantsNotEnqueued
)

func ConnectDB(url string) (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func CloseDB(client *mongo.Client) {
	client.Disconnect(context.TODO())
}

func AddContestant(client *mongo.Client, email string, name string) (*Contestant, error) {
	collection := client.Database(foildb).Collection(contestantsCollection)
	filter := bson.D{{"email", email}}
	var findResult Contestant
	err := collection.FindOne(context.TODO(), filter).Decode(&findResult)
	if err == nil {
		return &findResult, nil
	}

	contestant := Contestant{Id: primitive.NewObjectID(), Email: email, Name: name}
	result, err := collection.InsertOne(context.TODO(), contestant)
	if err != nil {
		return nil, err
	}
	println(result.InsertedID)

	return &contestant, nil
}

func GetContestant(client *mongo.Client, contestantId primitive.ObjectID) (*Contestant, error) {
	collection := client.Database(foildb).Collection(contestantsCollection)
	filter := bson.D{{Key: "_id", Value: contestantId}}
	var contestant Contestant
	err := collection.FindOne(context.TODO(), filter).Decode(&contestant)
	if err != nil {
		return nil, err
	}
	return &contestant, nil
}

func GetContestants(client *mongo.Client, getType GetContestantsType) ([]ContestantDTO, error) {
	collection := client.Database(foildb).Collection(contestantsCollection)
	result, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	var contestants []ContestantDTO
	for result.Next(context.TODO()) {
		var contestant Contestant
		err := result.Decode(&contestant)
		if err != nil {
			return nil, err
		}

		getContestant := false
		if getType != GetContestantsAll {
			qcollection := client.Database(foildb).Collection(queueCollection)
			filter := bson.D{{Key: "contestantId", Value: contestant.Id}}
			qresult, err := qcollection.Find(context.TODO(), filter)
			if err != nil {
				return nil, err
			}
			if qresult.Next(context.TODO()) {
				if getType == GetContestantsEnqueued {
					getContestant = true
				}
			} else {
				if getType == GetContestantsNotEnqueued {
					getContestant = true
				}
			}
		} else {
			getContestant = true
		}
		if getContestant {
			contestants = append(contestants, ContestantDTO{Email: contestant.Email, Name: contestant.Name})
		}
	}

	return contestants, nil
}

func EnqueueContestant(client *mongo.Client, contestantId primitive.ObjectID) (*QueueItem, error) {
	collection := client.Database(foildb).Collection(queueCollection)

	queueItem := QueueItem{Id: primitive.NewObjectID(), Timestamp: time.Now().UnixNano(), ContestantId: contestantId}
	result, err := collection.InsertOne(context.TODO(), queueItem)
	if err != nil {
		return nil, err
	}
	println(result.InsertedID)

	return &queueItem, nil
}

func DeleteQueueItem(client *mongo.Client, timestamp int64) error {
	collection := client.Database(foildb).Collection(queueCollection)

	filter := bson.D{{Key: "timestamp", Value: timestamp}}
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}

func GetQueue(client *mongo.Client) ([]QueueItemDTO, error) {
	collection := client.Database(foildb).Collection(queueCollection)

	queryOptions := options.Find()
	queryOptions.SetSort(bson.D{{Key: "timestamp", Value: 1}})
	result, err := collection.Find(context.TODO(), bson.D{}, queryOptions)
	if err != nil {
		return nil, err
	}

	var queueItems []QueueItemDTO
	for result.Next(context.TODO()) {
		var queueItem QueueItem
		err := result.Decode(&queueItem)
		if err != nil {
			return nil, err
		}

		contestant, err := GetContestant(client, queueItem.ContestantId)
		if err != nil {
			return nil, err
		}

		queueItems = append(queueItems, QueueItemDTO{queueItem.Timestamp, ContestantDTO{Email: contestant.Email, Name: contestant.Name}})
	}

	return queueItems, nil
}

func CreateOrUpdateLeaderboardEntry(client *mongo.Client, currentGame QueueItemDTO, gameResult GameResultDTO) (*LeaderboardEntry, error) {
	collection := client.Database(foildb).Collection(leaderboardCollection)

	filter := bson.D{{Key: "email", Value: currentGame.Contestant.Email}}
	var entry LeaderboardEntry
	result := collection.FindOne(context.TODO(), filter).Decode(&entry)
	if result == nil {
		filter = bson.D{{Key: "_id", Value: entry.Id}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "EndTime", Value: gameResult.EndTime}, {Key: "SplitTime", Value: gameResult.SplitTime}}}}
		_, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return nil, err
		}
		return &entry, nil
	} else {
		entry = LeaderboardEntry{primitive.NewObjectID(), currentGame.Contestant.Email, gameResult.EndTime, gameResult.SplitTime}
		_, err := collection.InsertOne(context.TODO(), entry)
		if err != nil {
			return nil, err
		}
		return &entry, nil
	}
}

func GetLeaderboard(client *mongo.Client) ([]LeaderboardEntryDTO, error) {
	collection := client.Database(foildb).Collection(leaderboardCollection)
	queryOptions := options.Find()
	queryOptions.SetSort(bson.D{{Key: "endTime", Value: 1}})
	result, err := collection.Find(context.TODO(), bson.D{}, queryOptions)
	if err != nil {
		return nil, err
	}

	var leaderboardEntries []LeaderboardEntryDTO = make([]LeaderboardEntryDTO, 0)
	for result.Next(context.TODO()) {
		var leaderboardEntry LeaderboardEntry
		err := result.Decode(&leaderboardEntry)
		if err != nil {
			return nil, err
		}
		leaderboardEntries = append(leaderboardEntries, LeaderboardEntryDTO{ContestantDTO{Email: leaderboardEntry.Email, Name: ""}, GameResultDTO{SplitTime: leaderboardEntry.SplitTime, EndTime: leaderboardEntry.EndTime}})
	}
	return leaderboardEntries, nil
}
