package jig_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	. "github.com/iancmcc/jig"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Jigfile", func() {

	Describe("Parsing Jigfile data", func() {

		var (
			data string
			err  error
			jf   *Jigfile
			spec *JigSpec
		)

		Context("With one valid spec", func() {

			BeforeEach(func() {
				data = `
				{
					"ubuntu": {
						"pre": [
							"/bin/bash dependencies.sh"
						],
						"spec": [
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

			It("should have the spec by name", func() {
				Expect(jf.Specs).To(HaveKey("ubuntu"))
			})

			Context("the spec", func() {
				BeforeEach(func() {
					spec = jf.Specs["ubuntu"]
				})
				It("should have a pre script", func() {
					Expect(spec.Pre).To(Equal(
						[]string{"/bin/bash dependencies.sh"},
					))
				})
				It("should have spec commands", func() {
					Expect(spec.Build).To(Equal(
						[]string{"configure", "make", "make install"},
					))
				})
				It("should have post commands", func() {
					Expect(spec.Post).To(Equal(
						[]string{"/bin/bash cleanup.sh"},
					))
				})
				It("should have output artifacts", func() {
					Expect(spec.Output).To(Equal(
						[]string{"myfile.tgz"},
					))
				})
				It("should have an image", func() {
					Expect(spec.Image).To(Equal("jigs/test"))
				})
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

	})

	Describe("Parsing a Jigfile", func() {
		var (
			file *os.File
			err  error
			jf   *Jigfile
			spec *JigSpec
		)
		BeforeEach(func() {
			data := `
				{
					"ubuntu": {
						"pre": [
							"/bin/bash dependencies.sh"
						],
						"spec": [
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
			jf, err = ParseJigfilePath(file.Name())
		})
		AfterEach(func() {
			if file != nil {
				syscall.Unlink(file.Name())
				file = nil
			}
		})
		It("should have the spec by name", func() {
			Expect(jf.Specs).To(HaveKey("ubuntu"))
		})

		Context("the spec", func() {
			BeforeEach(func() {
				spec = jf.Specs["ubuntu"]
			})
			It("should have a pre script", func() {
				Expect(spec.Pre).To(Equal(
					[]string{"/bin/bash dependencies.sh"},
				))
			})
			It("should have spec commands", func() {
				Expect(spec.Build).To(Equal(
					[]string{"configure", "make", "make install"},
				))
			})
			It("should have post commands", func() {
				Expect(spec.Post).To(Equal(
					[]string{"/bin/bash cleanup.sh"},
				))
			})
			It("should have output artifacts", func() {
				Expect(spec.Output).To(Equal(
					[]string{"myfile.tgz"},
				))
			})
			It("should have an image", func() {
				Expect(spec.Image).To(Equal("jigs/test"))
			})
		})

		It("should not error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Parsing a Jigfile in a directory", func() {
		var (
			dirname string
			err     error
			jf      *Jigfile
			spec    *JigSpec
		)
		BeforeEach(func() {
			data := `
				{
					"ubuntu": {
						"pre": [
							"/bin/bash dependencies.sh"
						],
						"spec": [
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
			dirname, _ = ioutil.TempDir("", "jig")
			fname := filepath.Join(dirname, "Jigfile")
			ioutil.WriteFile(fname, []byte(data), 0777)
			jf, err = ParseJigfilePath(dirname)
		})
		AfterEach(func() {
			if dirname != "" {
				syscall.Unlink(dirname)
				dirname = ""
			}
		})
		It("should have the spec by name", func() {
			Expect(jf.Specs).To(HaveKey("ubuntu"))
		})

		Context("the spec", func() {
			BeforeEach(func() {
				spec = jf.Specs["ubuntu"]
			})
			It("should have a pre script", func() {
				Expect(spec.Pre).To(Equal(
					[]string{"/bin/bash dependencies.sh"},
				))
			})
			It("should have spec commands", func() {
				Expect(spec.Build).To(Equal(
					[]string{"configure", "make", "make install"},
				))
			})
			It("should have post commands", func() {
				Expect(spec.Post).To(Equal(
					[]string{"/bin/bash cleanup.sh"},
				))
			})
			It("should have output artifacts", func() {
				Expect(spec.Output).To(Equal(
					[]string{"myfile.tgz"},
				))
			})
			It("should have an image", func() {
				Expect(spec.Image).To(Equal("jigs/test"))
			})
		})

		It("should not error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

})
