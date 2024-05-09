package repository

import (
	"context"
	"sample/common/model"
)

type ITracks interface {
	GetTracks(ctx context.Context, filter model.TrackFilter) (*[]model.Track, error)
	GetTrackById(ctx context.Context, trackUuid string) (*model.Track, error)
	PostTrack(ctx context.Context, track model.Track) error
	DeleteTrackById(ctx context.Context, trackUuid string) error
	PutTrackById(ctx context.Context, trackUuid string, trackUpdate model.Track) error
}

var TrackRepo ITracks
