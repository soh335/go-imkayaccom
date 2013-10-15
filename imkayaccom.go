package imkayaccom

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const (
	URL = "http://im.kayac.com"
)

type Client struct {
	builder Builder
	user    string
}

type Builder interface {
	build(message string, handler string) url.Values
}

type NoPasswordBuild struct {
}

type PasswordBuild struct {
	password string
}

type SecretBuild struct {
	secret string
}

func (builder *NoPasswordBuild) build(message string, handler string) url.Values {
	v := url.Values{}
	v.Set("message", message)
	v.Set("handler", handler)
	return v
}

func (builder *PasswordBuild) build(message string, handler string) url.Values {
	v := url.Values{}
	v.Set("message", message)
	v.Set("handler", handler)
	v.Set("password", builder.password)
	return v
}

func (builder *SecretBuild) build(message string, handler string) url.Values {
	hash := sha1.New()
	hash.Write([]byte(message + builder.secret))

	v := url.Values{}
	v.Set("message", message)
	v.Set("handler", handler)
	v.Set("sig", fmt.Sprintf("%x", hash.Sum(nil)))
	return v
}

func NewNoPasswordClient(user string) *Client {
	client := &Client{}
	client.user = user
	client.builder = &NoPasswordBuild{}
	return client
}

func NewPasswordClient(user string, password string) *Client {
	client := &Client{}
	client.user = user
	builder := &PasswordBuild{password}
	client.builder = builder
	return client
}

func NewSecrentClient(user string, secret string) *Client {
	client := &Client{}
	client.user = user
	builder := &SecretBuild{secret}
	client.builder = builder
	return client
}

func (client *Client) Post(message string, handler string) error {
	resp, err := http.PostForm(
		URL+"/api/post/"+client.user, client.builder.build(message, handler),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var data map[string]string
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&data); err != nil {
		return err
	}

	if err, ok := data["error"]; ok == true && err != "" {
		return errors.New(err)
	}

	return nil
}
