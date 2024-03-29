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

var _ = Describe("JSON output", func() {
	Context("create JSON output", func() {
		BeforeEach(func() {
			SetColorSettings(OFF, OFF)
		})

		AfterEach(func() {
			SetColorSettings(AUTO, AUTO)
		})

		It("should create JSON output for a simple list", func() {
			result, err := ToJSONString([]interface{}{
				"one",
				"two",
				"three",
			})

			Expect(err).To(BeNil())
			Expect(result).To(BeEquivalentTo(`["one", "two", "three"]`))
		})

		It("should create JSON output of nested maps", func() {
			example := yamlv2.MapSlice{
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
			}

			var result string
			var err error
			var outputProcessor = NewOutputProcessor(false, false, &DefaultColorSchema)

			result, err = outputProcessor.ToCompactJSON(example)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(BeEquivalentTo(`{"map": {"foo": {"bar": {"name": "foobar"}}}}`))

			result, err = outputProcessor.ToJSON(example)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(BeEquivalentTo(`{
  "map": {
    "foo": {
      "bar": {
        "name": "foobar"
      }
    }
  }
}`))
		})

		It("should create JSON output for empty structures", func() {
			example := yamlv2.MapSlice{
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
			}

			var result string
			var err error
			var outputProcessor = NewOutputProcessor(false, false, &DefaultColorSchema)

			result, err = outputProcessor.ToCompactJSON(example)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(BeEquivalentTo(`{"empty-map": {}, "empty-list": [], "empty-scalar": null}`))

			result, err = outputProcessor.ToJSON(example)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(BeEquivalentTo(`{
  "empty-map": {},
  "empty-list": [],
  "empty-scalar": null
}`))
		})

		It("should create compact JSON with correct quotes for different types", func() {
			result, err := ToJSONString([]interface{}{
				"string",
				42.0,
				42,
				true,
			})

			Expect(err).To(BeNil())
			Expect(result).To(BeEquivalentTo(`["string", 42, 42, true]`))
		})

		It("should create JSON with correct quotes for different types", func() {
			result, err := NewOutputProcessor(true, true, &DefaultColorSchema).ToJSON([]interface{}{
				"string",
				42.0,
				42,
				true,
			})

			Expect(err).To(BeNil())
			Expect(result).To(BeEquivalentTo(`[
  "string",
  42,
  42,
  true
]`))
		})

		It("should not create JSON output with unquoted timestamps (https://github.com/gonvenience/neat/issues/69)", func() {
			example := func() *yamlv3.Node {
				var node yamlv3.Node
				if err := yamlv3.Unmarshal([]byte(`timestamp: 2021-08-21T00:00:00Z`), &node); err != nil {
					Fail(err.Error())
				}

				return &node
			}()

			result, err := NewOutputProcessor(false, false, &DefaultColorSchema).ToCompactJSON(example)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(`{"timestamp": "2021-08-21T00:00:00Z"}`))

			result, err = NewOutputProcessor(false, false, &DefaultColorSchema).ToJSON(example)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(`{
  "timestamp": "2021-08-21T00:00:00Z"
}`))
		})

		It("should parse all YAML spec conform timestamps", func() {
			var example yamlv3.Node
			Expect(yamlv3.Unmarshal([]byte(`timestamp: 2033-12-20`), &example)).To(BeNil())

			result, err := NewOutputProcessor(false, false, &DefaultColorSchema).ToCompactJSON(example)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(`{"timestamp": "2033-12-20T00:00:00Z"}`))

			result, err = NewOutputProcessor(false, false, &DefaultColorSchema).ToJSON(example)
			Expect(err).ToNot(HaveOccurred())
			Expect(result).To(Equal(`{
  "timestamp": "2033-12-20T00:00:00Z"
}`))
		})
	})

	Context("create JSON output with colors", func() {
		BeforeEach(func() {
			SetColorSettings(ON, ON)
		})

		AfterEach(func() {
			SetColorSettings(AUTO, AUTO)
		})

		It("should create empty list output", func() {
			result, err := NewOutputProcessor(true, true, &DefaultColorSchema).ToJSON([]interface{}{})
			Expect(err).To(BeNil())
			Expect(result).To(BeEquivalentTo(Sprint("PaleGoldenrod{[]}")))

			result, err = NewOutputProcessor(true, true, &DefaultColorSchema).ToJSON(yml("[]"))
			Expect(err).To(BeNil())
			Expect(result).To(BeEquivalentTo(Sprint("PaleGoldenrod{[]}")))
		})
	})
})
