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
			data  string
			err   error
			jf    *Jigfile
			build Build
		)

		Context("With one valid build", func() {

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

			It("should have the build by name", func() {
				Expect(jf.Builds).To(HaveKey("ubuntu"))
			})

			Context("the build", func() {
				BeforeEach(func() {
					build = jf.Builds["ubuntu"]
				})
				It("should have a pre script", func() {
					Expect(build.Pre).To(Equal(
						[]string{"/bin/bash dependencies.sh"},
					))
				})
				It("should have build commands", func() {
					Expect(build.Build).To(Equal(
						[]string{"configure", "make", "make install"},
					))
				})
				It("should have post commands", func() {
					Expect(build.Post).To(Equal(
						[]string{"/bin/bash cleanup.sh"},
					))
				})
				It("should have output artifacts", func() {
					Expect(build.Output).To(Equal(
						[]string{"myfile.tgz"},
					))
				})
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

	})

})
