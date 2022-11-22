package repository

import (
	"context"
	"database/sql"
	"fmt"
	movieDomain "github.com/null-like/movie-backend/domain/movie"
	"github.com/null-like/movie-backend/movie"
	"github.com/sirupsen/logrus"
)

type mariaDBMovieRepository struct {
	logger    *logrus.Logger
	db        *sql.DB
	schemaMap map[string]string
}

func NewMariaDBMovieRepository(l *logrus.Logger, db *sql.DB, sm map[string]string) movie.Repository {
	return &mariaDBMovieRepository{
		logger:    l,
		db:        db,
		schemaMap: sm,
	}
}

func (r *mariaDBMovieRepository) ReadMovieById(ctx context.Context, movieId int) (movieDomain.Movie, error) {
	query := fmt.Sprintf(`
			SELECT id, adult, genres, title, language, overview, poster, production_companies, release_date, revenue,
				runtime, tagline, rating, votes
			FROM %s.Movie
			WHERE movie_id = %d
		`,
		r.schemaMap["movie"],
		movieId,
	)
	r.logger.Debug(query)
	row := r.db.QueryRowContext(ctx, query)

	var movieInfo movieDomain.Movie
	err := row.Scan(&movieInfo)
	if err != nil {
		r.logger.Error(err)
		return movieInfo, err
	}

	return movieInfo, nil
	//go
}
