package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/urfave/cli"

	"github.com/pkg/errors"
)

const (
	// APIBaseURL ...
	APIBaseURL = "https://gitlab.com/api/v4"
)

// UploadsResponse ...
type UploadsResponse struct {
	Alt      string `json:"alt"`
	URL      string `json:"url"`
	Markdown string `json:"markdown"`
}

func multiUploads(ctx *cli.Context, p string, dir string) ([]*UploadsResponse, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var resps []*UploadsResponse
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		path := filepath.Join(dir, file.Name())
		resp, err := uploads(ctx, p, path)
		if err != nil {
			return nil, err
		}
		resps = append(resps, resp)
	}

	return resps, nil
}

func uploads(ctx *cli.Context, p string, path string) (*UploadsResponse, error) {
	pid := url.QueryEscape(p)
	u := fmt.Sprintf("%s/projects/%s/uploads", APIBaseURL, pid)
	b := &bytes.Buffer{}
	w := multipart.NewWriter(b)
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "Os Open error")
	}
	defer f.Close()

	fw, err := w.CreateFormFile("file", path)
	if err != nil {
		return nil, errors.Wrap(err, "MultipartWriter CreateFromFile error")
	}
	_, err = io.Copy(fw, f)
	if err != nil {
		return nil, errors.Wrap(err, "IO Copy error")
	}
	w.Close()

	req, err := http.NewRequest("POST", u, b)
	if err != nil {
		return nil, errors.Wrap(err, "Http NewRequest error")
	}
	req.Header.Set("PRIVATE-TOKEN", ctx.String("token"))
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "HttpClient Do error")
	}

	defer res.Body.Close()
	ba, _ := ioutil.ReadAll(res.Body)
	ures := &UploadsResponse{}
	if err := json.Unmarshal(ba, ures); err != nil {
		return nil, errors.Wrap(err, "JSON Unmarshal error")
	}

	return ures, nil
}
