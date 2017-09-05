package main

import (
	"io/ioutil"
	"path/filepath"

	"github.com/urfave/cli"

	gitlab "github.com/longkey1/go-gitlab"
	"github.com/pkg/errors"
)

func multiUploads(ctx *cli.Context, p string, dir string) ([]*gitlab.ProjectFile, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var pfiles []*gitlab.ProjectFile
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		path := filepath.Join(dir, f.Name())
		resp, err := uploads(ctx, p, path)
		if err != nil {
			return nil, err
		}
		pfiles = append(pfiles, resp)
	}

	return pfiles, nil
}

func uploads(ctx *cli.Context, p string, path string) (*gitlab.ProjectFile, error) {
	gl := gitlab.NewClient(nil, ctx.String("token"))
	uf, _, err := gl.Projects.UploadFile(p, path)
	if err != nil {
		return nil, errors.Wrap(err, "GitlabProjects UploadFile error")
	}

	return uf, nil
}
