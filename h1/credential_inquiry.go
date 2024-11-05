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

// CredentialInquiry represents a credential inquiry.
//
// HackerOne API docs:  https://api.hackerone.com/customer-reference/#credential_inquiry
type CredentialInquiry struct {
	ID          *string `json:"id"`
	Type        *string `json:"type"`
	Description *string `json:"description"`
}

// Helper types for JSONUnmarshal
type credentialInquiry CredentialInquiry // Used to avoid recursion of JSONUnmarshal
type credentialInquiryUnmarshalHelper struct {
	credentialInquiry
	Attributes *credentialInquiry `json:"attributes"`
}

// UnmarshalJSON allows JSONAPI attributes and relationships to unmarshal cleanly.
func (ci *CredentialInquiry) UnmarshalJSON(b []byte) error {
	var helper credentialInquiryUnmarshalHelper
	helper.Attributes = &helper.credentialInquiry
	if err := json.Unmarshal(b, &helper); err != nil {
		return err
	}
	*ci = CredentialInquiry(helper.credentialInquiry)
	return nil
}
