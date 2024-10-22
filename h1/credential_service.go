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
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// CredentialService handles communication with the credential related methods of the H1 API.
type CredentialService service

func (s *CredentialService) ListCredentialInquiries(programID string, listOpts *ListOptions) ([]CredentialInquiry, *Response, error) {
	opts := struct{}{}
	// addOptions takes structs only so it can't fail
	u, _ := addOptions(fmt.Sprintf("programs/%s/credential_inquiries", programID), &opts, listOpts)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	inquiries := new([]CredentialInquiry)
	resp, err := s.client.Do(req, inquiries)
	if err != nil {
		return nil, resp, err
	}

	return *inquiries, resp, err
}

func (s *CredentialService) ListAllCredentialInquiries(programID string) ([]CredentialInquiry, *Response, error) {
	listOpts := &ListOptions{PageSize: defaultPageSize}
	data := []CredentialInquiry{}
	for {
		items, resp, err := s.ListCredentialInquiries(programID, listOpts)
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

func (s *CredentialService) ListCredentialInquiryResponses(programID, inquiryID string, listOpts *ListOptions) ([]CredentialInquiryResponse, *Response, error) {
	opts := struct{}{}
	// addOptions takes structs only so it can't fail
	u, _ := addOptions(fmt.Sprintf("programs/%s/credential_inquiries/%s/credential_inquiry_responses", programID, inquiryID), &opts, listOpts)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	responses := new([]CredentialInquiryResponse)
	resp, err := s.client.Do(req, responses)
	if err != nil {
		return nil, resp, err
	}

	return *responses, resp, err
}

func (s *CredentialService) ListAllCredentialInquiryResponses(programID, inquiryID string) ([]CredentialInquiryResponse, *Response, error) {
	listOpts := &ListOptions{PageSize: defaultPageSize}
	data := []CredentialInquiryResponse{}
	for {
		items, resp, err := s.ListCredentialInquiryResponses(programID, inquiryID, listOpts)
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

// CreateCredential creates a new credential for specified StructuredScope
//
// HackerOne API docs: https://api.hackerone.com/customer-resources/#credentials-create-a-credential
func (s *CredentialService) CreateCredential(structuredScopeID string, credentials interface{}, assignee string) (*Credential, *Response, error) {
	b, err := json.Marshal(credentials)
	if err != nil {
		return nil, nil, err
	}
	credential := &CreateCredential{
		StructuredScopeID: structuredScopeID,
		Data: &CreateCredentialData{
			Type: CredentialType,
			Attributes: &CreateCredentialAttributes{
				Credentials: string(b),
				Assignee:    assignee,
			},
		},
	}

	// NOTE cannot use s.client.NewRequest here because H1 does not follow JSON API spec
	// (structured scope id is outside of data)
	rel, err := url.Parse("credentials")
	if err != nil {
		return nil, nil, err
	}
	body := bytes.NewBuffer(nil)
	if err := json.NewEncoder(body).Encode(credential); err != nil {
		return nil, nil, err
	}
	req, err := http.NewRequest("POST", s.client.BaseURL.ResolveReference(rel).String(), body)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("User-Agent", s.client.UserAgent)
	req.Header.Add("Content-Type", "application/json")

	data := new(Credential)
	resp, err := s.client.Do(req, data)
	if err != nil {
		return nil, resp, err
	}

	return data, resp, err
}

func (s *CredentialService) DeleteCredential(ID string) (*Response, error) {
	req, err := s.client.NewRequest("DELETE", fmt.Sprintf("credentials/%s", ID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}
