package jig

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("jig")
)

func Build(jf *Jigfile) {
	log.Debug("Attempting to build using Jigfile %s",
		filepath.Join(jf.Path, "Jigfile"))
	for _, spec := range jf.Specs {
		log.Info(`Executing spec "%s"`, spec.Name)
		if err := execute(spec); err != nil {
			log.Critical("%v", err)
		}
	}
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
	containername := randstring(12)
	buildcmd := exec.Command("docker", "build", "-t", containername, cfg)
	// TODO: Enable if verbose
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
		fmt.Sprintf("%s:%s", rootpath, jb.Mount),
		containername,
		"/bin/bash", "-c", "sudo /bin/bash /tmp/pre && /bin/bash /tmp/build && sudo /bin/bash /tmp/post")
	// TODO: Enable if verbose
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr
	log.Info("%s: Kicking a build", jb.Name)
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		err = fmt.Errorf("Build did not complete successfully: %s\nStderr: %s",
			err, stderr.String())
		return err
	}
	log.Info("Build complete!")
	return nil
}
