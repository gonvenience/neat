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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/gonvenience/bunt"
	. "github.com/gonvenience/neat"

	yamlv2 "gopkg.in/yaml.v2"
	yamlv3 "gopkg.in/yaml.v3"
)

var _ = Describe("YAML output", func() {
	BeforeEach(func() {
		SetColorSettings(OFF, OFF)
	})

	AfterEach(func() {
		SetColorSettings(AUTO, AUTO)
	})

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
# start of document

# before map
map: # at map definition
  key: value # value

# before scalars
scalars: # at scalar definition
  boolean: true # true
  number: 42 # 42
  float: 47.11
  string: foobar
  data: !!binary Zm9vYmFyCg==

# before list
list: # at list definition
- one # one
- two # two

# before multiline
multiline: |-
  This is
  a multi
  line te
  xt.

# before zeros
zeros:
  map: {}
  list: []
  scalar: nil

# before anchors
anchors:
  scalar: &scalaranchor 42
  same-scalar: *scalaranchor
  list: &listanchor
  - one
  - two
  same-list: *listanchor
  empty-list: &emptylist []
  same-empty-list: *emptylist
  map: &mapanchor
    key: value
  same-map: *mapanchor
  empty-map: &emptymap {}
  same-empty-map: *emptymap

# end of document
`)

			expected := `---
# start of document

# before map
map:
  key: value # value
# before scalars
scalars:
  boolean: true # true
  number: 42 # 42
  float: 47.11
  string: foobar
  data: Zm9vYmFyCg==
# before list
list:
- one # one
- two # two
# before multiline
multiline: |
  This is
  a multi
  line te
  xt.
# before zeros
zeros:
  map: {}
  list: []
  scalar: nil
# before anchors
anchors:
  scalar: &scalaranchor 42
  same-scalar: *scalaranchor
  list: &listanchor
  - one
  - two
  same-list: *listanchor
  empty-list: &emptylist []
  same-empty-list: *emptylist
  map: &mapanchor
    key: value
  same-map: *mapanchor
  empty-map: &emptymap {}
  same-empty-map: *emptymap
# end of document
`

			var node yamlv3.Node
			Expect(yamlv3.Unmarshal(example, &node)).ToNot(HaveOccurred())

			output, err := ToYAMLString(node)
			Expect(err).ToNot(HaveOccurred())
			Expect(output).To(BeEquivalentTo(expected))
		})

		It("should create YAML output of so called named entry lists", func() {
			example := []byte(`---
yaml:
  named-entry-list-using-name:
  - name: A
  - name: B
  - name: C
  - name: X
  - name: Z

  named-entry-list-using-key:
  - key: A
  - key: B
  - key: C
  - key: X
  - key: Z

  named-entry-list-using-id:
  - id: A
  - id: B
  - id: C
  - id: X
  - id: Z
`)

			expected := `---
yaml:
  named-entry-list-using-name:
  - name: A
  - name: B
  - name: C
  - name: X
  - name: Z
  named-entry-list-using-key:
  - key: A
  - key: B
  - key: C
  - key: X
  - key: Z
  named-entry-list-using-id:
  - id: A
  - id: B
  - id: C
  - id: X
  - id: Z
`

			var node yamlv3.Node
			Expect(yamlv3.Unmarshal(example, &node)).ToNot(HaveOccurred())

			output, err := ToYAMLString(node)
			Expect(err).ToNot(HaveOccurred())
			Expect(output).To(BeEquivalentTo(expected))
		})

		It("should create YAML output of multi-line text", func() {
			example := []byte(`---
data:
  repos.yaml: |-
    repos:
    - apply_requirements:
      - approved
      - mergeable
      id: /.*/
      test: /.*/
`)

			expected := `---
data:
  repos.yaml: |
    repos:
    - apply_requirements:
      - approved
      - mergeable
      id: /.*/
      test: /.*/
`

			var node yamlv3.Node
			Expect(yamlv3.Unmarshal(example, &node)).ToNot(HaveOccurred())

			output, err := ToYAMLString(node)
			Expect(err).ToNot(HaveOccurred())
			Expect(output).To(BeEquivalentTo(expected))
		})

		It("should quote reserved keywords", func() {
			example := []byte(`---
data:
  foo: "true"
  bar: true
  and: null
  but: "null"
  dash: "-"
`)

			expected := `---
data:
  foo: "true"
  bar: true
  and: null
  but: "null"
  dash: "-"
`

			var node yamlv3.Node
			Expect(yamlv3.Unmarshal(example, &node)).ToNot(HaveOccurred())

			output, err := ToYAMLString(node)
			Expect(err).ToNot(HaveOccurred())
			Expect(output).To(BeEquivalentTo(expected))
		})
	})

	Context("create YAML output for type struct", func() {
		type Dependency struct {
			Name    string `yaml:"name"`
			Version string `yaml:"version"`
			Active  bool   `yaml:"active"`
		}

		type Example struct {
			Name         string       `yaml:"name"`
			Version      string       `yaml:"version"`
			Dependencies []Dependency `yaml:"dependencies"`
		}

		var yml = func(input string) yamlv3.Node {
			var result yamlv3.Node
			if err := yamlv3.Unmarshal([]byte(input), &result); err != nil {
				Fail(err.Error())
			}

			return result
		}

		It("should output a generic type struct", func() {
			SetColorSettings(ON, ON)
			defer SetColorSettings(AUTO, AUTO)

			expected, _ := ToYAMLString(yml(`---
name: foobar
version: v1.0.0
dependencies:
- name: foo
  version: v0.5.0
  active: true
- name: bar
  version: v0.5.0
  active: true
`))

			output, err := ToYAMLString(Example{
				Name:    "foobar",
				Version: "v1.0.0",
				Dependencies: []Dependency{
					{
						Name:    "foo",
						Version: "v0.5.0",
						Active:  true,
					},
					{
						Name:    "bar",
						Version: "v0.5.0",
						Active:  true,
					},
				},
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(output).To(Equal(expected))
		})

		It("should output a pointer to a generic type struct", func() {
			SetColorSettings(ON, ON)
			defer SetColorSettings(AUTO, AUTO)

			expected, _ := ToYAMLString(yml(`---
name: foobar
version: v1.0.0
dependencies:
- name: foo
  version: v0.5.0
  active: true
- name: bar
  version: v0.5.0
  active: true
`))

			output, err := ToYAMLString(&Example{
				Name:    "foobar",
				Version: "v1.0.0",
				Dependencies: []Dependency{
					{
						Name:    "foo",
						Version: "v0.5.0",
						Active:  true,
					},
					{
						Name:    "bar",
						Version: "v0.5.0",
						Active:  true,
					},
				},
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(output).To(Equal(expected))
		})
	})
})
