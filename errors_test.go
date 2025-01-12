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
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/pkg/errors"

	. "github.com/gonvenience/bunt"
	. "github.com/gonvenience/neat"
)

var _ = Describe("error rendering", func() {
	BeforeEach(func() {
		SetColorSettings(ON, ON)
	})

	AfterEach(func() {
		SetColorSettings(AUTO, AUTO)
	})

	Context("rendering errors", func() {
		It("should render simple errors", func() {
			Expect(SprintError(fmt.Errorf("failed to load"))).To(
				BeEquivalentTo(ContentBox(
					"Error",
					"failed to load",
					HeadlineColor(OrangeRed),
					ContentColor(Red),
				)))
		})

		It("should render a context error using a box", func() {
			cause := fmt.Errorf("failed to load X and Y")
			err := fmt.Errorf("unable to start Z: %w", cause)

			Expect(SprintError(err)).To(
				BeEquivalentTo(ContentBox(
					"Error: unable to start Z",
					cause.Error(),
					HeadlineColor(OrangeRed),
					ContentColor(Red),
				)))
		})

		It("should render to a writer", func() {
			buf := bytes.Buffer{}
			out := bufio.NewWriter(&buf)
			FprintError(out, fmt.Errorf("foo: %w", fmt.Errorf("failed to do X")))
			out.Flush()

			Expect(buf.String()).To(
				BeEquivalentTo(ContentBox(
					"Error: foo",
					"failed to do X",
					HeadlineColor(OrangeRed),
					ContentColor(Red),
				)))
		})

		It("should render to stdout, too", func() {
			captureStdout := func(f func()) string {
				r, w, err := os.Pipe()
				Expect(err).ToNot(HaveOccurred())

				tmp := os.Stdout
				defer func() {
					os.Stdout = tmp
				}()

				os.Stdout = w
				f()
				w.Close()

				var buf bytes.Buffer
				_, err = io.Copy(&buf, r)
				Expect(err).ToNot(HaveOccurred())

				return buf.String()
			}

			Expect(captureStdout(func() { PrintError(fmt.Errorf("foo: %w", fmt.Errorf("failed to do X"))) })).To(
				BeEquivalentTo(ContentBox(
					"Error: foo",
					"failed to do X",
					HeadlineColor(OrangeRed),
					ContentColor(Red),
				)))
		})

		It("should render a context error inside a context error", func() {
			root := fmt.Errorf("unable to load X")
			cause := fmt.Errorf("failed to start Y: %w", root)
			context := "cannot process Z"

			err := fmt.Errorf("cannot process Z: %w", cause)
			Expect(SprintError(err)).To(
				BeEquivalentTo(ContentBox(
					"Error: "+context,
					SprintError(cause),
					HeadlineColor(OrangeRed),
					ContentColor(Red),
				)))
		})

		It("should render github.com/pkg/errors package errors", func() {
			message := "unable to start Z"
			cause := fmt.Errorf("failed to load X and Y")
			err := errors.Wrap(cause, message)

			Expect(SprintError(err)).To(
				BeEquivalentTo(ContentBox(
					"Error: "+message,
					cause.Error(),
					HeadlineColor(OrangeRed),
					ContentColor(Red),
				)))
		})
	})
})
