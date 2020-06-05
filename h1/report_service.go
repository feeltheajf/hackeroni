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
	"time"
)

// ReportService handles communication with the report related methods of the H1 API.
type ReportService service

// Get fetches a Report by ID
func (s *ReportService) Get(ID string) (*Report, *Response, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("reports/%s", ID), nil)
	if err != nil {
		return nil, nil, err
	}

	rResp := new(Report)
	resp, err := s.client.Do(req, rResp)
	if err != nil {
		return nil, resp, err
	}

	return rResp, resp, err
}

// ChangeState transitions specified Report to a new state
//
// HackerOne API docs: https://api.hackerone.com/core-resources/#reports-change-state
func (s *ReportService) ChangeState(ID, message, state string, originalID *string) (*Report, *Response, error) {
	body := &StateChange{
		Message:          message,
		State:            state,
		OriginalReportID: originalID,
	}

	req, err := s.client.NewRequest("POST", fmt.Sprintf("reports/%s/state_changes", ID), body)
	if err != nil {
		return nil, nil, err
	}

	rResp := new(Report)
	resp, err := s.client.Do(req, rResp)
	if err != nil {
		return nil, resp, err
	}

	return rResp, resp, err
}

// CreateComment posts a new comment for specified Report
//
// HackerOne API docs: https://api.hackerone.com/core-resources/#reports-create-comment
func (s *ReportService) CreateComment(ID, message string, internal bool) (*Activity, *Response, error) {
	body := &CreateComment{
		Message:  message,
		Internal: internal,
	}

	req, err := s.client.NewRequest("POST", fmt.Sprintf("reports/%s/activities", ID), body)
	if err != nil {
		return nil, nil, err
	}

	rResp := new(Activity)
	resp, err := s.client.Do(req, rResp)
	if err != nil {
		return nil, resp, err
	}

	return rResp, resp, err
}

// URL returns report HTML URL
func (s *ReportService) URL(ID string) string {
	return fmt.Sprintf("%s/reports/%s", s.client.BaseURL, ID)
}

// ReportListFilter specifies optional parameters to the ReportService.List method.
//
// HackerOne API docs: https://api.hackerone.com/reference/#reports/query
type ReportListFilter struct {
	Program                           []string  `url:"program,brackets"`
	State                             []string  `url:"state,brackets,omitempty"`
	ID                                []uint64  `url:"id,brackets,omitempty"`
	CreatedAtGreaterThan              time.Time `url:"created_at__gt,omitempty"`
	CreatedAtLessThan                 time.Time `url:"created_at__lt,omitempty"`
	TriagedAtGreaterThan              time.Time `url:"triaged_at__gt,omitempty"`
	TriagedAtLessThan                 time.Time `url:"triaged_at__lt,omitempty"`
	TriagedAtNull                     bool      `url:"triaged_at__null,omitempty"`
	ClosedAtGreaterThan               time.Time `url:"closed_at__gt,omitempty"`
	ClosedAtLessThan                  time.Time `url:"closed_at__lt,omitempty"`
	ClosedAtNull                      bool      `url:"closed_at__null,omitempty"`
	DisclosedAtGreaterThan            time.Time `url:"disclosed_at__gt,omitempty"`
	DisclosedAtLessThan               time.Time `url:"disclosed_at__lt,omitempty"`
	DisclosedAtNull                   bool      `url:"disclosed_at__null,omitempty"`
	BountyAwardedAtGreaterThan        time.Time `url:"bounty_awarded_at__gt,omitempty"`
	BountyAwardedAtLessThan           time.Time `url:"bounty_awarded_at__lt,omitempty"`
	BountyAwardedAtNull               bool      `url:"bounty_awarded_at__null,omitempty"`
	SwagAtGreaterThan                 time.Time `url:"swag_at__gt,omitempty"`
	SwagAtLessThan                    time.Time `url:"swag_at__lt,omitempty"`
	SwagAtNull                        bool      `url:"swag_at__null,omitempty"`
	LastReporterActivityAtGreaterThan time.Time `url:"last_reporter_activity_at__gt,omitempty"`
	LastReporterActivityAtLessThan    time.Time `url:"last_reporter_activity_at__lt,omitempty"`
	LastReporterActivityAtNull        bool      `url:"last_reporter_activity_at__null,omitempty"`
	FirstProgramActivityAtGreaterThan time.Time `url:"first_program_activity_at__gt,omitempty"`
	FirstProgramActivityAtLessThan    time.Time `url:"first_program_activity_at__lt,omitempty"`
	FirstProgramActivityAtNull        bool      `url:"first_program_activity_at__null,omitempty"`
	LastProgramActivityAtGreaterThan  time.Time `url:"last_program_activity_at__gt,omitempty"`
	LastProgramActivityAtLessThan     time.Time `url:"last_program_activity_at__lt,omitempty"`
	LastActivityAtGreaterThan         time.Time `url:"last_activity_at__gt,omitempty"`
	LastActivityAtLessThan            time.Time `url:"last_activity_at__lt,omitempty"`
}

// List returns Reports matching the specified criteria
//
// HackerOne API docs: https://api.hackerone.com/core-resources/#reports-get-all-reports
func (s *ReportService) List(filterOpts ReportListFilter, listOpts *ListOptions) ([]Report, *Response, error) {
	opts := struct {
		Filter ReportListFilter `url:"filter,brackets"`
	}{
		Filter: filterOpts,
	}
	// addOptions takes structs only so it can't fail
	u, _ := addOptions("reports", &opts, listOpts)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	reports := new([]Report)
	resp, err := s.client.Do(req, reports)
	if err != nil {
		return nil, resp, err
	}

	return *reports, resp, err
}

// ListAll returns all Reports matching the specified criteria
//
// HackerOne API docs: https://api.hackerone.com/core-resources/#reports-get-all-reports
func (s *ReportService) ListAll(filterOpts ReportListFilter) ([]Report, *Response, error) {
	listOpts := &ListOptions{PageSize: defaultPageSize}
	reports := []Report{}
	for {
		reportList, resp, err := s.List(filterOpts, listOpts)
		if err != nil {
			return nil, resp, err
		}
		reports = append(reports, reportList...)
		if resp.Links.Next == "" {
			break
		}
		listOpts.Page = resp.Links.NextPageNumber()
	}
	return reports, nil, nil
}
