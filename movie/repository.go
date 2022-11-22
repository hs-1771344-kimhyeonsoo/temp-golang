package movie

import (
	"context"
	movieDomain "github.com/null-like/movie-backend/domain/movie"
)

type Repository interface {
	ReadMovieById(ctx context.Context, movieId int) (movieDomain.Movie, error)
}
