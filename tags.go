package main

import (
	gitlab "github.com/longkey1/go-gitlab"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func getTag(ctx *cli.Context, pid string, t string) (*gitlab.Tag, error) {
	gl := gitlab.NewClient(nil, ctx.String("token"))
	tag, resp, err := gl.Tags.GetTag(pid, t)
	if resp.StatusCode == 404 {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "GitlabTags GetTag error")
	}

	return tag, nil
}

func createTag(ctx *cli.Context, pid string, tagName string, ref string, msg string, desc string) (*gitlab.Tag, error) {
	gl := gitlab.NewClient(nil, ctx.String("token"))
	tag, _, err := gl.Tags.CreateTag(pid, &gitlab.CreateTagOptions{
		TagName:            gitlab.String(tagName),
		Ref:                gitlab.String(ref),
		Message:            gitlab.String(msg),
		ReleaseDescription: gitlab.String(desc),
	})
	if err != nil {
		return nil, errors.Wrap(err, "GitlabTags CreateTag error")
	}

	return tag, nil
}

func deleteTag(ctx *cli.Context, pid string, t string) error {
	gl := gitlab.NewClient(nil, ctx.String("token"))
	if _, err := gl.Tags.DeleteTag(pid, t); err != nil {
		return errors.Wrap(err, "GitlabTags DeleteTag error")
	}

	return nil
}
