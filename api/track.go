package api

import (
	"fmt"
	"net/http"
	"os"
	"sample/common/log"
	"sample/common/model"
	"sample/common/response"
	"sample/common/util"
	"sample/service"

	"github.com/gin-gonic/gin"
)

type Track struct {
	trackService service.ITrackService
}

func APIMusicTrackHandler(r *gin.Engine, trackService service.ITrackService) {
	handler := &Track{
		trackService: trackService,
	}
	Group := r.Group("v1/track")
	{
		Group.GET("", handler.GetTracks)
		Group.GET(":id", handler.GetTrackById)
		Group.POST("", handler.PostTrack)
		Group.PUT(":id", handler.PutTrackById)
		Group.DELETE(":id", handler.DeleteTrackById)
		Group.GET(":id/download", handler.DownloadTrackById)
	}
}

// GetTracks godoc
// @Summary Get tracks
// @Description Get tracks
// @Tags track
// @Id get-track
// @Accept json
// @Produce json
// @Param title query string false "title"
// @Param artist query string false "artist"
// @Param album query string false "album"
// @Param genre query string false "genre"
// @Success 200 {object} model.Track
// @Router /track [get]
func (m *Track) GetTracks(c *gin.Context) {
	filter := model.TrackFilter{
		Title:  c.Query("title"),
		Artist: c.Query("artist"),
		Album:  c.Query("album"),
		Genre:  c.Query("genre"),
		Limit:  util.ParseInt(c.Query("limit")),
		Offset: util.ParseInt(c.Query("offset")),
	}
	code, result := m.trackService.GetTracks(c, filter)
	c.JSON(code, result)
}

// PostTracks godoc
// @Summary Post tracks
// @Description Post tracks
// @Tags track
// @Id post-track
// @Accept json
// @Produce json
// @Param track body model.TrackRequest true "track"
// @Param mp3_file formData file true "mp3_file"
// @Success 200 {object} model.Track
// @Router /track [post]
func (m *Track) PostTrack(c *gin.Context) {
	file, err := c.FormFile("mp3_file")
	if err != nil {
		if err == http.ErrMissingFile {
			c.JSON(response.BadRequestMsg("mp3_file is required"))
			return
		} else {
			log.Error(err)
			c.JSON(response.ServiceUnavailableMsg(err.Error()))
			return
		}
	}

	// Check file upload valid audio file
	fileType := file.Header.Get("Content-type")
	if fileType != "audio/mp3" && fileType != "audio/mpeg" && fileType != "audio/wav" {
		c.JSON(response.BadRequestMsg("Invalid file format. Please upload an MP3 audio file."))
		return
	}

	trackRequest := model.TrackRequest{
		Title:       c.PostForm("title"),
		Artist:      c.PostForm("artist"),
		Album:       c.PostForm("album"),
		Genre:       c.PostForm("genre"),
		ReleaseYear: util.ParseInt(c.PostForm("release_year")),
	}
	if err := trackRequest.Validate(); err != nil {
		code, _ := response.BadRequest()
		c.JSON(code, err)
		return
	}
	code, result := m.trackService.PostTrack(c, trackRequest, file)
	if code == http.StatusOK {
		// if create success, upload file to audio dir
		track := result.(model.Track)
		audioDir := util.GetAudioDir() + "/" + track.ID
		if _, err := os.Stat(audioDir); os.IsNotExist(err) {
			os.MkdirAll(audioDir, 0755)
		}
		if err := c.SaveUploadedFile(file, audioDir+"/"+track.MP3File); err != nil {
			log.Error(err)
			c.JSON(response.ServiceUnavailableMsg(err.Error()))
			return
		}
	}
	c.JSON(code, result)
}

// GetTrackById godoc
// @Summary Get track by id
// @Description Get track by id
// @Tags track
// @Id get-track-id
// @Accept json
// @Produce json
// @Param id path string true "Track ID"
// @Success 200 {object} model.Track
// @Router /track/{id} [get]
func (m *Track) GetTrackById(c *gin.Context) {
	trackUuid := c.Param("id")
	if trackUuid == "" {
		c.JSON(response.BadRequestMsg("id is missing"))
		c.Abort()
		return
	}
	code, result := m.trackService.GetTrackById(c, trackUuid)
	c.JSON(code, result)
}

