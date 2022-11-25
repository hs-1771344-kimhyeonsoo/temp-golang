package user

import (
	"context"
	userDomain "github.com/null-like/movie-backend/domain/user"
)

type Usecase interface {
	RegisterUser(ctx context.Context, user userDomain.User) error
	CheckUser(ctx context.Context, email string) (bool, error)
	AuthUser(ctx context.Context, email string, password string) (userDomain.UserInfo, error)
	GetAllUsers(ctx context.Context) ([]userDomain.AllUserInfo, error)
	DeleteAndGetAllUsers(ctx context.Context, id int) ([]userDomain.AllUserInfo, error)
	UpdateAndGetAllUsers(ctx context.Context, id int, rank string) ([]userDomain.AllUserInfo, error)
	GetNickName(ctx context.Context, id int) (string, error)
	GetIsFavorite(ctx context.Context, userId int, movieId int, mediaType string) (bool, error)
	GetFavorites(ctx context.Context, userId int) ([]userDomain.Favorite, error)
	ChangeIsLiked(ctx context.Context, userId int, movieId int, isLiked int, mediaType string) error
	GetRating(ctx context.Context, userId int, movieId int, mediaType string) (int, error)
	GetRatingList(ctx context.Context, userId int) ([]userDomain.Rate, error)
	GetChangedRatingList(ctx context.Context, userId int, movieId int, rating int, mediaType string) ([]userDomain.Rate, error)
	GetAllPlaylists(ctx context.Context) ([]userDomain.Playlist, error)
	ChangePlaylistAndGetAllPlaylists(ctx context.Context, id int, title string, playlist string) ([]userDomain.Playlist, error)
	AddPlaylistAndGetAllPlaylists(ctx context.Context, name string, playlist string, mediaType string) ([]userDomain.Playlist, error)
	DeletePlaylistAndGetAllPlaylists(ctx context.Context, id int) ([]userDomain.Playlist, error)
}
