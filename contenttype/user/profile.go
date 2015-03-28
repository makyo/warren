package user

import (
	"bytes"
	"html/template"

	"gopkg.in/mgo.v2/bson"

	"github.com/warren-community/warren/contenttype/text"
)

type Profile struct {
	ProfileText string
	Website     string
}

func NewProfile(in bson.M) Profile {
	return Profile{
		ProfileText: in["profiletext"].(string),
		Website:     in["website"].(string),
	}
}

// Since users are managed through markdown, they are a safe content type.
func (c *Profile) Safe() bool {
	return true
}

// Render the profile using markdown
// TODO Users may need additional fields in the future.
func (c *Profile) RenderDisplayContent(content interface{}) (string, error) {
	profile := content.(Profile)
	profileText := template.HTML(text.RenderMarkdown(profile.ProfileText))
	buf := new(bytes.Buffer)
	tmpl := template.Must(template.ParseFiles("templates/contenttype/user/profile.tmpl"))
	tmpl.Execute(buf, map[string]interface{}{
		"ProfileText": profileText,
		"Website":     profile.Website,
	})
	return buf.String(), nil
}

// Simply return the markdown content.
func (c *Profile) RenderIndexContent(content interface{}) (string, error) {
	return (content.(Profile)).ProfileText, nil
}
