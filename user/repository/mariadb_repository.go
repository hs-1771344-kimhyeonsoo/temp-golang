package repository

import (
	"context"
	"database/sql"
	"fmt"
	userDomain "github.com/null-like/movie-backend/domain/user"
	"github.com/null-like/movie-backend/user"
	"github.com/sirupsen/logrus"
)

type mariaDBUserRepository struct {
	logger    *logrus.Logger
	Conn      *sql.DB
	schemaMap map[string]string
}

func NewMariaDBUserRepository(l *logrus.Logger, Conn *sql.DB, sm map[string]string) user.Repository {
	return &mariaDBUserRepository{
		logger:    l,
		Conn:      Conn,
		schemaMap: sm,
	}
}

func (r *mariaDBUserRepository) InsertUser(ctx context.Context, user userDomain.User) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.User (email, password, nickname, rank, signup_date)
		VALUES ('%s', '%s', '%s', '회원', now());
		`,
		r.schemaMap["movie"],
		user.Email,
		user.Password,
		user.Nickname,
	)
	r.logger.Debug(query)

	_, err := r.Conn.ExecContext(ctx, query)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}

func (r *mariaDBUserRepository) FindAllUser(ctx context.Context) ([]userDomain.AllUserInfo, error) {
	query := fmt.Sprintf(`
		SELECT id, email, nickname, rank, signup_date
		FROM %s.User
		`,
		r.schemaMap["movie"],
	)

	rows, err := r.Conn.QueryContext(ctx, query)
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			r.logger.Error(err)
		}
	}()

	var users []userDomain.AllUserInfo
	for rows.Next() {
		var user userDomain.AllUserInfo
		err = rows.Scan(&user.Id, &user.Email, &user.Nickname, &user.Rank, &user.SignUpDate)
		if err != nil {
			r.logger.Error(err)
			return nil, err
		}
		users = append(users, user)
	}

	r.logger.Debug(query)
	return users, nil
}

func (r *mariaDBUserRepository) DeleteUser(ctx context.Context, id int) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.User
		WHERE id = %d;
		`,
		r.schemaMap["movie"],
		id,
	)
	r.logger.Debug(query)

	_, err := r.Conn.ExecContext(ctx, query)
	return err
}

func (r *mariaDBUserRepository) UpdateUser(ctx context.Context, id int, rank string) error {
	query := fmt.Sprintf(`
		UPDATE %s.User
		SET rank = '%s'
		WHERE id = %d;
		`,
		r.schemaMap["movie"],
		rank,
		id,
	)
	r.logger.Debug(query)

	_, err := r.Conn.ExecContext(ctx, query)
	return err
}

func (r *mariaDBUserRepository) FindIdByEmail(ctx context.Context, email string) (bool, error) {
	query := fmt.Sprintf(`
		SELECT id
		FROM %s.User
		WHERE email = '%s';
		`,
		r.schemaMap["movie"],
		email,
	)

	var id int
	r.logger.Debug(query)
	row := r.Conn.QueryRowContext(ctx, query)
	err := row.Scan(&id)

	if err != nil {
		return false, nil
	}

	return true, nil
}

func (r *mariaDBUserRepository) FindIdAndPasswdByEmail(ctx context.Context, email string) (int, string, string, string, string, error) {
	query := fmt.Sprintf(`
		SELECT id, password, nickname, rank
		FROM %s.User
		WHERE email = '%s';
		`,
		r.schemaMap["movie"],
		email,
	)

	var id int
	var nickname string
	var rank string
	var hashPassword string
	r.logger.Debug(query)
	row := r.Conn.QueryRowContext(ctx, query)
	err := row.Scan(&id, &hashPassword, &nickname, &rank)

	if err != nil {
		r.logger.Error(err)
		return -1, email, "", "", "", err
	}

	return id, email, hashPassword, nickname, rank, nil
}

func (r *mariaDBUserRepository) FindNicknameByUserId(ctx context.Context, userId int) (string, error) {
	query := fmt.Sprintf(`
		SELECT nickname
		FROM %s.User
		WHERE id = %d;
		`,
		r.schemaMap["movie"],
		userId,
	)

	var nickname string
	r.logger.Debug(query)
	row := r.Conn.QueryRowContext(ctx, query)
	err := row.Scan(&nickname)

	if err != nil {
		return "", err
	}
	return nickname, nil
}

