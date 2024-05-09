package service

import (
	"context"
	"io"
	"mime/multipart"
	"sample/common/model"
	"sample/common/response"
	"sample/repository"
	"strings"

	"sample/common/log"

	"github.com/tcolgate/mp3"

	"github.com/google/uuid"
)

type ITrackService interface {
	GetTracks(ctx context.Context, filter model.TrackFilter) (int, any)
	GetTrackById(ctx context.Context, trackUuid string) (int, any)
	PostTrack(ctx context.Context, track model.TrackRequest, fileUpload *multipart.FileHeader) (int, any)
	DeleteTrackById(ctx context.Context, trackUuid string) (int, any)
	PutTrackById(ctx context.Context, trackUuid string, trackUpdate model.TrackRequest, fileUpload *multipart.FileHeader) (int, any)
}

type Track struct {
}

func NewTrack() ITrackService {
	return &Track{}
}

func (s *Track) GetTracks(ctx context.Context, filter model.TrackFilter) (int, any) {
	tracks, err := repository.TrackRepo.GetTracks(ctx, filter)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.OK(tracks)
}

func (s *Track) PostTrack(ctx context.Context, trackRequest model.TrackRequest, fileUpload *multipart.FileHeader) (int, any) {
	// Parse audio file duration

	fileExtension := fileUpload.Filename[strings.LastIndex(fileUpload.Filename, ".")+1:]
	duration, err := HandleParseAudioDuration(fileUpload)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	track := model.Track{
		ID:          uuid.NewString(),
		Title:       trackRequest.Title,
		Artist:      trackRequest.Artist,
		Album:       trackRequest.Album,
		Genre:       trackRequest.Genre,
		ReleaseYear: trackRequest.ReleaseYear,
		Duration:    duration,
		MP3File:     trackRequest.Title + "." + fileExtension,
	}
	err = repository.TrackRepo.PostTrack(ctx, track)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.OK(track)
}

func (s *Track) GetTrackById(ctx context.Context, trackUuid string) (int, any) {
	track, err := repository.TrackRepo.GetTrackById(ctx, trackUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.OK(track)
}

func (s *Track) DeleteTrackById(ctx context.Context, trackUuid string) (int, any) {
	// check exits track id
	if _, err := repository.TrackRepo.GetTrackById(ctx, trackUuid); err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	err := repository.TrackRepo.DeleteTrackById(ctx, trackUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}
	return response.OK(map[string]interface{}{
		"delete success track id": trackUuid,
	})
}

func (s *Track) PutTrackById(ctx context.Context, trackUuid string, trackRequest model.TrackRequest, fileUpload *multipart.FileHeader) (int, any) {
	// check exits track id
	trackExist, err := repository.TrackRepo.GetTrackById(ctx, trackUuid)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if trackExist == nil {
		log.Error(err)
		return response.BadRequestMsg("track not found")
	}

	if fileUpload != nil {
		duration, err := HandleParseAudioDuration(fileUpload)
		if err != nil {
			log.Error(err)
			return response.ServiceUnavailableMsg(err.Error())
		}
		trackExist.Duration = duration

		fileExtension := fileUpload.Filename[strings.LastIndex(fileUpload.Filename, ".")+1:]
		trackExist.MP3File = trackRequest.Title + "." + fileExtension
	}

	trackUpdate := model.Track{
		ID:          trackUuid,
		Title:       trackRequest.Title,
		Artist:      trackRequest.Artist,
		Album:       trackRequest.Album,
		Genre:       trackRequest.Genre,
		ReleaseYear: trackRequest.ReleaseYear,
		Duration:    trackExist.Duration,
		MP3File:     trackExist.MP3File,
	}

	err = repository.TrackRepo.PutTrackById(ctx, trackUuid, trackUpdate)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	}

	return response.OK(trackUpdate)
}

func HandleParseAudioDuration(file *multipart.FileHeader) (float64, error) {
	var duration float64
	fd, err := file.Open()
	if err != nil {
		log.Error(err)
		return 0, err
	}
	defer fd.Close()

	decode := mp3.NewDecoder(fd)
	var frame mp3.Frame
	skipped := 0
	for {
		if err := decode.Decode(&frame, &skipped); err != nil {
			if err == io.EOF {
				break
			}
			log.Error(err)
			return 0, err
		}
		duration = duration + frame.Duration().Seconds()
	}

	return duration, nil
}
