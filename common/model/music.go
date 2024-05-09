package model

import (
	"gopkg.in/validator.v2"
)

type Track struct {
	ID          string  `json:"id" bson:"_id,omitempty"`
	Title       string  `json:"title" bson:"title"`
	Artist      string  `json:"artist" bson:"artist"`
	Album       string  `json:"album" bson:"album"`
	Genre       string  `json:"genre" bson:"genre"`
	ReleaseYear int     `json:"release_year" bson:"release_year"`
	Duration    float64 `json:"duration" bson:"duration"`
	MP3File     string  `json:"mp3_file" bson:"mp3_file"`
}

type TrackRequest struct {
	Title       string  `json:"title" bson:"title" validate:"nonzero"`
	Artist      string  `json:"artist" bson:"artist"`
	Album       string  `json:"album" bson:"album"`
	Genre       string  `json:"genre" bson:"genre"`
	ReleaseYear int     `json:"release_year" bson:"release_year"`
	Duration    float64 `json:"duration" bson:"duration"`
}

func (track *TrackRequest) Validate() error {
	if errs := validator.Validate(track); errs != nil {
		return errs
	}
	return nil
}

type Playlist struct {
	ID           string     `json:"id" bson:"_id,omitempty"`
	Name         string     `json:"name" bson:"name"`
	TrackIds     []TrackIds `json:"track_ids" bson:"track_ids"`
	PlaybackMode string     `json:"playback_mode" bson:"playback_mode"` // 'priority' or 'random'
}

type PlaylistRequest struct {
	Name         string     `json:"name" bson:"name" validate:"nonzero"`
	TrackIds     []TrackIds `json:"track_ids" bson:"track_ids"`
	PlaybackMode string     `json:"playback_mode" bson:"playback_mode"`
}

func (playlist *PlaylistRequest) Validate() error {
	if errs := validator.Validate(playlist); errs != nil {
		return errs
	}
	return nil
}

type TrackIds struct {
	TrackID  string `json:"track_id" bson:"track_id"`
	Priority int    `json:"priority" bson:"priority"`
}

type TrackFilter struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
	Genre  string `json:"genre"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

type PlaylistFilter struct {
	Name   string `json:"name"  bson:"name"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}
