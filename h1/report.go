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

// ReportState represent possible states for a report
//
// HackerOne API docs: https://api.hackerone.com/reference/#report
const (
	ReportStateNew           string = "new"
	ReportStateTriaged       string = "triaged"
	ReportStateNeedsMoreInfo string = "needs-more-info"
	ReportStateResolved      string = "resolved"
	ReportStateNotApplicable string = "not-applicable"
	ReportStateInformative   string = "informative"
	ReportStateDuplicate     string = "duplicate"
	ReportStateSpam          string = "spam"
)

// Report represents a report.
//
// HackerOne API docs: https://api.hackerone.com/reference/#report
type Report struct {
	ID                            *string          `json:"id"`
	Type                          *string          `json:"type"`
	Title                         *string          `json:"title"`
	VulnerabilityInformation      *string          `json:"vulnerability_information,omitempty"`
	State                         *string          `json:"state"`
	CreatedAt                     *Timestamp       `json:"created_at"`
	TriagedAt                     *Timestamp       `json:"triaged_at,omitempty"`
	ClosedAt                      *Timestamp       `json:"closed_at,omitempty"`
	LastReporterActivityAt        *Timestamp       `json:"last_reporter_activity_at,omitempty"`
	FirstProgramActivityAt        *Timestamp       `json:"first_program_activity_at,omitempty"`
	LastProgramActivityAt         *Timestamp       `json:"last_program_activity_at,omitempty"`
	LastActivityAt                *Timestamp       `json:"last_activity_at,omitempty"`
	LastPublicActivityAt          *Timestamp       `json:"last_public_activity_at,omitempty"`
	BountyAwardedAt               *Timestamp       `json:"bounty_awarded_at,omitempty"`
	SwagAwardedAt                 *Timestamp       `json:"swag_awarded_at,omitempty"`
	DisclosedAt                   *Timestamp       `json:"disclosed_at,omitempty"`
	ReporterAgreedOnGoingPublicAt *Timestamp       `json:"reporter_agreed_on_going_public_at,omitempty"`
	IssueTrackerReferenceID       *string          `json:"issue_tracker_reference_id,omitempty"`
	IssueTrackerReferenceURL      *string          `json:"issue_tracker_reference_url,omitempty"`
	CVEIDs                        []string         `json:"cve_ids,omitempty"`
	Source                        *string          `json:"source,omitempty"`
	Program                       *Program         `json:"program"`
	RawAssignee                   json.RawMessage  `json:"assignee,omitempty"` // Used by Assignee()
	Attachments                   []Attachment     `json:"attachments,omitempty"`
	Swag                          []Swag           `json:"swag,omitempty"`
	Weakness                      *Weakness        `json:"weakness,omitempty"`
	StructuredScope               *StructuredScope `json:"structured_scope,omitempty"`
	Severity                      *Severity        `json:"severity,omitempty"`
	Reporter                      *User            `json:"reporter,omitempty"`
	// TriggeredPreSubmissionTrigger
	Activities []Activity      `json:"activities,omitempty"`
	Bounties   []Bounty        `json:"bounties,omitempty"`
	Summaries  []ReportSummary `json:"summaries,omitempty"`
	// CustomFieldValues
}

// Helper types for JSONUnmarshal
type report Report // Used to avoid recursion of JSONUnmarshal
type reportUnmarshalHelper struct {
	report
	Attributes    *report `json:"attributes"`
	Relationships struct {
		Program struct {
			Data *Program `json:"data"`
		} `json:"program"`
		RawAssignee struct {
			Data json.RawMessage `json:"data"`
		} `json:"assignee"`
		Attachments struct {
			Data []Attachment `json:"data"`
		} `json:"attachments"`
		Swag struct {
			Data []Swag `json:"data"`
		} `json:"swag"`
		Weakness struct {
			Data *Weakness `json:"data"`
		} `json:"weakness"`
		Severity struct {
			Data *Severity `json:"data"`
		} `json:"severity"`
		Reporter struct {
			Data *User `json:"data"`
		} `json:"reporter"`
		Activities struct {
			Data []Activity `json:"data"`
		} `json:"activities"`
		Bounties struct {
			Data []Bounty `json:"data"`
		} `json:"bounties"`
		Summaries struct {
			Data []ReportSummary `json:"data"`
		} `json:"summaries"`
		StructuredScope struct {
			Data *StructuredScope `json:"data"`
		} `json:"structured_scope"`
	} `json:"relationships"`
}

// UnmarshalJSON allows JSONAPI attributes and relationships to unmarshal cleanly.
func (r *Report) UnmarshalJSON(b []byte) error {
	var helper reportUnmarshalHelper
	helper.Attributes = &helper.report
	if err := json.Unmarshal(b, &helper); err != nil {
		return err
	}
	*r = Report(helper.report)
	r.Program = helper.Relationships.Program.Data
	r.RawAssignee = helper.Relationships.RawAssignee.Data
	r.Attachments = helper.Relationships.Attachments.Data
	r.Swag = helper.Relationships.Swag.Data
	r.Weakness = helper.Relationships.Weakness.Data
	r.Severity = helper.Relationships.Severity.Data
	r.Reporter = helper.Relationships.Reporter.Data
	r.Activities = helper.Relationships.Activities.Data
	for idx := range r.Activities {
		r.Activities[idx].report = r
	}
	r.Bounties = helper.Relationships.Bounties.Data
	r.Summaries = helper.Relationships.Summaries.Data
	r.StructuredScope = helper.Relationships.StructuredScope.Data
	return nil
}

// Assignee returns returns the parsed assignee. For recognized assignee types, a value of the corresponding struct type will be returned.
func (r *Report) Assignee() (assignee interface{}) {
	var obj unknownResource
	if err := json.Unmarshal(r.RawAssignee, &obj); err != nil {
		panic(err.Error())
	}
	if obj.Type == nil {
		return nil
	}
	switch *obj.Type {
	case UserType:
		assignee = &User{}
	case GroupType:
		assignee = &Group{}
	}
	if err := json.Unmarshal(r.RawAssignee, &assignee); err != nil {
		panic(err.Error())
	}
	return assignee
}

// Helper function for Participants
func appendUserIfMissing(slice []User, u User) []User {
	for _, ele := range slice {
		if *ele.ID == *u.ID {
			return slice
		}
	}
	return append(slice, u)
}

// Participants returns a list of participants in the report. It does not include the reporter
func (r *Report) Participants(internal bool) (participants []User) {
	// Loop each activity
	for _, activity := range r.Activities {
		// Skip internal if we are not interested in them
		if *activity.Internal && !internal {
			continue
		}
		// Get the actor (if it's not a user skip it)
		user, success := activity.Actor().(*User)
		if !success {
			continue
		}
		// Skip the reporter
		if *r.Reporter.ID == *user.ID {
			continue
		}
		// Add to known participants list
		participants = appendUserIfMissing(participants, *user)
	}
	// Return the results
	return participants
}
