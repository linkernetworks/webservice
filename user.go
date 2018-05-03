package acl

import (
	"fmt"
	"net/http"

	"bitbucket.org/linkernetworks/aurora/src/entity"
	"bitbucket.org/linkernetworks/aurora/src/service/mongo"
	"bitbucket.org/linkernetworks/aurora/src/service/session"
	restful "github.com/emicklei/go-restful"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const SessionKey = "ses"

func GetCurrentUserRestful(ses *mongo.Session, req *restful.Request) (*entity.User, error) {
	return GetCurrentUser(ses, req.Request)
}

// GetCurrentUser get current user data with login session and return user data
// excluding sensitive data like password.
func GetCurrentUser(ses *mongo.Session, req *http.Request) (*entity.User, error) {
	email, err := GetCurrentUserEmail(req)
	if err != nil {
		return nil, err
	}

	user := entity.User{}
	q := bson.M{"email": email}
	projection := bson.M{"password": 0}
	if err := ses.C(entity.UserCollectionName).Find(q).Select(projection).One(&user); err != nil {
		if err == mgo.ErrNotFound {
			return nil, fmt.Errorf("user document not found.")
		}
		return nil, err
	}

	return &user, nil
}

// GetCurrentUserWithPassword get current user data with login session and return all user data
// including sensitive data like encrypted password.
func GetCurrentUserWithPassword(ses *mongo.Session, req *http.Request) (*entity.User, error) {
	email, err := GetCurrentUserEmail(req)
	if err != nil {
		return nil, err
	}

	user := entity.User{}
	q := bson.M{"email": email}
	if err := ses.C(entity.UserCollectionName).Find(q).One(&user); err != nil {
		if err == mgo.ErrNotFound {
			return nil, fmt.Errorf("user document not found.")
		}
		return nil, err
	}

	return &user, nil
}

func GetCurrentUserEmail(req *http.Request) (string, error) {
	// FIXME: token is not used, we should use the token to load the actual user.
	token := req.Header.Get("Authorization")
	if len(token) == 0 {
		return "", fmt.Errorf("Authorization token is missing.")
	}

	session, err := session.Service.Store.Get(req, SessionKey)
	if err != nil {
		return "", err
	}

	val, found := session.Values["email"]
	if !found {
		return "", fmt.Errorf("session email is not set.")
	}

	email, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("session email value type is invalid.")
	}
	return email, err
}
