package service

import (
	"context"
	"sample/common/model"
	"sample/common/response"
	"sample/repository"

	"sample/common/log"

	"github.com/google/uuid"
)

type IPlaylistService interface {
	GetPlaylists(ctx context.Context, filter model.PlaylistFilter) (int, any)
	GetPlaylistById(ctx context.Context, playlistUuid string) (int, any)
	PostPlaylist(ctx context.Context, playlistRequest model.PlaylistRequest) (int, any)
	DeletePlaylistById(ctx context.Context, playlistUuid string) (int, any)
	PutPlaylistById(ctx context.Context, playlistUuid string, playlistRequest model.PlaylistRequest) (int, any)
}

type Playlist struct {
}

func NewPlaylist() IPlaylistService {
	return &Playlist{}
}

func (s *Playlist) GetPlaylists(ctx context.Context, filter model.PlaylistFilter) (int, any) {
	tracks, err := repository.PlaylistRepo.GetPlaylists(ctx, filter)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.OK(tracks)
}

func (s *Playlist) PostPlaylist(ctx context.Context, playlistRequest model.PlaylistRequest) (int, any) {
	// default value playback mode and update if playlist request set playback mode
	playbackMode := "priority"
	if len(playlistRequest.PlaybackMode) > 0 {
		playbackMode = playlistRequest.PlaybackMode
	}
	playlist := model.Playlist{
		ID:           uuid.NewString(),
		Name:         playlistRequest.Name,
		TrackIds:     playlistRequest.TrackIds,
		PlaybackMode: playbackMode,
	}
	err := repository.PlaylistRepo.PostPlaylist(ctx, playlist)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.OK(playlist)
}

func (s *Playlist) GetPlaylistById(ctx context.Context, playlistUuid string) (int, any) {
	playlist, err := repository.PlaylistRepo.GetPlaylistById(ctx, playlistUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.OK(playlist)
}

func (s *Playlist) DeletePlaylistById(ctx context.Context, playlistUuid string) (int, any) {
	// check exits playlist id
	if _, err := repository.PlaylistRepo.GetPlaylistById(ctx, playlistUuid); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	err := repository.PlaylistRepo.DeletePlaylistById(ctx, playlistUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.OK(map[string]interface{}{
		"delete success playlist id": playlistUuid,
	})
}

func (s *Playlist) PutPlaylistById(ctx context.Context, playlistUuid string, playlistRequest model.PlaylistRequest) (int, any) {
	// check exits playlist id
	if _, err := repository.PlaylistRepo.GetPlaylistById(ctx, playlistUuid); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	playlistUpdate := model.Playlist{
		ID:           playlistUuid,
		Name:         playlistRequest.Name,
		TrackIds:     playlistRequest.TrackIds,
		PlaybackMode: playlistRequest.PlaybackMode,
	}

	err := repository.PlaylistRepo.PutPlaylistById(ctx, playlistUuid, playlistUpdate)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(map[string]interface{}{
		"update success playlist id": playlistUuid,
	})
}
