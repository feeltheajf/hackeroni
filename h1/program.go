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

// Program represents a overall program.
//
// HackerOne API docs: https://api.hackerone.com/reference/#program
type Program struct {
	ID        *string    `json:"id"`
	Type      *string    `json:"type"`
	Handle    *string    `json:"handle"`
	Policy    *string    `json:"policy"`
	CreatedAt *Timestamp `json:"created_at"`
	UpdatedAt *Timestamp `json:"updated_at"`
	Groups    []*Group   `json:"groups,omitempty"`
	Members   []*Member  `json:"member,omitempty"`
	// CustomFieldAttributes
	// PolicyAttachments
	// Transactions
}

// Helper types for JSONUnmarshal
type program Program // Used to avoid recursion of JSONUnmarshal
type programUnmarshalHelper struct {
	program
	Attributes    *program `json:"attributes"`
	Relationships struct {
		Groups struct {
			Data []*Group `json:"data"`
		} `json:"groups"`
		Members struct {
			Data []*Member `json:"data"`
		} `json:"members"`
	} `json:"relationships"`
}

// UnmarshalJSON allows JSONAPI attributes and relationships to unmarshal cleanly.
func (p *Program) UnmarshalJSON(b []byte) error {
	var helper programUnmarshalHelper
	helper.Attributes = &helper.program
	if err := json.Unmarshal(b, &helper); err != nil {
		return err
	}
	*p = Program(helper.program)
	p.Groups = helper.Relationships.Groups.Data
	p.Members = helper.Relationships.Members.Data
	return nil
}
