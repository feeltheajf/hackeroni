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

// StructuredScope represents an asset defined by the program.
//
// HackerOne API docs: https://api.hackerone.com/reference/#structured-scope
type StructuredScope struct {
	ID                         *string    `json:"id"`
	Type                       *string    `json:"type"`
	AssetIdentifier            string     `json:"asset_identifier"`
	AssetType                  string     `json:"asset_type"`
	EligibleForBounty          bool       `json:"eligible_for_bounty"`
	EligibleForSubmission      bool       `json:"eligible_for_submission"`
	Instruction                *string    `json:"instruction"`
	ConfidentialityRequirement *string    `json:"confidentiality_requirement"`
	IntegrityRequirement       *string    `json:"integrity_requirement"`
	AvailabilityRequirement    *string    `json:"availability_requirement"`
	MaxSeverity                string     `json:"max_severity"`
	CreatedAt                  *Timestamp `json:"created_at"`
	UpdatedAt                  *Timestamp `json:"updated_at"`
	Reference                  *string    `json:"reference"`
}

// Helper types for JSONUnmarshal
type structuredScope StructuredScope // Used to avoid recursion of JSONUnmarshal
type structuredScopeUnmarshalHelper struct {
	structuredScope
	Attributes *structuredScope `json:"attributes"`
}

// UnmarshalJSON allows JSONAPI attributes and relationships to unmarshal cleanly.
func (s *StructuredScope) UnmarshalJSON(b []byte) error {
	var helper structuredScopeUnmarshalHelper
	helper.Attributes = &helper.structuredScope
	if err := json.Unmarshal(b, &helper); err != nil {
		return err
	}
	*s = StructuredScope(helper.structuredScope)
	return nil
}
