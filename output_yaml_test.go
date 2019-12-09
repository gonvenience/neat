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

package neat_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/gonvenience/neat"

	yamlv2 "gopkg.in/yaml.v2"
	yamlv3 "gopkg.in/yaml.v3"
)

var _ = Describe("YAML output", func() {
	Context("process input JSON for YAML output", func() {
		It("should convert JSON to YAML", func() {
			var content yamlv2.MapSlice
			if err := yamlv2.Unmarshal([]byte(`{ "name": "foobar", "list": [A, B, C] }`), &content); err != nil {
				Fail(err.Error())
			}

			result, err := ToYAMLString(content)
			Expect(err).To(BeNil())

			Expect(result).To(BeEquivalentTo(`name: foobar
list:
- A
- B
- C
`))
		})

		It("should preserve the order inside the structure", func() {
			var content yamlv2.MapSlice
			err := yamlv2.Unmarshal([]byte(`{ "list": [C, B, A], "name": "foobar" }`), &content)
			if err != nil {
				Fail(err.Error())
			}

			result, err := ToYAMLString(content)
			Expect(err).To(BeNil())

			Expect(result).To(Equal(`list:
- C
- B
- A
name: foobar
`))
		})
	})

	Context("create YAML output (go-yaml v2)", func() {
		It("should create YAML output for a simple list", func() {
			result, err := ToYAMLString([]interface{}{
				"one",
				"two",
				"three",
			})

			Expect(err).To(BeNil())
			Expect(result).To(BeEquivalentTo(`- one
- two
- three
`))
		})

		It("should create YAML output for a specific YAML v2 MapSlice list", func() {
			result, err := ToYAMLString([]yamlv2.MapSlice{
				{
					yamlv2.MapItem{
						Key:   "name",
						Value: "one",
					},
				},
				{
					yamlv2.MapItem{
						Key:   "name",
						Value: "two",
					},
				},
			})

			Expect(err).To(BeNil())
			Expect(result).To(BeEquivalentTo(`- name: one
- name: two
`))
		})

		It("should create YAML output of nested maps", func() {
			result, err := ToYAMLString(yamlv2.MapSlice{
				yamlv2.MapItem{
					Key: "map",
					Value: yamlv2.MapSlice{
						yamlv2.MapItem{
							Key: "foo",
							Value: yamlv2.MapSlice{
								yamlv2.MapItem{
									Key: "bar",
									Value: yamlv2.MapSlice{
										yamlv2.MapItem{
											Key:   "name",
											Value: "foobar",
										},
									},
								},
							},
						},
					},
				},
			})

			Expect(err).To(BeNil())
			Expect(result).To(BeEquivalentTo(`map:
  foo:
    bar:
      name: foobar
`))
		})

		It("should create YAML output for empty structures", func() {
			result, err := ToYAMLString(yamlv2.MapSlice{
				yamlv2.MapItem{
					Key:   "empty-map",
					Value: yamlv2.MapSlice{},
				},

				yamlv2.MapItem{
					Key:   "empty-list",
					Value: []interface{}{},
				},

				yamlv2.MapItem{
					Key:   "empty-scalar",
					Value: nil,
				},
			})

			Expect(err).To(BeNil())
			Expect(result).To(BeEquivalentTo(`empty-map: {}
empty-list: []
empty-scalar: null
`))
		})
	})

	Context("create YAML output (go-yaml v3)", func() {
		It("should create YAML output based on YAML node structure", func() {
			example := []byte(`---
# before document

# before map
map: # at map definition
  key: value # value
# after map

# before scalar
scalars: # at scalar definition
  boolean: true # true
  number: 42 # 42
  float: 47.11
  string: foobar
  data: !!binary Zm9vYmFyCg==
# after scaler

# before list
list: # at list definition
- one # one
- two # two
# after list

# before multiline
multiline: |
  This is
  a multi
  line te
  xt.
# after multiline

# before zeros
zeros:
  map: {}
  list: []
  scalar: nil
# after zeros

# after document
`)

			expected := `---
# before document

# before map
map:
  key: value # value
# after map
# before scalar
scalars:
  boolean: true # true
  number: 42 # 42
  float: 47.11
  string: foobar
  data: Zm9vYmFyCg==
# after scaler
# before list
list:
- one # one
- two # two
# after list
# before multiline
multiline: |
  This is
  a multi
  line te
  xt.
# after multiline
# before zeros
zeros:
  map: {}
  list: []
  scalar: nil
# after zeros
# after document
`

			var node yamlv3.Node
			Expect(yamlv3.Unmarshal(example, &node)).ToNot(HaveOccurred())

			output, err := ToYAMLString(node)
			Expect(err).ToNot(HaveOccurred())
			Expect(output).To(BeEquivalentTo(expected))
		})
	})
})
