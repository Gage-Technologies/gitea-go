// Copyright 2014 The Gogs Authors. All rights reserved.
// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
)

// AccessToken represents an API access token.
type AccessToken struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Token          string `json:"sha1"`
	TokenLastEight string `json:"token_last_eight"`
}

// ListAccessTokensOptions options for listing a users's access tokens
type ListAccessTokensOptions struct {
	ListOptions
}

// ListAccessTokens lists all the access tokens of user
func (c *Client) ListAccessTokens(opts ListAccessTokensOptions) ([]*AccessToken, *Response, error) {
	c.mutex.RLock()
	username := c.username
	c.mutex.RUnlock()
	if len(username) == 0 {
		return nil, nil, fmt.Errorf("\"username\" not set: only BasicAuth allowed")
	}
	opts.setDefaults()
	tokens := make([]*AccessToken, 0, opts.PageSize)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/users/%s/tokens?%s", url.PathEscape(username), opts.getURLQuery().Encode()), jsonHeader, nil, &tokens)
	return tokens, resp, err
}

// ListAccessTokensAdmin lists all the access tokens of any user
func (c *Client) ListAccessTokensAdmin(user string, opts ListAccessTokensOptions) ([]*AccessToken, *Response, error) {
	opts.setDefaults()
	tokens := make([]*AccessToken, 0, opts.PageSize)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/users/%s/tokens?%s", url.PathEscape(user), opts.getURLQuery().Encode()), jsonHeader, nil, &tokens)
	return tokens, resp, err
}

// CreateAccessTokenOption options when create access token
type CreateAccessTokenOption struct {
	Name string `json:"name"`
}

// CreateAccessToken create one access token with options
func (c *Client) CreateAccessToken(opt CreateAccessTokenOption) (*AccessToken, *Response, error) {
	c.mutex.RLock()
	username := c.username
	c.mutex.RUnlock()
	if len(username) == 0 {
		return nil, nil, fmt.Errorf("\"username\" not set: only BasicAuth allowed")
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}
	t := new(AccessToken)
	resp, err := c.getParsedResponse("POST", fmt.Sprintf("/users/%s/tokens", url.PathEscape(username)), jsonHeader, bytes.NewReader(body), t)
	return t, resp, err
}

// CreateAccessTokenAdmin create one access for any user token with options
func (c *Client) CreateAccessTokenAdmin(user string, opt CreateAccessTokenOption) (*AccessToken, *Response, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}
	t := new(AccessToken)
	resp, err := c.getParsedResponse("POST", fmt.Sprintf("/users/%s/tokens", url.PathEscape(user)), jsonHeader, bytes.NewReader(body), t)
	return t, resp, err
}

// DeleteAccessToken delete token, identified by ID and if not available by name
func (c *Client) DeleteAccessToken(value interface{}) (*Response, error) {
	c.mutex.RLock()
	username := c.username
	c.mutex.RUnlock()
	if len(username) == 0 {
		return nil, fmt.Errorf("\"username\" not set: only BasicAuth allowed")
	}

	token := ""

	switch reflect.ValueOf(value).Kind() {
	case reflect.Int64:
		token = fmt.Sprintf("%d", value.(int64))
	case reflect.String:
		if err := c.checkServerVersionGreaterThanOrEqual(version1_13_0); err != nil {
			return nil, err
		}
		token = value.(string)
	default:
		return nil, fmt.Errorf("only string and int64 supported")
	}

	_, resp, err := c.getResponse("DELETE", fmt.Sprintf("/users/%s/tokens/%s", url.PathEscape(username), url.PathEscape(token)), jsonHeader, nil)
	return resp, err
}

// DeleteAccessTokenAdmin delete token for any user, identified by ID and if not available by name
func (c *Client) DeleteAccessTokenAdmin(user string, value interface{}) (*Response, error) {
	token := ""

	switch reflect.ValueOf(value).Kind() {
	case reflect.Int64:
		token = fmt.Sprintf("%d", value.(int64))
	case reflect.String:
		if err := c.checkServerVersionGreaterThanOrEqual(version1_13_0); err != nil {
			return nil, err
		}
		token = value.(string)
	default:
		return nil, fmt.Errorf("only string and int64 supported")
	}

	_, resp, err := c.getResponse("DELETE", fmt.Sprintf("/users/%s/tokens/%s", url.PathEscape(user), url.PathEscape(token)), jsonHeader, nil)
	return resp, err
}
