package delivery

import (
	"github.com/labstack/echo/v4"
	UserDomain "github.com/null-like/movie-backend/domain/user"
	"github.com/null-like/movie-backend/user"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type userHandler struct {
	Usecase user.Usecase
	logger  *logrus.Logger
}

type ResponseError struct {
	Message string `json:"message"`
}

func NewUserHandler(g *echo.Group, u user.Usecase, logger *logrus.Logger) {
	handler := &userHandler{
		Usecase: u,
		logger:  logger,
	}

	g.POST("/sign-up", handler.SignUp)
	g.POST("/sign-in", handler.SignIn)
	g.POST("/check", handler.CheckDuplicates)
	g.GET("/all-user", handler.SendAllUser)
	g.GET("/delete-user", handler.SendDeletedAllUser)
	g.GET("/update-user", handler.SendUpdatedAllUser)
	g.GET("/nickname", handler.SendNickname)
	g.GET("/favorite", handler.SendFavorite)
	g.GET("/is-favorite", handler.SendIsFavorite)
	g.GET("/toggle-fav", handler.ToggleIsLiked)
	g.GET("/rating", handler.SendRating)
	g.GET("/rating-list", handler.SendRatings)
	g.GET("/rating-list-changed", handler.SendChangedRatings)
	g.GET("/playlist", handler.SendPlaylists)
	g.GET("/add-playlist", handler.SendAddedPlaylists)
	g.GET("/change-playlist", handler.SendChangedPlaylists)
	g.GET("/delete-playlist", handler.SendDeletedPlaylists)
}

func (h *userHandler) SignUp(c echo.Context) error {
	ctx := c.Request().Context()

	params := c.QueryParams()

	user := UserDomain.User{
		Email:    params.Get("email"),
		Password: params.Get("password"),
		Nickname: params.Get("nickname"),
	}

	err := h.Usecase.RegisterUser(ctx, user)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(http.StatusBadRequest, -1)
	}

	return c.JSON(http.StatusOK, 1)
}

func (h *userHandler) CheckDuplicates(c echo.Context) error {
	ctx := c.Request().Context()
	params := c.QueryParams()
	isExist, _ := h.Usecase.CheckUser(ctx, params.Get("email"))

	if isExist {
		return c.JSON(http.StatusOK, -1)
	} else {
		return c.JSON(http.StatusOK, 1)
	}
}

func (h *userHandler) SignIn(c echo.Context) error {
	ctx := c.Request().Context()
	params := c.QueryParams()
	userInfo, err := h.Usecase.AuthUser(ctx, params.Get("email"), params.Get("password"))

	if err != nil {
		h.logger.Error(err)
		return c.JSON(http.StatusBadRequest, userInfo)
	}
	return c.JSON(http.StatusOK, userInfo)
}

func (h *userHandler) SendAllUser(c echo.Context) error {
	ctx := c.Request().Context()
	users, err := h.Usecase.GetAllUsers(ctx)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(http.StatusOK, nil)
	}

	return c.JSON(http.StatusOK, users)
}

func (h *userHandler) SendDeletedAllUser(c echo.Context) error {
	ctx := c.Request().Context()
	params := c.QueryParams()
	id, err := strconv.Atoi(params.Get("id"))
	users, err := h.Usecase.DeleteAndGetAllUsers(ctx, id)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(http.StatusOK, nil)
	}

	return c.JSON(http.StatusOK, users)
}

func (h *userHandler) SendUpdatedAllUser(c echo.Context) error {
	ctx := c.Request().Context()
	params := c.QueryParams()
	id, err := strconv.Atoi(params.Get("id"))
	users, err := h.Usecase.UpdateAndGetAllUsers(ctx, id, params.Get("rank"))
	if err != nil {
		h.logger.Error(err)
		return c.JSON(http.StatusOK, nil)
	}

	return c.JSON(http.StatusOK, users)
}

