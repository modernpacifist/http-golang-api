package types

import (
	"encoding/json"
	"encoding/xml"

	"github.com/pelletier/go-toml"
)

// TODO: the string type must actually be a list of bytes <17-10-23, modernpacifist> //
type Data struct {
	JsonField string `json:"json"`
	XmlField  string `xml: "xml"`
	TomlField string `toml: "toml"`
}

type Marshaler interface {
	Marshal(v interface{}) ([]byte, error)
}

type JSONMarshaler struct{}

func (j *JSONMarshaler) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

type XMLMarshaler struct{}

func (x *XMLMarshaler) Marshal(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

type TOMLMarshaler struct{}

func (t *TOMLMarshaler) Marshal(v interface{}) ([]byte, error) {
	return toml.Marshal(v)
}

type User struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Salary     int    `json:"salary"`
	Occupation string `json:"occupation"`
}

type EmptyJson struct {
	Field string `json:"id"`
}
