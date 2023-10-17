package types

import (
	"encoding/json"
	"encoding/xml"

	"github.com/pelletier/go-toml"
)

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
