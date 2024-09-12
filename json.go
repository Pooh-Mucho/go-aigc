package aigc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Pooh-Mucho/go-aigc/internal"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

// Json schema document: https://json-schema.org/

type JsonType string

const (
	// {"type": "boolean"}
	// {"type": "integer"}
	// {"type": "number"}
	// {"type": "string"}
	// {"type": "array", "items": {"type": "string"}}
	// {"type": "object", "required": ["name", "age", "sex"], "properties": {"name": {"type": "string"}, "age": {"type": "integer"}, "

	JsonBoolean JsonType = "boolean"
	JsonInteger JsonType = "integer"
	JsonNumber  JsonType = "number"
	JsonString  JsonType = "string"
	JsonArray   JsonType = "array"
	JsonObject  JsonType = "object"
)

type JsonSchemaProperty struct {
	Name string     `json:"name,omitempty"`
	Type JsonSchema `json:"type,omitempty"`
}

type JsonSchemaProperties []JsonSchemaProperty

// JsonSchema is a JSON schema object
// ATTENTION: only implements json.Marshaler, doesn't implement json.Unmarshaler
type JsonSchema struct {
	// "boolean" | "integer" | number | "string" | "array" | "object"
	Type JsonType `json:"type,omitempty"`
	// Object description
	Description string `json:"description,omitempty"`
	// Array of possible values
	Enum []any `json:"enum,omitempty"`

	// Only for array type. E.G. {"type":"array","items":{"type": "string"}}
	Items *JsonSchema `json:"items,omitempty"`

	// Only for object type. E.G. {"type":"object","properties":{"location":{"type":"string"}}}
	Properties JsonSchemaProperties `json:"properties,omitempty"`

	// Only for object type. E.G. {"type":"object","properties":{...},"required":["location"]}
	Required []string `json:"required,omitempty"`
}

type jsonConverter struct{}

// JsonConverter is a converter for helping LLM function argument parse.
var JsonConverter = jsonConverter{}

func (p JsonSchemaProperties) MarshalJSON() ([]byte, error) {
	if len(p) == 0 {
		return []byte("{}"), nil
	}
	var err error
	var buf bytes.Buffer
	var encoder *json.Encoder

	encoder = json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)

	buf.WriteByte('{')
	for i, _ := range p {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteByte('"')
		buf.WriteString(p[i].Name)
		buf.WriteString(`":`)
		err = encoder.Encode(p[i].Type)
		if err != nil {
			return nil, err
		}
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func (s *JsonSchema) DeepCopy() JsonSchema {
	var z = JsonSchema{
		Type:        s.Type,
		Description: s.Description,
	}

	if len(s.Enum) > 0 {
		z.Enum = make([]any, len(s.Enum))
		for i, _ := range s.Enum {
			z.Enum[i] = internal.DeepCopyPrimitive(s.Enum[i])
		}
	}

	if s.Items != nil {
		z.Items = new(JsonSchema)
		*z.Items = s.Items.DeepCopy()
	}

	if len(s.Properties) > 0 {
		z.Properties = make(JsonSchemaProperties, len(s.Properties))
		for i, _ := range s.Properties {
			z.Properties[i] = JsonSchemaProperty{
				Name: s.Properties[i].Name,
				Type: s.Properties[i].Type.DeepCopy(),
			}
		}
	}

	if len(s.Required) > 0 {
		z.Required = make([]string, len(s.Required))
		copy(z.Required, s.Required)
	}

	return z
}

func EncodeJson(jsonBuffer *bytes.Buffer, value any) error {
	encoder := json.NewEncoder(jsonBuffer)
	encoder.SetEscapeHTML(false)
	// encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

func DecodeJson(jsonBuffer *bytes.Buffer, value any) error {
	decoder := json.NewDecoder(jsonBuffer)
	return decoder.Decode(value)
	// return json.Unmarshal(jsonBuffer.Bytes(), value)
}

func (c jsonConverter) ToNumber(v any) (float64, error) {
	if v == nil {
		return 0, nil
	}

	var rv = reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Float32, reflect.Float64:
		return rv.Float(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(rv.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(rv.Uint()), nil
	case reflect.Bool:
		if rv.Bool() {
			return 1, nil
		} else {
			return 0, nil
		}
	case reflect.String:
		var s = rv.String()
		s = strings.TrimSpace(s)
		if s == "" {
			return 0, nil
		}
		var f, err = strconv.ParseFloat(s, 64)
		if err == nil {
			return f, nil
		}
	}
	return 0, fmt.Errorf("[JsonConverter.ToNumber] can not convert %v (%T) to number", v, v)
}

func (c jsonConverter) FormatJsonObject(v any) (any, error) {
	if v == nil {
		return nil, nil
	}

	var rv = reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return nil, nil
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.String:
		return rv.String(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10), nil
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'f', -1, 32), nil
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'f', -1, 64), nil
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool()), nil
	}

	var timeType = reflect.TypeFor[time.Time]()

	if rv.CanConvert(timeType) {
		var t = rv.Convert(timeType).Interface().(time.Time)
		if t.IsZero() {
			return "", nil
		}
		return t.Format("2006-01-02 15:04:05 -07:00 Monday"), nil
	}

	return v, nil
}

