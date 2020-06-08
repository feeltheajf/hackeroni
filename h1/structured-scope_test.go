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
	"github.com/stretchr/testify/assert"

	"testing"
)

// TODO

func Test_StructuredScope(t *testing.T) {
	var actual ReportSummary
	loadResource(t, &actual, "tests/resources/report-summary.json")
	expected := ReportSummary{
		ID:        String("1337"),
		Type:      String(ReportSummaryType),
		Content:   String("There was a cross-site scripting vulnerability in our login form."),
		Category:  String("team"),
		CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		UpdatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		User: &User{
			ID:       String("1337"),
			Type:     String(UserType),
			Disabled: Bool(false),
			Username: String("api-example"),
			Name:     String("API Example"),
			ProfilePicture: UserProfilePicture{
				Size62x62:   String("/assets/avatars/default.png"),
				Size82x82:   String("/assets/avatars/default.png"),
				Size110x110: String("/assets/avatars/default.png"),
				Size260x260: String("/assets/avatars/default.png"),
			},
			CreatedAt: NewTimestamp("2016-02-02T04:05:06.000Z"),
		},
	}
	assert.Equal(t, expected, actual)
}
