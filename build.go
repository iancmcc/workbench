package jig

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("jig")
)

func Build(jf *Jigfile) {
	var wg sync.WaitGroup
	log.Debug("Attempting to build using Jigfile %s",
		filepath.Join(jf.Path, "Jigfile"))
	for _, spec := range jf.Specs {
		wg.Add(1)
		log = logging.MustGetLogger(spec.Name)
		log.Info(`Executing spec "%s"`, spec.Name)
		go func() {
			defer wg.Done()
			if err := execute(spec); err != nil {
				log.Critical("%v", err)
			}
		}()
	}
	wg.Wait()
}

func execute(jb *JigSpec) error {
	var (
		stderr bytes.Buffer
	)
	if err := CreateDockerfile(jb); err != nil {
		return err
	}
	rootpath := jb.Jigfile.Path
	cfg, err := jb.ConfigDir()
	if err != nil {
		return err
	}
	log.Info("%s: Building the base container", jb.Name)
	jb.CreateScripts()
	containername := fmt.Sprintf("jig-%s", jb.Name)
	buildcmd := exec.Command("docker", "build", "-t", containername, cfg)
	//buildcmd.Stdout = os.Stdout
	buildcmd.Stderr = &stderr
	if err := buildcmd.Start(); err != nil {
		return err
	}
	if err := buildcmd.Wait(); err != nil {
		err = fmt.Errorf("Error building jig container: %s\nStderr: %s",
			err, stderr.String())
		return err
	}
	cmd := exec.Command("docker", "run", "--rm", "-v",
		fmt.Sprintf("%s:/mnt/jig", rootpath),
		containername,
		"/bin/bash", "/tmp/build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr
	log.Info("%s: Kicking a build", jb.Name)
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		err = fmt.Errorf("Error exporting: %s\nStderr: %s",
			err, stderr.String())
		return err
	}
	return nil
}
