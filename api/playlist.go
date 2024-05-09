package api

import (
	"sample/common/model"
	"sample/common/response"
	"sample/common/util"
	"sample/service"

	"github.com/gin-gonic/gin"
)

type Playlist struct {
	playListService service.IPlaylistService
}

func APIPlaylistHandler(r *gin.Engine, playListService service.IPlaylistService) {
	handler := &Playlist{
		playListService: playListService,
	}
	Group := r.Group("v1/playlist")
	{
		Group.GET("", handler.GetPlaylists)
		Group.GET(":id", handler.GetPlaylistById)
		Group.POST("", handler.PostPlaylist)
		Group.DELETE(":id", handler.DeletePlaylistById)
		Group.PUT(":id", handler.PutPlaylistById)
	}
}

// GetPlaylist godoc
// @Summary Get playlists
// @Description Get playlists
// @Tags playlist
// @Id get-playlist
// @Accept json
// @Produce json
// @Param name query string false "name"
// @Success 200 {object} model.Playlist
// @Router /playlist [get]
func (p *Playlist) GetPlaylists(c *gin.Context) {
	filter := model.PlaylistFilter{
		Name:   c.Query("name"),
		Limit:  util.ParseInt(c.Query("limit")),
		Offset: util.ParseInt(c.Query("offset")),
	}
	code, result := p.playListService.GetPlaylists(c, filter)
	c.JSON(code, result)
}

// PostPlaylist godoc
// @Summary Post playlist
// @Description Post playlist
// @Tags playlist
// @Id post-playlist
// @Accept json
// @Produce json
// @Param playlist body model.PlaylistRequest true "playlist"
// @Success 200 {object} model.Playlist
// @Router /playlist [post]
func (p *Playlist) PostPlaylist(c *gin.Context) {
	playlist := model.PlaylistRequest{}
	if err := c.BindJSON(&playlist); err != nil {
		code, result := response.BadRequest()
		c.JSON(code, result)
		return
	}
	if err := playlist.Validate(); err != nil {
		code, _ := response.BadRequest()
		c.JSON(code, err)
		return
	}

	code, result := p.playListService.PostPlaylist(c, playlist)
	c.JSON(code, result)
}

// GetPlaylistById godoc
// @Summary Get playlist by id
// @Description Get playlist by id
// @Tags playlist
// @Id get-playlist-id
// @Accept json
// @Produce json
// @Param id path string true "Playlist ID"
// @Success 200 {object} model.Playlist
// @Router /playlist/{id} [get]
func (p *Playlist) GetPlaylistById(c *gin.Context) {
	trackUuid := c.Param("id")
	if trackUuid == "" {
		c.JSON(response.BadRequestMsg("id is missing"))
		c.Abort()
		return
	}
	code, result := p.playListService.GetPlaylistById(c, trackUuid)
	c.JSON(code, result)
}

// DeletePlaylistById godoc
// @Summary Delete playlist by id
// @Description Delete playlist by id
// @Tags playlist
// @Id Delete-playlist-id
// @Accept json
// @Produce json
// @Param id path string true "Playlist ID"
// @Success 200 {object} model.Playlist
// @Router /playlist/{id} [delete]
func (p *Playlist) DeletePlaylistById(c *gin.Context) {
	playlistUuid := c.Param("id")
	if playlistUuid == "" {
		c.JSON(response.BadRequestMsg("id is missing"))
		c.Abort()
		return
	}
	code, result := p.playListService.DeletePlaylistById(c, playlistUuid)
	c.JSON(code, result)
}

// PutPlaylistById godoc
// @Summary Put playlist by id
// @Description Put playlist by id
// @Tags playlist
// @Id put-playlist
// @Accept json
// @Produce json
// @Param id path string true "Playlist ID"
// @Param playlist body model.PlaylistRequest true "playlist"
// @Success 200 {object} model.Playlist
// @Router /playlist/{id} [Put]
func (p *Playlist) PutPlaylistById(c *gin.Context) {
	playlistUuid := c.Param("id")
	playlistUpdate := model.PlaylistRequest{}
	if err := c.BindJSON(&playlistUpdate); err != nil {
		code, result := response.BadRequest()
		c.JSON(code, result)
		return
	}
	if err := playlistUpdate.Validate(); err != nil {
		code, _ := response.BadRequest()
		c.JSON(code, err)
		return
	}

	code, result := p.playListService.PutPlaylistById(c, playlistUuid, playlistUpdate)
	c.JSON(code, result)
}