// DeleteTrackById godoc
// @Summary Delete track by id
// @Description Delete track by id
// @Tags track
// @Id Delete-track-id
// @Accept json
// @Produce json
// @Param id path string true "Track ID"
// @Success 200 {object} model.Track
// @Router /track/{id} [delete]
func (m *Track) DeleteTrackById(c *gin.Context) {
	trackUuid := c.Param("id")
	if trackUuid == "" {
		c.JSON(response.BadRequestMsg("id is missing"))
		c.Abort()
		return
	}
	code, result := m.trackService.DeleteTrackById(c, trackUuid)
	c.JSON(code, result)
}

// PutTrackById godoc
// @Summary Put track by id
// @Description Put track by id
// @Tags track
// @Id put-track
// @Accept json
// @Produce json
// @Param id path string true "Track ID"
// @Param track body model.TrackRequest true "track"
// @Param mp3_file formData file false "mp3_file"
// @Success 200 {object} model.Track
// @Router /track/{id} [Put]
func (m *Track) PutTrackById(c *gin.Context) {
	trackUuid := c.Param("id")
	file, err := c.FormFile("mp3_file")
	if err != nil && err != http.ErrMissingFile {
		log.Error(err)
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}

	// Check file upload valid audio file
	if file != nil {
		fileType := file.Header.Get("Content-type")
		if fileType != "audio/mp3" && fileType != "audio/mpeg" && fileType != "audio/wav" {
			c.JSON(response.BadRequestMsg("Invalid file format. Please upload an MP3 audio file."))
			return
		}
	}

	trackPost := model.TrackRequest{
		Title:       c.PostForm("title"),
		Artist:      c.PostForm("artist"),
		Album:       c.PostForm("album"),
		Genre:       c.PostForm("genre"),
		ReleaseYear: util.ParseInt(c.PostForm("release_year")),
	}
	if err := trackPost.Validate(); err != nil {
		code, _ := response.BadRequest()
		c.JSON(code, err)
		return
	}

	code, result := m.trackService.PutTrackById(c, trackUuid, trackPost, file)
	if code == http.StatusOK && file != nil {
		// if update success, upload file to audio dir
		track := result.(model.Track)
		audioDir := util.GetAudioDir() + "/" + track.ID
		if _, err := os.Stat(audioDir); os.IsNotExist(err) {
			os.MkdirAll(audioDir, 0755)
		}
		if err := c.SaveUploadedFile(file, audioDir+"/"+track.MP3File); err != nil {
			log.Error(err)
			c.JSON(response.ServiceUnavailableMsg(err.Error()))
			return
		}
	}
	c.JSON(code, result)
}

// DownloadTrackById godoc
// @Summary download track by id
// @Description download track by id
// @Tags track
// @Id download-track-id
// @Accept json
// @Produce json
// @Param id path string true "Track ID"
// @Success 200 {object} model.Track
// @Router /track/{id}/download [get]
func (m *Track) DownloadTrackById(c *gin.Context) {
	trackUuid := c.Param("id")
	if trackUuid == "" {
		c.JSON(response.BadRequestMsg("id is missing"))
		c.Abort()
		return
	}
	code, result := m.trackService.GetTrackById(c, trackUuid)
	if code != http.StatusOK {
		c.JSON(code, result)
		return
	}
	track := result.(*model.Track)
	audioDir := util.GetAudioDir() + "/" + trackUuid + "/" + track.MP3File
	audioByte, err := os.ReadFile(audioDir)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	f, err := os.Open(audioDir)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	contentType := http.DetectContentType(audioByte[:512])
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%q", track.MP3File))
	c.Writer.Header().Add("Content-Type", contentType)
	c.Writer.Header().Add("Cache-Control", "no-cache, must-revalidate")
	c.Writer.Header().Add("Accept-Ranges", "bytes")
	c.Writer.Header().Add("Content-Length", fmt.Sprintf("%d", fi.Size()))
	c.Data(http.StatusOK, contentType, audioByte)
	return
}
