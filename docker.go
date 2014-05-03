package jig

import (
	"bytes"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"text/template"
)

const dftemplate string = `
FROM {{.Image}}
MAINTAINER github.com/iancmcc/jig

RUN groupadd -f -g {{.Gid}} jig
RUN useradd -d /home/jig -m \
	-s /bin/bash \
	-u {{.Uid}} \
	-g {{.Gid}} \
	jig
RUN echo "jig ALL=(ALL:ALL) NOPASSWD:ALL" >> /etc/sudoers
ADD pre /tmp/pre
ADD build /tmp/build
ADD post /tmp/post
RUN /bin/bash /tmp/pre
USER jig
WORKDIR {{.Workdir}}
{{range $k, $v := .Environ}}ENV {{$k}} {{$v}}
{{end}}
`

type context struct {
	Uid      string
	Gid      string
	Image    string
	SpecName string
	Environ  map[string]string
	Workdir  string
}

func createContext(spec *JigSpec) (*context, error) {
	c := &context{
		Image: spec.Image,
	}
	u, err := user.Current()
	if err != nil {
		return nil, err
	}
	c.Uid = u.Uid
	c.Gid = u.Gid
	c.SpecName = spec.Name
	c.Environ = spec.Environ
	c.Workdir = spec.Workdir
	if c.Workdir == "" {
		c.Workdir = spec.Mount
	}
	return c, nil
}

func CreateDockerfile(spec *JigSpec) error {
	var (
		cfgdir string
		err    error
		buffer *bytes.Buffer = &bytes.Buffer{}
	)
	if cfgdir, err = spec.ConfigDir(); err != nil {
		return err
	}
	dfpath := filepath.Join(cfgdir, "Dockerfile")
	log.Debug("Creating Dockerfile at path %s", dfpath)
	tmpl, err := template.New("dockerfile").Parse(dftemplate)
	if err != nil {
		return err
	}
	ctx, err := createContext(spec)
	if err != nil {
		return err
	}
	log.Debug("Building Dockerfile using context %v", ctx)
	// TODO: Use file as io.Writer instead of buffer
	if err = tmpl.Execute(buffer, ctx); err != nil {
		return err
	}
	ioutil.WriteFile(dfpath, buffer.Bytes(), 0664)
	return nil
}