func (h *userHandler) SendNickname(c echo.Context) error {
	ctx := c.Request().Context()
	params := c.QueryParams()
	id, err := strconv.Atoi(params.Get("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}
	nickname, err := h.Usecase.GetNickName(ctx, id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, "")
	} else {
		return c.JSON(http.StatusOK, nickname)
	}
}

func (h *userHandler) SendIsFavorite(c echo.Context) error {
	ctx := c.Request().Context()
	params := c.QueryParams()
	userId, _ := strconv.Atoi(params.Get("user_id"))
	movieId, _ := strconv.Atoi(params.Get("movie_id"))
	mediaType := params.Get("type")

	isFavorite, _ := h.Usecase.GetIsFavorite(ctx, userId, movieId, mediaType)
	return c.JSON(http.StatusOK, isFavorite)
}

func (h *userHandler) SendFavorite(c echo.Context) error {
	ctx := c.Request().Context()
	params := c.QueryParams()
	id, err := strconv.Atoi(params.Get("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	movies, err := h.Usecase.GetFavorites(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, movies)
}

func (h *userHandler) ToggleIsLiked(c echo.Context) error {
	ctx := c.Request().Context()
	params := c.QueryParams()
	userId, _ := strconv.Atoi(params.Get("user_id"))
	movieId, _ := strconv.Atoi(params.Get("movie_id"))
	isLiked, _ := strconv.Atoi(params.Get("is_liked"))
	mediaType := params.Get("type")

	err := h.Usecase.ChangeIsLiked(ctx, userId, movieId, isLiked, mediaType)
	if err != nil {
		h.logger.Error(err)
		return err
	}

	movies, err := h.Usecase.GetFavorites(ctx, userId)
	if err != nil {
		h.logger.Error(err)
		return err
	}

	return c.JSON(http.StatusOK, movies)
}

func (h *userHandler) SendRating(c echo.Context) error {
	ctx := c.Request().Context()
	params := c.QueryParams()
	userId, _ := strconv.Atoi(params.Get("user_id"))
	movieId, _ := strconv.Atoi(params.Get("movie_id"))
	mediaType := params.Get("type")

	rating, err := h.Usecase.GetRating(ctx, userId, movieId, mediaType)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(http.StatusOK, 0)
	}

	return c.JSON(http.StatusOK, rating)
}

func (h *userHandler) SendRatings(c echo.Context) error {
	ctx := c.Request().Context()
	params := c.QueryParams()
	userId, _ := strconv.Atoi(params.Get("user_id"))

	ratings, err := h.Usecase.GetRatingList(ctx, userId)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(http.StatusOK, nil)
	}

	return c.JSON(http.StatusOK, ratings)
}

func (h *userHandler) SendChangedRatings(c echo.Context) error {
	ctx := c.Request().Context()
	params := c.QueryParams()
	userId, _ := strconv.Atoi(params.Get("user_id"))
	movieId, _ := strconv.Atoi(params.Get("movie_id"))
	rating, _ := strconv.Atoi(params.Get("rating"))
	mediaType := params.Get("type")

	ratings, err := h.Usecase.GetChangedRatingList(ctx, userId, movieId, rating, mediaType)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(http.StatusOK, nil)
	}

	return c.JSON(http.StatusOK, ratings)
}

func (h *userHandler) SendPlaylists(c echo.Context) error {
	ctx := c.Request().Context()

	playlists, err := h.Usecase.GetAllPlaylists(ctx)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(http.StatusOK, nil)
	}

	return c.JSON(http.StatusOK, playlists)
}

func (h *userHandler) SendChangedPlaylists(c echo.Context) error {
	ctx := c.Request().Context()
	params := c.QueryParams()
	id, _ := strconv.Atoi(params.Get("id"))
	name := params.Get("name")
	playlist := params.Get("playlist")

	playlists, err := h.Usecase.ChangePlaylistAndGetAllPlaylists(ctx, id, name, playlist)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(http.StatusOK, nil)
	}

	return c.JSON(http.StatusOK, playlists)
}

func (h *userHandler) SendAddedPlaylists(c echo.Context) error {
	ctx := c.Request().Context()
	params := c.QueryParams()
	name := params.Get("name")
	playlist := params.Get("playlist")
	mediaType := params.Get("type")

	playlists, err := h.Usecase.AddPlaylistAndGetAllPlaylists(ctx, name, playlist, mediaType)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(http.StatusOK, nil)
	}

	return c.JSON(http.StatusOK, playlists)
}

func (h *userHandler) SendDeletedPlaylists(c echo.Context) error {
	ctx := c.Request().Context()
	params := c.QueryParams()
	id, _ := strconv.Atoi(params.Get("id"))

	playlists, err := h.Usecase.DeletePlaylistAndGetAllPlaylists(ctx, id)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(http.StatusOK, nil)
	}

	return c.JSON(http.StatusOK, playlists)
}
