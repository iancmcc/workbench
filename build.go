package jig

import (
	"errors"
	"fmt"
	"sync"

	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("jig")
)

func Build(jf *Jigfile) {
	var wg sync.WaitGroup
	log.Debug(fmt.Sprintf("Attempting to build using Jigfile %s", jf.Path))
	for _, spec := range jf.Specs {
		wg.Add(1)
		log = logging.MustGetLogger(spec.Name)
		log.Info(fmt.Sprintf(`Executing spec "%s"`, spec.Name))
		go func() {
			defer wg.Done()
			if err := execute(spec); err != nil {
				log.Critical(fmt.Sprintf("%v", err))
			}
		}()
	}
	wg.Wait()
}

func execute(jb *JigSpec) error {
	// First create the Dockerfile
	// Then create the makefile
	return errors.New("Shit got bad")
}
