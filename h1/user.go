// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package h1

import (
	"encoding/json"
)

// UserProfilePicture is a nested struct of the User struct
type UserProfilePicture struct {
	Size62x62   *string `json:"62x62"`
	Size82x82   *string `json:"82x82"`
	Size110x110 *string `json:"110x110"`
	Size260x260 *string `json:"260x260"`
}

// User represents an individual user.
//
// HackerOne API docs: https://api.hackerone.com/reference/#user
type User struct {
	ID                    *string            `json:"id"`
	Type                  *string            `json:"type"`
	Disabled              *bool              `json:"disabled"`
	Username              *string            `json:"username"`
	Name                  *string            `json:"name"`
	ProfilePicture        UserProfilePicture `json:"profile_picture"`
	Bio                   *string            `json:"bio,omitempty"`
	Website               *string            `json:"website,omitempty"`
	Location              *string            `json:"location,omitempty"`
	Reputation            *uint64            `json:"reputation,omitempty"`
	Signal                *float64           `json:"signal,omitempty"`
	Impact                *float64           `json:"impact,omitempty"`
	HackeroneTriager      *bool              `json:"hackerone_triager,omitempty"`
	CreatedAt             *Timestamp         `json:"created_at"`
	ParticipatingPrograms []Program          `json:"participating_programs,omitempty"`
}

// Helper types for JSONUnmarshal
type user User // Used to avoid recursion of JSONUnmarshal
type userUnmarshalHelper struct {
	user
	Attributes    *user `json:"attributes"`
	Relationships struct {
		ParticipatingPrograms struct {
			Data []Program `json:"data"`
		} `json:"participating_programs"`
	}
}

// UnmarshalJSON allows JSONAPI attributes and relationships to unmarshal cleanly.
func (u *User) UnmarshalJSON(b []byte) error {
	var helper userUnmarshalHelper
	helper.Attributes = &helper.user
	if err := json.Unmarshal(b, &helper); err != nil {
		return err
	}
	*u = User(helper.user)
	u.ParticipatingPrograms = helper.Relationships.ParticipatingPrograms.Data
	return nil
}
