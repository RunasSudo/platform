// Copyright (c) 2017 RunasSudo (Yingtong Li).
// Based on Mattermost Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package oauthgoogle

import (
	"encoding/json"
	"github.com/mattermost/platform/einterfaces"
	"github.com/mattermost/platform/model"
	"io"
	"strings"
)

type GoogleProvider struct {
}

type GoogleUser struct {
	Id     string        `json:"id"`
	Emails []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"emails"`
	Name   string        `json:"displayName"`
}

func init() {
	provider := &GoogleProvider{}
	einterfaces.RegisterOauthProvider(model.USER_AUTH_SERVICE_GOOGLE, provider)
}

func userFromGoogleUser(glu *GoogleUser) *model.User {
	user := &model.User{}
	username := strings.Split(glu.Emails[0].Value, "@")[0]
	user.Username = model.CleanUsername(username)
	splitName := strings.Split(glu.Name, " ")
	if len(splitName) == 2 {
		user.FirstName = splitName[0]
		user.LastName = splitName[1]
	} else if len(splitName) >= 2 {
		user.FirstName = splitName[0]
		user.LastName = strings.Join(splitName[1:], " ")
	} else {
		user.FirstName = glu.Name
	}
	strings.TrimSpace(glu.Emails[0].Value)
	user.Email = glu.Emails[0].Value
	user.AuthData = &glu.Id
	user.AuthService = model.USER_AUTH_SERVICE_GOOGLE

	return user
}

func googleUserFromJson(data io.Reader) *GoogleUser {
	decoder := json.NewDecoder(data)
	var glu GoogleUser
	err := decoder.Decode(&glu)
	if err == nil {
		return &glu
	} else {
		return nil
	}
}

func (glu *GoogleUser) ToJson() string {
	b, err := json.Marshal(glu)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

func (glu *GoogleUser) IsValid() bool {
	if len(glu.Id) == 0 {
		return false
	}

	if len(glu.Emails) == 0 {
		return false
	}

	return true
}

func (glu *GoogleUser) getAuthData() string {
	return glu.Id
}

func (m *GoogleProvider) GetIdentifier() string {
	return model.USER_AUTH_SERVICE_GOOGLE
}

func (m *GoogleProvider) GetUserFromJson(data io.Reader) *model.User {
	glu := googleUserFromJson(data)
	if glu.IsValid() {
		return userFromGoogleUser(glu)
	}

	return &model.User{}
}

func (m *GoogleProvider) GetAuthDataFromJson(data io.Reader) string {
	glu := googleUserFromJson(data)

	if glu.IsValid() {
		return glu.getAuthData()
	}

	return ""
}
