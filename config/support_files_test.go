package config

import (
	"bytes"
	"github.com/dsalazar32/pave/helper/strparser"
	"path/filepath"
	"testing"
)

const (
	BUILDFILE = `#!/bin/bash
set -e

APP_NS=ProjectMock
EXPORT_DIR=/home/node/app/outfile
BUILD_VARS=build.variables

run_as_node() {
  su -c "$1" node
}

if [ -f $BUILD_VARS ]; then
  echo "$BUILD_VARS exist."
  while IFS='=' read -r key value; do
    eval ${key}\=${value}
  done < "$BUILD_VARS"
else
  apt-get install -qy git
  GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
  GIT_COMMIT=$(git rev-parse HEAD)
fi

if [ -d $EXPORT_DIR ]; then
  rm -rf $EXPORT_DIR
fi

cd /home/node/app
chown -R node.node /home/node/app

run_as_node "yarn install --production"
run_as_node "mkdir -p $EXPORT_DIR"
run_as_node "mkdir -p /home/node/outfile"
run_as_node "tar -czvf /home/node/outfile/${GIT_COMMIT}.tar.gz --exclude .git --exclude .idea --exclude outfile ."
run_as_node "mv /home/node/outfile/${GIT_COMMIT}.tar.gz ${EXPORT_DIR}/"

cd $EXPORT_DIR
run_as_node "shasum -a 256 ${GIT_COMMIT}.tar.gz > ${GIT_COMMIT}.tar.gz.sha256"
`

DOCKERFILE = `FROM 063112144237.dkr.ecr.us-east-1.amazonaws.com/carecloud/node_server:10.13.0

ENV APP_NS=ProjectMock
ENV NODE_DIR=/home/node
ENV APP_DIR=$NODE_DIR/$APP_NS

ARG ARTIFACT

WORKDIR /home/node

USER node
COPY outfile $NODE_DIR
RUN  shasum -a 256 -c ${ARTIFACT}.tar.gz.sha256 \
  && mkdir ${NODE_DIR}/${ARTIFACT} \
  && tar -xzvf *.tar.gz -C ${NODE_DIR}/${ARTIFACT} \
  && ln -s ${ARTIFACT} $APP_DIR \
  && ln -s $APP_DIR app \
  && rm -rf *.tar.gz

USER root
CMD ["/sbin/my_init"]
`
)

var (
	support      *Support
	supportError error
)

func TestSupportFiles_For(t *testing.T) {
	type tt struct {
		given func() string
		with  strparser.TemplatePackage
		want  string
	}

	support, supportError = LoadSupportFile(filepath.Join("..", "fixtures", "support.yml"))
	if supportError != nil {
		t.Error(supportError)
	}

	data := &Config{
		"1",
		&Pave{
			"ProjectMock",
			"node:10.13.0",
			true,
			false,
		},
	}

	tmpl := strparser.TemplatePackage{
		Ns:   "TemplateMock",
		Data: data,
		FuncMap: strparser.FuncMap{
			"baseImage": func(l string) string {
				return support.BaseImageLookup(l)
			},
		},
	}

	// NOTE: There is some tests here that the index integrity matters
	// for test comparison.  If for whatever the order is changed the tests will
	// start failing.
	tc := []tt{
		{func() string { return "{{.Pave.ProjectName}}" }, tmpl, "ProjectMock"},
		{func() string { return "{{.Pave.ProjectLang}}" }, tmpl, "node:10.13.0"},
		{func() string { return "{{.Pave.DockerEnabled}}" }, tmpl, "true"},
		{func() string { return "{{.Pave.TerraformEnabled}}" }, tmpl, "false"},
		{func() string { return "{{baseImage .Pave.ProjectLang}}" }, tmpl,
			"063112144237.dkr.ecr.us-east-1.amazonaws.com/carecloud/node_server:10.13.0"},
		{func() string {
			outs, err := supportFilesFor(data.Pave.ProjectLang, "docker")
			if err != nil {
				t.Error(err)
			}

			// 0 index
			// Returns Buildfile
			return (*outs)[0].Content
		}, tmpl, BUILDFILE,
		},
		{func() string {
			outs, err := supportFilesFor(data.Pave.ProjectLang, "docker")
			if err != nil {
				t.Error(err)
			}

			// 1 index
			// Returns Dockerfile
			return (*outs)[1].Content
		}, tmpl, DOCKERFILE,
		},
	}

	for _, c := range tc {
		b := &bytes.Buffer{}
		given := c.given()
		if err := strparser.ParseTemplate(given, c.with, b); err != nil {
			t.Error(err)
		}
		got := b.String()

		if got != c.want {
			t.Errorf("error parsing template\n%s \nwant \n\n%s \nbut got: \n\n%s\n", given, c.want, got)
		}
	}
}

func supportFilesFor(l, p string) (*Outfiles, error) {
	platform, err := support.For(l)
	if err != nil {
		return nil, err
	}

	out, err := platform.Get(p)
	if err != nil {
		return nil, err
	}

	return out, nil
}

