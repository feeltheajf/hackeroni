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

// Credential represents a credential.
//
// HackerOne API docs:  https://api.hackerone.com/customer-reference/#credential
type Credential struct {
	ID               *string      `json:"id"`
	Type             *string      `json:"type"`
	Credentials      *Credentials `json:"credentials"`
	Revoked          *bool        `json:"revoked"`
	AssigneeID       *string      `json:"assignee_id"`
	AssigneeUsername *string      `json:"assignee_username"`
}

type Credentials struct {
	Table *CredentialsTable `json:"table"`
}

type CredentialsTable struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}

// Helper types for JSONUnmarshal
type credential Credential // Used to avoid recursion of JSONUnmarshal
type credentialUnmarshalHelper struct {
	credential
	Attributes *credential `json:"attributes"`
}

// UnmarshalJSON allows JSONAPI attributes and relationships to unmarshal cleanly.
func (ci *Credential) UnmarshalJSON(b []byte) error {
	var helper credentialUnmarshalHelper
	helper.Attributes = &helper.credential
	if err := json.Unmarshal(b, &helper); err != nil {
		return err
	}
	*ci = Credential(helper.credential)
	return nil
}
