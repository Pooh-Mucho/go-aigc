package main

import (
	_ "github.com/Pooh-Mucho/go-aigc"
	_ "github.com/Pooh-Mucho/go-aigc/chat"
	"unsafe"
)

import (
	"bytes"
	"encoding/json"
	"github.com/Pooh-Mucho/go-aigc"
	"strings"
)

type AnthropicTool struct {
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	InputSchema *aigc.JsonSchema `json:"input_schema,omitempty"`
}

func test1() {
	var err error
	var buf []byte
	var t1 AnthropicTool

	t1.Name = "get_weather"
	t1.Description = "Get the current weather in a given location"
	t1.InputSchema = new(aigc.JsonSchema)

	t1.InputSchema.Type = aigc.JsonObject
	t1.InputSchema.Properties = append(t1.InputSchema.Properties, aigc.JsonSchemaProperty{
		Name: "location",
		Type: aigc.JsonSchema{
			Type:        aigc.JsonString,
			Description: "The city and state, e.g. San Francisco, CA",
		},
	})
	t1.InputSchema.Properties = append(t1.InputSchema.Properties, aigc.JsonSchemaProperty{
		Name: "unit",
		Type: aigc.JsonSchema{
			Type:        aigc.JsonString,
			Description: "The unit of temperature, either 'celsius' or 'fahrenheit'",
			Enum:        []any{"celsius", "fahrenheit"},
		},
	})

	buf, err = json.Marshal(t1)
	if err != nil {
		panic(err)
	}
	println(string(buf))
}

func test2() {
	var err error
	var buf []byte
	var data struct {
		Buffer  *bytes.Buffer    `json:"buffer,omitempty"`
		Builder *strings.Builder `json:"builder,omitempty"`
	}

	data.Buffer = &bytes.Buffer{}
	data.Buffer.Grow(100)
	data.Builder = &strings.Builder{}
	data.Builder.Grow(100)

	// data.buffer.WriteByte(0)
	data.Buffer.WriteString("hello buffer")
	data.Builder.WriteString("hello builder")

	buf, err = json.MarshalIndent(data, "", "    ")
	if err != nil {
		panic(err)
	}
	println(string(buf))
}

func test3() {
	type Base struct {
		A int `json:"a"`
		B int `json:"b,omitempty"`
	}

	type Sub1 struct {
		C int `json:"c"`
		B int `json:"b"`
		Base
	}

	type Sub2 struct {
		D int `json:"d"`
		Base
	}

	s1 := Sub1{C: 0, Base: Base{A: 1, B: 2}}

	buf, err := json.MarshalIndent(s1, "", "    ")
	if err != nil {
		panic(err)
	}
	println(string(buf))
}

func test4() {
	type Z struct {
		z1 string `json:"z1,omitempty"`
		z2 string `json:"z2,omitempty"`
	}

	type AA struct {
		XX string `json:"xx,omitempty"`
		YY struct {
			Y1 string `json:"y1,omitempty"`
			Y2 string `json:"y2,omitempty"`
		} `json:"yy,omitempty"`
		ZZ any `json:"zzs,omitempty"`
	}

	var a AA
	a.XX = "xx"

	buf, err := json.MarshalIndent(a, "", "    ")
	if err != nil {
		panic(err)
	}
	println(string(buf))
}

func test5() {
	type X struct {
		Y      *bool           `json:"y,omitempty"`
		Z      *map[string]any `json:"z,omitempty"`
		yValue bool
		zValue map[string]any
	}

	var x X
	x.Y = &x.yValue
	x.zValue = make(map[string]any)
	x.Z = &x.zValue
	// x.Z = make(map[string]any, 1)
	buf, err := json.MarshalIndent(x, "", "    ")
	if err != nil {
		panic(err)
	}
	println(string(buf))
}

type vector *[1]float32

func toVector(v []float32) vector {
	var p = unsafe.Pointer(unsafe.SliceData(v))
	return vector(p)
}

func main() {
	// test5()
	var x = []float32{99: 0.0}
	var v = toVector(x)
	for i := 0; i < 10; i++ {
		v[i] = float32(i)
	}
	println(x)
}
