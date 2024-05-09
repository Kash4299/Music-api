package db

import (
	"context"
	"sample/common/model"
	"sample/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var trackCollection *mongo.Collection

type Track struct {
}

func NewTrack(client *mongo.Client) repository.ITracks {
	trackCollection = client.Database("music").Collection("tracks")
	return &Track{}
}

func (repo *Track) GetTracks(ctx context.Context, filter model.TrackFilter) (*[]model.Track, error) {
	tracks := new([]model.Track)
	query := bson.D{}
	findOptions := options.Find()

	if len(filter.Title) > 0 {
		query = append(query, bson.E{Key: "title", Value: filter.Title})
	}
	if len(filter.Artist) > 0 {
		query = append(query, bson.E{Key: "artist", Value: filter.Artist})
	}
	if len(filter.Album) > 0 {
		query = append(query, bson.E{Key: "album", Value: filter.Album})
	}
	if len(filter.Genre) > 0 {
		query = append(query, bson.E{Key: "genre", Value: filter.Genre})
	}

	if filter.Limit > 0 {
		findOptions.SetLimit(int64(filter.Limit))
		findOptions.SetSkip(int64(filter.Offset))
	}

	cursor, err := trackCollection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, tracks)
	if err != nil {
		return nil, err
	}

	return tracks, nil
}

func (repo *Track) PostTrack(ctx context.Context, track model.Track) error {
	_, err := trackCollection.InsertOne(ctx, track)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Track) GetTrackById(ctx context.Context, trackUuid string) (*model.Track, error) {
	track := new(model.Track)
	err := trackCollection.FindOne(ctx, bson.M{"_id": trackUuid}).Decode(track)
	if err != nil {
		return nil, err
	}
	return track, nil
}

func (repo *Track) DeleteTrackById(ctx context.Context, trackUuid string) error {
	_, err := trackCollection.DeleteOne(ctx, bson.M{"_id": trackUuid})
	if err != nil {
		return err
	}
	return nil
}

func (repo *Track) PutTrackById(ctx context.Context, trackUuid string, trackUpdate model.Track) error {
	_, err := trackCollection.UpdateOne(ctx, bson.M{"_id": trackUuid}, bson.M{"$set": trackUpdate})
	if err != nil {
		return err
	}
	return nil
}
