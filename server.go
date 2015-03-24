// Copyright 2015 The Warren Authors
// Use of this source code is governed by an MIT license that can be found in
// the LICENSE file.

package main

import (
	"io/ioutil"
	"os"

	"github.com/go-martini/martini"
	"github.com/gorilla/sessions"
	"github.com/martini-contrib/secure"
	"gopkg.in/mgo.v2"
	"gopkg.in/yaml.v2"

	"github.com/makyo/warren-go/handlers"
)

type Mongo struct {
	Host string `yaml:"host"`
	DB   string `yaml:"db"`
}

type Config struct {
	EnvironmentType string `yaml:"env-type"`
	AuthKey         string `yaml:"auth-key"`
	EncryptionKey   string `yaml:"encryption-key"`
	Mongo           Mongo
}

var (
	store sessions.Store
	db    *mgo.Database
)

func init() {
	var config Config
	file := os.Args[1]
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	martini.Env = config.EnvironmentType

	store = sessions.NewCookieStore([]byte(config.AuthKey), []byte(config.EncryptionKey))

	dbSession, err := mgo.Dial(config.Mongo.Host)
	if err != nil {
		panic(err)
	}

	db = dbSession.DB(config.Mongo.DB)
}

func main() {
	m := martini.Classic()

	h := handlers.New(store, db)

	m.Get("/", h.Front)

	m.Get("/login", h.DisplayLogin)
	m.Post("/login", h.Login)
	m.Get("/logout", h.Logout)
	m.Get("/register", h.DisplayRegister)
	m.Post("/register", h.Register)
	m.Get("/confirm/:confirmation", h.Confirm)
	m.Get("/~:username", h.DisplayUser)
	m.Post("/~:username/follow", h.FollowUser)
	m.Post("/~:username/unfollow", h.UnfollowUser)
	m.Post("/~:username/friend/request", h.RequestFriendship)
	m.Post("/~:username/friend/confirm", h.ConfirmFriendship)
	m.Post("/~:username/friend/reject", h.RejectFriendship)
	m.Post("/~:username/friend/cancel", h.CancelFriendship)

	m.Get("/(?P<post>\\d+)", h.DisplayPost)
	m.Get("/(?P<post>\\d+)/delete", h.DisplayDeletePost)
	m.Get("/(?P<post>\\d+)/delete", h.DeletePost)
	m.Post("/(?P<post>\\d+)/share", h.SharePost)
	m.Get("/post", h.DisplayCreatePost)
	m.Post("/post", h.CreatePost)

	m.Get("/posts", h.ListAll)
	m.Get("/posts/following", h.ListFollowing)
	m.Get("/posts/friends", h.ListFriends)

	m.Use(secure.Secure(secure.Options{
		//AllowedHosts:          []string{"example.com", "ssl.example.com"},
		//SSLHost:               "ssl.example.com",
		SSLRedirect:           true,
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
	}))
	m.Use(h.SessionMiddleware)
	m.Use(h.AuthenticationMiddleware)

	m.Run()
}
