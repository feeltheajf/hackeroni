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
	"fmt"
)

// ProgramService handles communication with the program related methods of the H1 API.
type ProgramService service

// Me fetches a list of programs available to the client
func (s *ProgramService) Me() ([]Program, *Response, error) {
	req, err := s.client.NewRequest("GET", "me/programs", nil)
	if err != nil {
		return nil, nil, err
	}

	data := new([]Program)
	resp, err := s.client.Do(req, data)
	if err != nil {
		return nil, resp, err
	}

	return *data, resp, err
}

// Get fetches a Program by ID
func (s *ProgramService) Get(ID string) (*Program, *Response, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("programs/%s", ID), nil)
	if err != nil {
		return nil, nil, err
	}

	data := new(Program)
	resp, err := s.client.Do(req, data)
	if err != nil {
		return nil, resp, err
	}

	return data, resp, err
}

// ListStructuredScopes fetches a list of structured scopes for the given program
func (s *ProgramService) ListStructuredScopes(programID string, listOpts *ListOptions) ([]StructuredScope, *Response, error) {
	opts := struct{}{}
	// addOptions takes structs only so it can't fail
	u, _ := addOptions(fmt.Sprintf("programs/%s/structured_scopes", programID), &opts, listOpts)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	data := new([]StructuredScope)
	resp, err := s.client.Do(req, data)
	if err != nil {
		return nil, resp, err
	}

	return *data, resp, err
}

// ListAllStructuredScopes fetches a list of all structured scopes for the given program
func (s *ProgramService) ListAllStructuredScopes(programID string) ([]StructuredScope, *Response, error) {
	listOpts := &ListOptions{PageSize: defaultPageSize}
	data := []StructuredScope{}
	for {
		items, resp, err := s.ListStructuredScopes(programID, listOpts)
		if err != nil {
			return nil, resp, err
		}
		data = append(data, items...)
		if resp.Links.Next == "" {
			break
		}
		listOpts.Page = resp.Links.NextPageNumber()
	}
	return data, nil, nil
}