func (r *mariaDBUserRepository) FindIsFavorite(ctx context.Context, userId int, movieId int, mediaType string) (bool, error) {
	query := fmt.Sprintf(`
		SELECT user_id
		FROM %s.Favorite
		WHERE user_id = %d and movie_id = %d and type = '%s';
		`,
		r.schemaMap["movie"],
		userId,
		movieId,
		mediaType,
	)

	var id int
	r.logger.Debug(query)
	row := r.Conn.QueryRowContext(ctx, query)
	err := row.Scan(&id)

	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *mariaDBUserRepository) FindFavoriteByUserId(ctx context.Context, userId int) ([]userDomain.Favorite, error) {
	query := fmt.Sprintf(`
		SELECT movie_id, type
		FROM %s.Favorite
		WHERE user_id = %d;
		`,
		r.schemaMap["movie"],
		userId,
	)

	rows, err := r.Conn.QueryContext(ctx, query)
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			r.logger.Error(err)
		}
	}()

	var movies []userDomain.Favorite
	for rows.Next() {
		var movie userDomain.Favorite
		err = rows.Scan(&movie.Id, &movie.Type)
		if err != nil {
			r.logger.Error(err)
			return nil, err
		}
		movies = append(movies, movie)
	}

	r.logger.Debug(query)
	return movies, nil
}

func (r *mariaDBUserRepository) InsertFavorite(ctx context.Context, userId int, movieId int, mediaType string) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.Favorite (user_id, movie_id, type)
		VALUES (%d, %d, '%s');
		`,
		r.schemaMap["movie"],
		userId,
		movieId,
		mediaType,
	)
	r.logger.Debug(query)

	_, err := r.Conn.ExecContext(ctx, query)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}

func (r *mariaDBUserRepository) DeleteFavorite(ctx context.Context, userId int, movieId int, mediaType string) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.Favorite
		WHERE user_id = %d and movie_id = %d and type = '%s';
		`,
		r.schemaMap["movie"],
		userId,
		movieId,
		mediaType,
	)
	r.logger.Debug(query)

	_, err := r.Conn.ExecContext(ctx, query)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}

func (r *mariaDBUserRepository) InsertRating(ctx context.Context, userId int, movieId int, rating int, mediaType string) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.Rate VALUES (%d, %d, %d, '%s', now())
		ON DUPLICATE KEY UPDATE rating = VALUES(rating), apply_date = VALUES(apply_date);
		`,
		r.schemaMap["movie"],
		userId,
		movieId,
		rating,
		mediaType,
	)
	r.logger.Debug(query)

	_, err := r.Conn.ExecContext(ctx, query)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}

func (r *mariaDBUserRepository) FindRatingByMovieId(ctx context.Context, userId int, movieId int, mediaType string) (int, error) {
	query := fmt.Sprintf(`
		SELECT rating
		FROM %s.Rate
		WHERE user_id = %d and movie_id = %d and type = '%s';
		`,
		r.schemaMap["movie"],
		userId,
		movieId,
		mediaType,
	)

	var rating int
	r.logger.Debug(query)
	row := r.Conn.QueryRowContext(ctx, query)
	err := row.Scan(&rating)

	if err != nil {
		return 0, err
	}
	return rating, nil
}

func (r *mariaDBUserRepository) FindRatingsByUserId(ctx context.Context, userId int) ([]userDomain.Rate, error) {
	query := fmt.Sprintf(`
		SELECT movie_id, rating, type, apply_date
		FROM %s.Rate
		WHERE user_id = %d;
		`,
		r.schemaMap["movie"],
		userId,
	)

	rows, err := r.Conn.QueryContext(ctx, query)

	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			r.logger.Error(err)
		}
	}()

	var movieRatings []userDomain.Rate
	for rows.Next() {
		var movieRating userDomain.Rate
		err = rows.Scan(&movieRating.Id, &movieRating.Rating, &movieRating.Type, &movieRating.ApplyDate)
		if err != nil {
			r.logger.Error(err)
			return nil, err
		}
		movieRatings = append(movieRatings, movieRating)
	}

	r.logger.Debug(query)
	return movieRatings, nil
}

func (r *mariaDBUserRepository) AllPlaylist(ctx context.Context) ([]userDomain.Playlist, error) {
	query := fmt.Sprintf(`
		SELECT *
		FROM %s.Playlist
		`,
		r.schemaMap["movie"],
	)

	rows, err := r.Conn.QueryContext(ctx, query)

	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			r.logger.Error(err)
		}
	}()

	var playlists []userDomain.Playlist
	for rows.Next() {
		var playlist userDomain.Playlist
		err = rows.Scan(&playlist.Id, &playlist.Name, &playlist.Playlist, &playlist.Type)
		if err != nil {
			r.logger.Error(err)
			return nil, err
		}
		playlists = append(playlists, playlist)
	}

	r.logger.Debug(query)
	return playlists, nil
}

func (r *mariaDBUserRepository) InsertPlaylist(ctx context.Context, id int, name string, playlist string) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.Playlist (id, name, list) VALUES (%d, '%s', '%s')
		ON DUPLICATE KEY UPDATE name = VALUES(name), list = VALUES(list);
		`,
		r.schemaMap["movie"],
		id,
		name,
		playlist,
	)
	r.logger.Debug(query)

	_, err := r.Conn.ExecContext(ctx, query)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}

