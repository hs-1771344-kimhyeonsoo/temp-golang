package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"net"

	_movieDelivery "github.com/null-like/movie-backend/movie/delivery"
	_movieRepo "github.com/null-like/movie-backend/movie/repository"
	_movieUsecase "github.com/null-like/movie-backend/movie/usecase"

	_userDelivery "github.com/null-like/movie-backend/user/delivery"
	_userRepo "github.com/null-like/movie-backend/user/repository"
	_userUsecase "github.com/null-like/movie-backend/user/usecase"

	"net/http"
	"os"
	"time"
)

var log *logrus.Logger
var db *sql.DB
var schemaMap map[string]string
var env string

type ViaSSHDialer struct {
	client *ssh.Client
}

func (self *ViaSSHDialer) Dial(addr string) (net.Conn, error) {
	return self.client.Dial("tcp", addr)
}

func init() {
	initLogger()
	initConfig()
	env = os.Args[1]
	db = initDBConnection()
}

func initLogger() {
	log = logrus.New()
	f, err := os.OpenFile("log.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	log.SetOutput(f)
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.DebugLevel)
	log.SetReportCaller(true)
}

func initConfig() {
	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func initDBConnection() *sql.DB {
	var agentClient agent.Agent

	if conn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		//defer conn.Close()
		agentClient = agent.NewClient(conn)
	}

	sshConfig := &ssh.ClientConfig{
		User:            viper.GetString("ssh.user"),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            []ssh.AuthMethod{},
	}

	if agentClient != nil {
		sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeysCallback(agentClient.Signers))
	}

	sshPass := viper.GetString("ssh.pass")
	if sshPass != "" {
		sshConfig.Auth = append(sshConfig.Auth, ssh.PasswordCallback(func() (string, error) {
			return sshPass, nil
		}))
	}

	sshHost := viper.GetString("ssh.host")
	sshPort := viper.GetInt("ssh.port")
	sshcon, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", sshHost, sshPort), sshConfig)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	mysql.RegisterDialContext("mysql+tcp", func(_ context.Context, addr string) (net.Conn, error) {
		dialer := &ViaSSHDialer{sshcon}
		return dialer.Dial(addr)
	})

	conf := mysql.NewConfig()
	conf.User = viper.GetString("movie_db.user")
	conf.Passwd = viper.GetString("movie_db.pass")
	conf.Net = "mysql+tcp"
	conf.Addr = viper.GetString("movie_db.host")
	conf.ParseTime = true

	DB, err := sql.Open("mysql", conf.FormatDSN())
	if err != nil {
		log.Error(err)
		DB.Close()
		os.Exit(1)
	}
	const shortDuration = 1 * time.Second
	d := time.Now().Add(shortDuration)
	ctx, cancel := context.WithDeadline(context.TODO(), d)
	defer cancel()
	err = DB.PingContext(ctx)
	if err != nil {
		log.Error(err)
		DB.Close()
		os.Exit(1)
	}
	schemaMap = viper.GetStringMapString(fmt.Sprintf("movie_db.schema.%s", env))
	return DB
}

func main() {
	e := echo.New()
	if env == "development" {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			Skipper:          nil,
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{http.MethodHead, http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
			AllowHeaders:     []string{"X-Requested-With", "Content-Type", "Authorization"},
			AllowCredentials: false,
			ExposeHeaders:    nil,
			MaxAge:           0,
		}))
	}

	v1 := e.Group("/v1")

	mr := _movieRepo.NewMariaDBMovieRepository(log, db, schemaMap)
	mu := _movieUsecase.NewMovieUsecase(log, mr)
	_movieDelivery.NewMovieHandler(v1, mu)

	ur := _userRepo.NewMariaDBUserRepository(log, db, schemaMap)
	uu := _userUsecase.NewUserUsecase(log, ur)
	_userDelivery.NewUserHandler(v1, uu, log)

	log.Fatal(e.Start(viper.GetString(`server.address`)))
}
