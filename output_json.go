// Copyright © 2019 The Homeport Team
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

package neat

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gonvenience/bunt"
	yamlv2 "gopkg.in/yaml.v2"
	yamlv3 "gopkg.in/yaml.v3"
)

// ToJSONString marshals the provided object into JSON with text decorations
// and is basically just a convenience function to create the output processor
// and call its `ToJSON` function.
func ToJSONString(obj interface{}) (string, error) {
	return NewOutputProcessor(true, true, &DefaultColorSchema).ToCompactJSON(obj)
}

// ToJSON processes the provided input object and tries to neatly output it as
// human readable JSON honoring the preferences provided to the output processor
func (p *OutputProcessor) ToJSON(obj interface{}) (string, error) {
	var out string
	var err error

	if out, err = p.neatJSON("", obj); err != nil {
		return "", err
	}

	return out, nil
}

// ToCompactJSON processed the provided input object and tries to create a as
// compact as possible output
func (p *OutputProcessor) ToCompactJSON(obj interface{}) (string, error) {
	switch tobj := obj.(type) {
	case *yamlv3.Node:
		switch tobj.Kind {
		case yamlv3.DocumentNode:
			return p.ToCompactJSON(tobj.Content[0])

		case yamlv3.MappingNode:
			tmp := []string{}
			for i := 0; i < len(tobj.Content); i += 2 {
				k, v := tobj.Content[i], tobj.Content[i+1]

				key, err := p.ToCompactJSON(k)
				if err != nil {
					return "", err
				}

				value, err := p.ToCompactJSON(v)
				if err != nil {
					return "", err
				}

				tmp = append(tmp, fmt.Sprintf("%s: %s", key, value))
			}

			return fmt.Sprintf("{%s}", strings.Join(tmp, ", ")), nil

		case yamlv3.SequenceNode:
			tmp := []string{}
			for _, e := range tobj.Content {
				entry, err := p.ToCompactJSON(e)
				if err != nil {
					return "", err
				}

				tmp = append(tmp, entry)
			}

			return fmt.Sprintf("[%s]", strings.Join(tmp, ", ")), nil

		case yamlv3.ScalarNode:
			switch tobj.Tag {
			case "!!str":
				return fmt.Sprintf("\"%s\"", tobj.Value), nil
			}

			return tobj.Value, nil
		}

	case []interface{}:
		result := make([]string, 0)
		for _, i := range tobj {
			value, err := p.ToCompactJSON(i)
			if err != nil {
				return "", err
			}
			result = append(result, value)
		}

		return fmt.Sprintf("[%s]", strings.Join(result, ", ")), nil

	case yamlv2.MapSlice:
		result := make([]string, 0)
		for _, i := range tobj {
			value, err := p.ToCompactJSON(i)
			if err != nil {
				return "", err
			}
			result = append(result, value)
		}

		return fmt.Sprintf("{%s}", strings.Join(result, ", ")), nil

	case yamlv2.MapItem:
		key, keyError := p.ToCompactJSON(tobj.Key)
		if keyError != nil {
			return "", keyError
		}

		value, valueError := p.ToCompactJSON(tobj.Value)
		if valueError != nil {
			return "", valueError
		}

		return fmt.Sprintf("%s: %s", key, value), nil
	}

	bytes, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (p *OutputProcessor) neatJSON(prefix string, obj interface{}) (string, error) {
	switch t := obj.(type) {
	case yamlv2.MapSlice:
		if err := p.neatJSONofYAMLMapSlice(prefix, t); err != nil {
			return "", err
		}

	case []interface{}:
		if err := p.neatJSONofSlice(prefix, t); err != nil {
			return "", err
		}

	case []yamlv2.MapSlice:
		if err := p.neatJSONofSlice(prefix, p.simplify(t)); err != nil {
			return "", err
		}

	default:
		if err := p.neatJSONofScalar(prefix, obj); err != nil {
			return "", nil
		}
	}

	p.out.Flush()
	return p.data.String(), nil
}

func (p *OutputProcessor) neatJSONofYAMLMapSlice(prefix string, mapslice yamlv2.MapSlice) error {
	if len(mapslice) == 0 {
		p.out.WriteString(p.colorize("{}", "emptyStructures"))
		return nil
	}

	p.out.WriteString(bunt.Style("{", bunt.Bold()))
	p.out.WriteString("\n")

	for idx, mapitem := range mapslice {
		keyString := fmt.Sprintf("\"%v\": ", mapitem.Key)

		p.out.WriteString(prefix + p.prefixAdd())
		p.out.WriteString(p.colorize(keyString, "keyColor"))

		if p.isScalar(mapitem.Value) {
			p.neatJSONofScalar("", mapitem.Value)

		} else {
			p.neatJSON(prefix+p.prefixAdd(), mapitem.Value)
		}

		if idx < len(mapslice)-1 {
			p.out.WriteString(",")
		}

		p.out.WriteString("\n")
	}

	p.out.WriteString(prefix)
	p.out.WriteString(bunt.Style("}", bunt.Bold()))

	return nil
}

func (p *OutputProcessor) neatJSONofSlice(prefix string, list []interface{}) error {
	if len(list) == 0 {
		p.out.WriteString(p.colorize("[]", "emptyStructures"))
		return nil
	}

	p.out.WriteString(bunt.Style("[", bunt.Bold()))
	p.out.WriteString("\n")

	for idx, value := range list {
		if p.isScalar(value) {
			p.neatJSONofScalar(prefix+p.prefixAdd(), value)

		} else {
			p.out.WriteString(prefix + p.prefixAdd())
			p.neatJSON(prefix+p.prefixAdd(), value)
		}

		if idx < len(list)-1 {
			p.out.WriteString(",")
		}

		p.out.WriteString("\n")
	}

	p.out.WriteString(prefix)
	p.out.WriteString(bunt.Style("]", bunt.Bold()))

	return nil
}

func (p *OutputProcessor) neatJSONofScalar(prefix string, obj interface{}) error {
	if obj == nil {
		p.out.WriteString(p.colorize("null", "nullColor"))
		return nil
	}

	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	color := p.determineColorByType(obj)

	p.out.WriteString(prefix)
	parts := strings.Split(string(data), "\\n")
	for idx, part := range parts {
		p.out.WriteString(p.colorize(part, color))

		if idx < len(parts)-1 {
			p.out.WriteString(p.colorize("\\n", "emptyStructures"))
		}
	}

	return nil
}
