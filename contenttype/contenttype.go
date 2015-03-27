// Copyright 2015 The Warren Authors
// Use of this source code is governed by an MIT license that can be found in
// the LICENSE file.

package contenttype

// A content type encodes the data in an entity for display on the site.
type ContentType interface {
	RenderDisplayContent(content string) (string, error)
	RenderIndexContent(content string) (string, error)
	Safe() bool
}
