package repository

import (
	"context"
	"sample/common/model"
)

type IPlaylist interface {
	GetPlaylists(ctx context.Context, filter model.PlaylistFilter) (*[]model.Playlist, error)
	GetPlaylistById(ctx context.Context, playlistUuid string) (*model.Playlist, error)
	PostPlaylist(ctx context.Context, playlist model.Playlist) error
	DeletePlaylistById(ctx context.Context, playlistUuid string) error
	PutPlaylistById(ctx context.Context, playlistUuid string, playlistUpdate model.Playlist) error
}

var PlaylistRepo IPlaylist
