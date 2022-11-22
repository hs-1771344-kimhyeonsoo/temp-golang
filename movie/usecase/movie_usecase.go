package usecase

import (
	"context"
	movieDomain "github.com/null-like/movie-backend/domain/movie"
	"github.com/null-like/movie-backend/movie"
	"github.com/sirupsen/logrus"
)

type movieUsecase struct {
	logger    *logrus.Logger
	movieRepo movie.Repository
}

func NewMovieUsecase(l *logrus.Logger, r movie.Repository) movie.Usecase {
	return &movieUsecase{
		logger:    l,
		movieRepo: r,
	}
}

func (u *movieUsecase) GetMovieInfo(ctx context.Context, movieId int) (movieDomain.Movie, error) {
	movieInfo, err := u.movieRepo.ReadMovieById(ctx, movieId)
	if err != nil {
		u.logger.Error(err)
		return movieInfo, err
	}
	return movieInfo, nil
}
