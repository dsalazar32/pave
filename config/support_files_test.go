package config

import (
	"github.com/davecgh/go-spew/spew"
	"path/filepath"
	"testing"
)

type outputs map[string]string

func TestSupportFiles_For(t *testing.T) {
	type tt struct {
		c    Pave
		want map[string]string
	}

	tc := []tt{
		{Pave{"Pave", "node:10.13.0", true, false}, outputs{
			"Dockerfile": dockerfile,
			"Buildfile": buildfile,
		}},
	}
	_ = tc

	s, err := LoadSupportFile(filepath.Join("..", "fixtures", "support.yml"))
	if err != nil {
		t.Fatal(err)
	}

	p, err := s.SupportFiles.For("node:10.13.0")
	if err != nil {
		t.Error(err)
	}

	spew.Dump(p.Get("docker"))
}

var buildfile = `
#!/bin/bash
set -e

APP_NS={{.projectName}}
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

var dockerfile = `
FROM {{.baseImage}}

ENV APP_NS={{.projectName}}
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
