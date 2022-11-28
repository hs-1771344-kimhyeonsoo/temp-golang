package user

import (
	"context"
	userDomain "github.com/null-like/movie-backend/domain/user"
)

type Repository interface {
	InsertUser(ctx context.Context, user userDomain.User) error
	FindIdByEmail(ctx context.Context, email string) (bool, error)
	FindIdAndPasswdByEmail(ctx context.Context, email string) (int, string, string, string, string, error)
	FindAllUser(ctx context.Context) ([]userDomain.AllUserInfo, error)
	DeleteUser(ctx context.Context, id int) error
	UpdateUser(ctx context.Context, id int, rank string) error
	FindNicknameByUserId(ctx context.Context, userId int) (string, error)

	FindIsFavorite(ctx context.Context, userId int, movieId int, mediaType string) (bool, error)
	FindFavoriteByUserId(ctx context.Context, userId int) ([]userDomain.Favorite, error)
	InsertFavorite(ctx context.Context, userId int, movieId int, mediaType string) error
	DeleteFavorite(ctx context.Context, userId int, movieId int, mediaType string) error

	InsertRating(ctx context.Context, userId int, movieId int, rating int, mediaType string) error
	FindRatingByMovieId(ctx context.Context, userId int, movieId int, mediaType string) (int, error)
	FindRatingsByUserId(ctx context.Context, userId int) ([]userDomain.Rate, error)

	AllPlaylist(ctx context.Context) ([]userDomain.Playlist, error)
	InsertPlaylist(ctx context.Context, id int, title string, playlist string) error
	InsertPlaylist2(ctx context.Context, name string, playlist string, mediaType string) error
	DeletePlaylist(ctx context.Context, id int) error

	AllBanner(ctx context.Context) ([]userDomain.Banner, error)
	UpdateBanner(ctx context.Context, id int, movieId int, title string, mediaType string, comment string) error
	InsertBanner(ctx context.Context, movieId int, title string, mediaType string, comment string) error
	DeleteBanner(ctx context.Context, id int) error
}
