package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	userDomain "github.com/null-like/movie-backend/domain/user"
	"github.com/null-like/movie-backend/user"
	"github.com/sirupsen/logrus"
)

type userUsecase struct {
	logger   *logrus.Logger
	userRepo user.Repository
}

func NewUserUsecase(l *logrus.Logger, r user.Repository) user.Usecase {
	return &userUsecase{
		logger:   l,
		userRepo: r,
	}
}

func (u *userUsecase) RegisterUser(ctx context.Context, user userDomain.User) error {
	hash := sha256.New()
	hash.Write([]byte(user.Password))
	user.Password = hex.EncodeToString(hash.Sum(nil))
	err := u.userRepo.InsertUser(ctx, user)
	if err != nil {
		u.logger.Error(err)
		return err
	}

	return nil
}

func (u *userUsecase) CheckUser(ctx context.Context, email string) (bool, error) {
	isExist, err := u.userRepo.FindIdByEmail(ctx, email)
	return isExist, err
}

func (u *userUsecase) AuthUser(ctx context.Context, email string, password string) (userDomain.UserInfo, error) {
	var hashPassword string
	hash := sha256.New()
	hash.Write([]byte(password))
	hashPassword = hex.EncodeToString(hash.Sum(nil))
	id, email, dbPassword, nickname, rank, err := u.userRepo.FindIdAndPasswdByEmail(ctx, email)

	var userInfo userDomain.UserInfo

	if err != nil {
		u.logger.Error(err)
		userInfo.Id = -1
		userInfo.Email = ""
		userInfo.Nickname = ""
		userInfo.Rank = ""
		return userInfo, err
	}

	if dbPassword == hashPassword {
		userInfo.Id = id
		userInfo.Email = email
		userInfo.Nickname = nickname
		userInfo.Rank = rank
		return userInfo, nil
	} else {
		userInfo.Id = -1
		userInfo.Email = ""
		userInfo.Nickname = ""
		userInfo.Rank = ""
		return userInfo, nil
	}
}

func (u *userUsecase) GetAllUsers(ctx context.Context) ([]userDomain.AllUserInfo, error) {
	users, err := u.userRepo.FindAllUser(ctx)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}
	return users, nil
}

func (u *userUsecase) DeleteAndGetAllUsers(ctx context.Context, id int) ([]userDomain.AllUserInfo, error) {
	err := u.userRepo.DeleteUser(ctx, id)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}

	users, err := u.userRepo.FindAllUser(ctx)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}
	return users, nil
}

func (u *userUsecase) UpdateAndGetAllUsers(ctx context.Context, id int, rank string) ([]userDomain.AllUserInfo, error) {
	err := u.userRepo.UpdateUser(ctx, id, rank)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}

	users, err := u.userRepo.FindAllUser(ctx)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}
	return users, nil
}

func (u *userUsecase) GetNickName(ctx context.Context, id int) (string, error) {
	nickname, err := u.userRepo.FindNicknameByUserId(ctx, id)
	if err != nil {
		u.logger.Error(err)
		return nickname, err
	}
	return nickname, nil
}

func (u *userUsecase) GetIsFavorite(ctx context.Context, userId int, movieId int, mediaType string) (bool, error) {
	isFavorite, err := u.userRepo.FindIsFavorite(ctx, userId, movieId, mediaType)
	return isFavorite, err
}

func (u *userUsecase) GetFavorites(ctx context.Context, userId int) ([]userDomain.Favorite, error) {
	movieId, err := u.userRepo.FindFavoriteByUserId(ctx, userId)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}
	return movieId, nil
}

func (u *userUsecase) ChangeIsLiked(ctx context.Context, userId int, movieId int, isLiked int, mediaType string) error {
	var err error

	if isLiked == 1 {
		err = u.userRepo.InsertFavorite(ctx, userId, movieId, mediaType)
	} else {
		err = u.userRepo.DeleteFavorite(ctx, userId, movieId, mediaType)
	}

	if err != nil {
		u.logger.Error(err)
		return err
	}
	return nil
}

func (u *userUsecase) GetRating(ctx context.Context, userId int, movieId int, mediaType string) (int, error) {
	rating, err := u.userRepo.FindRatingByMovieId(ctx, userId, movieId, mediaType)
	if err != nil {
		u.logger.Error(err)
		return 0, err
	}
	return rating, nil
}

func (u *userUsecase) GetRatingList(ctx context.Context, userId int) ([]userDomain.Rate, error) {
	movieRatings, err := u.userRepo.FindRatingsByUserId(ctx, userId)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}
	return movieRatings, nil
}

func (u *userUsecase) GetChangedRatingList(ctx context.Context, userId int, movieId int, rating int, mediaType string) ([]userDomain.Rate, error) {
	err := u.userRepo.InsertRating(ctx, userId, movieId, rating, mediaType)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}

	movieRatings, err := u.userRepo.FindRatingsByUserId(ctx, userId)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}
	return movieRatings, nil
}

func (u *userUsecase) GetAllPlaylists(ctx context.Context) ([]userDomain.Playlist, error) {
	playlists, err := u.userRepo.AllPlaylist(ctx)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}

	return playlists, nil
}

func (u *userUsecase) ChangePlaylistAndGetAllPlaylists(ctx context.Context, id int, name string, playlist string) ([]userDomain.Playlist, error) {
	err := u.userRepo.InsertPlaylist(ctx, id, name, playlist)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}

	playlists, err := u.userRepo.AllPlaylist(ctx)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}

	return playlists, nil
}

func (u *userUsecase) AddPlaylistAndGetAllPlaylists(ctx context.Context, name string, playlist string, mediaType string) ([]userDomain.Playlist, error) {
	err := u.userRepo.InsertPlaylist2(ctx, name, playlist, mediaType)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}

	playlists, err := u.userRepo.AllPlaylist(ctx)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}

	return playlists, nil
}

func (u *userUsecase) DeletePlaylistAndGetAllPlaylists(ctx context.Context, id int) ([]userDomain.Playlist, error) {
	err := u.userRepo.DeletePlaylist(ctx, id)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}

	playlists, err := u.userRepo.AllPlaylist(ctx)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}

	return playlists, nil
}

func (u *userUsecase) GetAllBanners(ctx context.Context) ([]userDomain.Banner, error) {
	banners, err := u.userRepo.AllBanner(ctx)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}

	return banners, nil
}

func (u *userUsecase) UpdateAndGetAllBanners(ctx context.Context, id int, movieId int, title string, mediaType string, comment string) ([]userDomain.Banner, error) {
	err := u.userRepo.UpdateBanner(ctx, id, movieId, title, mediaType, comment)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}

	banners, err := u.userRepo.AllBanner(ctx)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}

	return banners, nil
}

func (u *userUsecase) AddAndGetAllBanners(ctx context.Context, movieId int, title string, mediaType string, comment string) ([]userDomain.Banner, error) {
	err := u.userRepo.InsertBanner(ctx, movieId, title, mediaType, comment)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}

	banners, err := u.userRepo.AllBanner(ctx)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}

	return banners, nil
}

func (u *userUsecase) DeleteAndGetAllBanners(ctx context.Context, id int) ([]userDomain.Banner, error) {
	err := u.userRepo.DeleteBanner(ctx, id)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}

	banners, err := u.userRepo.AllBanner(ctx)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}

	return banners, nil
}
