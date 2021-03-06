// Copyright 2015 The Warren Authors
// Use of this source code is governed by an MIT license that can be found in
// the LICENSE file.

package text

import (
	"fmt"
	"html/template"
)

// The text/plain content type.
type Plain struct{}

// Since the display content is sanitized, this content type is safe.
func (c *Plain) Safe() bool {
	return true
}

// Sanitize the output, replace newlines with HTML line breaks, and return
// the modified content.
func (c *Plain) RenderDisplayContent(content interface{}) (string, error) {
	contentStr := template.HTMLEscapeString(content.(string))
	return fmt.Sprintf("<pre>%s</pre>", contentStr), nil
}

// Simply return the content.
func (c *Plain) RenderIndexContent(content interface{}) (string, error) {
	return content.(string), nil
}