func (c jsonConverter) FormatJsonString(v any) (string, error) {
	if v == nil {
		return "", nil
	}
	var rv = reflect.ValueOf(v)

	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return "", nil
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.String:
		return rv.String(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10), nil
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'f', -1, 32), nil
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'f', -1, 64), nil
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool()), nil
	}

	var timeType = reflect.TypeFor[time.Time]()

	if rv.CanConvert(timeType) {
		var t = rv.Convert(timeType).Interface().(time.Time)
		if t.IsZero() {
			return "", nil
		}
		return t.Format("2006-01-02 15:04:05 -07:00 Monday"), nil
	}

	var err, ok = v.(error)
	if ok {
		return err.Error(), nil
	}

	var buffer = AllocBuffer()
	defer FreeBuffer(buffer)

	err = EncodeJson(buffer, v)
	if err != nil {
		return "", fmt.Errorf("[jsonConverter.FormatJsonString] %w", err)
	}
	return buffer.String(), nil
}

func (c jsonConverter) ToArray(v any) ([]any, error) {
	var err error
	var ok bool
	var array []any

	if v == nil {
		return nil, nil
	}

	if array, ok = v.([]any); ok {
		return array, nil
	}

	var rv = reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil, nil
		}
		rv = rv.Elem()
	}

	if rv.Kind() == reflect.String {
		var s = rv.String()
		var wrapped any
		err = json.Unmarshal(unsafe.Slice(unsafe.StringData(s), len(s)), &wrapped)
		if err != nil {
			return nil, fmt.Errorf("can not convert %v (%T) to array", v, v)
		}
		v = wrapped
		rv = reflect.ValueOf(v)
	}

	if rv.Kind() != reflect.Slice {
		return nil, fmt.Errorf("can not convert %v (%T) to array", v, v)
	}
	array = make([]any, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		array[i] = rv.Index(i).Interface()
	}
	return array, nil
}

func (c jsonConverter) ToMap(v any) (map[string]any, error) {
	var ok bool
	var dict map[string]any

	if v == nil {
		return nil, nil
	}

	if dict, ok = v.(map[string]any); ok {
		return dict, nil
	}

	var rDict = reflect.ValueOf(v)
	if rDict.Kind() == reflect.Ptr {
		if rDict.IsNil() {
			return nil, nil
		}
		rDict = rDict.Elem()
	}

	if rDict.Kind() != reflect.Map {
		return nil, fmt.Errorf("can not convert %v (%T) to map", v, v)
	}
	dict = make(map[string]any, rDict.Len())
	var iter = rDict.MapRange()
	for iter.Next() {
		var k = iter.Key()
		if k.Kind() != reflect.String {
			return nil, fmt.Errorf("can not convert %v (%T) to string", v, v)
		}
		dict[k.String()] = iter.Value().Interface()
	}
	return dict, nil
}
