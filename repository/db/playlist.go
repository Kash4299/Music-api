package db

import (
	"context"
	"sample/common/model"
	"sample/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var playlistCollection *mongo.Collection

type Playlist struct {
}

func NewPlaylist(client *mongo.Client) repository.IPlaylist {
	playlistCollection = client.Database("music").Collection("playlist")
	return &Playlist{}
}

func (repo *Playlist) GetPlaylists(ctx context.Context, filter model.PlaylistFilter) (*[]model.Playlist, error) {
	playlists := new([]model.Playlist)
	query := bson.D{}
	findOptions := options.Find()

	if len(filter.Name) > 0 {
		query = append(query, bson.E{Key: "name", Value: filter.Name})
	}

	if filter.Limit > 0 {
		findOptions.SetLimit(int64(filter.Limit))
		findOptions.SetSkip(int64(filter.Offset))
	}

	cursor, err := playlistCollection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, playlists)
	if err != nil {
		return nil, err
	}

	return playlists, nil
}

func (repo *Playlist) PostPlaylist(ctx context.Context, playlist model.Playlist) error {
	_, err := playlistCollection.InsertOne(ctx, playlist)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Playlist) GetPlaylistById(ctx context.Context, playlistUuid string) (*model.Playlist, error) {
	playlist := new(model.Playlist)
	err := playlistCollection.FindOne(ctx, bson.M{"_id": playlistUuid}).Decode(playlist)
	if err != nil {
		return nil, err
	}
	return playlist, nil
}

func (repo *Playlist) DeletePlaylistById(ctx context.Context, playlistUuid string) error {
	_, err := playlistCollection.DeleteOne(ctx, bson.M{"_id": playlistUuid})
	if err != nil {
		return err
	}
	return nil
}

func (repo *Playlist) PutPlaylistById(ctx context.Context, playlistUuid string, playlistUpdate model.Playlist) error {
	_, err := playlistCollection.UpdateOne(ctx, bson.M{"_id": playlistUuid}, bson.M{"$set": playlistUpdate})
	if err != nil {
		return err
	}
	return nil
}