func (r *mariaDBUserRepository) InsertPlaylist2(ctx context.Context, name string, playlist string, mediaType string) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.Playlist (name, list, type) VALUES ('%s', '%s', '%s');
		`,
		r.schemaMap["movie"],
		name,
		playlist,
		mediaType,
	)
	r.logger.Debug(query)

	_, err := r.Conn.ExecContext(ctx, query)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}

func (r *mariaDBUserRepository) DeletePlaylist(ctx context.Context, id int) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.Playlist
		WHERE id = %d;
		`,
		r.schemaMap["movie"],
		id,
	)
	r.logger.Debug(query)

	_, err := r.Conn.ExecContext(ctx, query)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}

func (r *mariaDBUserRepository) AllBanner(ctx context.Context) ([]userDomain.Banner, error) {
	query := fmt.Sprintf(`
		SELECT *
		FROM %s.Banner
		`,
		r.schemaMap["movie"],
	)

	rows, err := r.Conn.QueryContext(ctx, query)

	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			r.logger.Error(err)
		}
	}()

	var banners []userDomain.Banner
	for rows.Next() {
		var banner userDomain.Banner
		err = rows.Scan(&banner.Id, &banner.MovieId, &banner.Title, &banner.Type, &banner.Comment)
		if err != nil {
			r.logger.Error(err)
			return nil, err
		}
		banners = append(banners, banner)
	}

	r.logger.Debug(query)
	return banners, nil
}

func (r *mariaDBUserRepository) UpdateBanner(ctx context.Context, id int, movieId int, title string, mediaType string, comment string) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.Banner (id, movie_id, title, type, comment) VALUES (%d, %d, '%s', '%s', '%s')
		ON DUPLICATE KEY UPDATE movie_id = VALUES(movie_id), title = VALUES(title), type = VALUES(type), comment = VALUES(comment);
		`,
		r.schemaMap["movie"],
		id,
		movieId,
		title,
		mediaType,
		comment,
	)
	r.logger.Debug(query)

	_, err := r.Conn.ExecContext(ctx, query)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}

func (r *mariaDBUserRepository) InsertBanner(ctx context.Context, movieId int, title string, mediaType string, comment string) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.Banner (movie_id, title, type, comment) VALUES (%d, '%s', '%s', '%s');
		`,
		r.schemaMap["movie"],
		movieId,
		title,
		mediaType,
		comment,
	)
	r.logger.Debug(query)

	_, err := r.Conn.ExecContext(ctx, query)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}

func (r *mariaDBUserRepository) DeleteBanner(ctx context.Context, id int) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.Banner
		WHERE id = %d;
		`,
		r.schemaMap["movie"],
		id,
	)
	r.logger.Debug(query)

	_, err := r.Conn.ExecContext(ctx, query)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	return nil
}
