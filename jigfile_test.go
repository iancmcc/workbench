package jig_test

import (
	"io/ioutil"
	"os"
	"strings"
	"syscall"

	. "github.com/iancmcc/jig"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Jigfile", func() {

	Describe("Parsing Jigfile data", func() {

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
						],
						"image": "jigs/test"
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
				It("should have an image", func() {
					Expect(build.Image).To(Equal("jigs/test"))
				})
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

	})

	Describe("Parsing a Jigfile", func() {
		var (
			file  *os.File
			err   error
			jf    *Jigfile
			build Build
		)
		BeforeEach(func() {
			data := `
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
						],
						"image": "jigs/test"
					}
				}
				`
			file, _ = ioutil.TempFile("", "Jigfile")
			file.WriteString(data)
			file.Sync()
			jf, err = ParseJigfileFromFile(file.Name())
		})
		AfterEach(func() {
			if file != nil {
				syscall.Unlink(file.Name())
				file = nil
			}
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
			It("should have an image", func() {
				Expect(build.Image).To(Equal("jigs/test"))
			})
		})

		It("should not error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

})
