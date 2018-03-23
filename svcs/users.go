package svcs

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"notif2"
)

type UserService interface {
	FindAll() ([]notif2.User, error)
}

type userServiceImpl struct {
	url *url.URL
}

func (s *userServiceImpl) FindAll() ([]notif2.User, error) {
	type UserDTO struct {
		Name  string `json:"name"`
		EMail string `json:"email"`
	}
	log.Printf("Requesting: %s", s.url.String())
	resp, err := http.Get(s.url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Println("Received", string(body))
	var users []UserDTO
	err = json.Unmarshal(body, &users)
	if err != nil {
		return nil, err
	}
	var result []notif2.User = make([]notif2.User, len(users))
	for i, u := range users {
		result[i].Name = u.Name
		result[i].EMail = u.EMail
	}
	return result, nil
}

func NewBasicUserService(svc string) (UserService, error) {
	if len(svc) > 0 && svc[len(svc)-1] != '/' {
		svc += "/"
	}
	serviceUrl, err := url.Parse(svc)
	if err != nil {
		return nil, err
	}
	usersPath, _ := url.Parse("users")
	userResourceURL := serviceUrl.ResolveReference(usersPath)
	return &userServiceImpl{url: userResourceURL}, nil
}
