package jig_test

import (
	"strings"

	. "github.com/iancmcc/jig"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Jigfile", func() {

	Describe("Parsing a Jigfile", func() {

		var (
			data string
			err  error
			jf   *Jigfile
		)

		Context("With one valid builder", func() {

			BeforeEach(func() {
				data = `
				{
					"ubuntu": {
						"pre": [
							"/bin/bash dependencies.sh"
						],
						"build": [
							"configure", 
							"make", 
							"make install"
						],
						"post": [
							"/bin/bash cleanup.sh"
						],
						"output": [
							"myfile.tgz"
						]
					}
				}
				`
				jf, err = ParseJigfile(strings.NewReader(data))
			})

			It("should have the builder by name", func() {
				Expect(jf.Builds).To(HaveKey("ubuntu"))
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

	})

})
