package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/null-like/movie-backend/movie"
	"net/http"
	"strconv"
)

type movieHandler struct {
	Usecase movie.Usecase
}

type ResponseError struct {
	Message string `json:"message"`
}

func NewMovieHandler(g *echo.Group, u movie.Usecase) {
	handler := &movieHandler{
		Usecase: u,
	}
	g.GET("/movie/movie-info", handler.GetMovieInfo)
}

func (h *movieHandler) GetMovieInfo(c echo.Context) error {
	ctx := c.Request().Context()
	pathParam := c.Param("movie_id")
	movieId, err := strconv.Atoi(pathParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	movieInfo, err := h.Usecase.GetMovieInfo(ctx, movieId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, movieInfo)
}
