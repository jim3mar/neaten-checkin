package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	mod "github.com/jim3mar/tidy/models/user"
	util "github.com/jim3mar/tidy/utilities"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"encoding/json"
	"log"
	//"strconv"
	"time"
)

type UserResource struct {
	Mongo    *mgo.Session
	CollUser *mgo.Collection
}

func (ur *UserResource) Init(session *mgo.Session) {
	ur.Mongo = session
	ur.CollUser = ur.Mongo.DB("tidy").C("user")
}

// NewUser add a user into mongo/tidy/user
// return current timestamp if success
func (ur *UserResource) NewUser(c *gin.Context) {
	now := time.Now()
	//col := ur.Mongo.DB("tidy").C("user")
	//content := c.PostForm("content")
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")

	if username == "" || password == "" || email == "" {
		c.JSON(http.StatusBadRequest, "Invalid parameter")
	}
	log.Print("New username:" + username)
	log.Print("New password" + password)
	log.Print("New email" + email)
	err := ur.CollUser.Insert(&mod.User{
		Id_:        bson.NewObjectId(),
		UserName:   username,
		Password:   util.Md5Sum(password),
		EMail:      email,
		CreateAt:   now,
		Timestamp:  now.Unix(),
		Portrait:   "avantar.png",
		Continuous: 0,
		//LastCheckIn:  ,
	})

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, now.Unix())
}

type AuthReponse struct {
	AuthToken string   `json:"auth_token"`
	UserInfo  mod.User `json:"user_info"`
}

func (ur *UserResource) AuthWithPassword(c *gin.Context) {
	username := c.DefaultQuery("username", "")
	password := c.DefaultQuery("password", "")
	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, "Invalid username or password")
		return
	}
	password = util.Md5Sum(password)
	user := new(mod.User)
	err := ur.CollUser.Find(
		bson.M{
			"username": username,
			"password": password,
		}).One(user)
	if err != nil {
		c.JSON(http.StatusForbidden, err)
		return
	}
	tokenString, err := util.NewToken(map[string]string{"uid": user.Id_.Hex()})
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK,
		AuthReponse{
			AuthToken: tokenString,
			UserInfo:  *user,
		})
	return
}