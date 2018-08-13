package main

import (
	"encoding/csv"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"github.com/tenntenn/goplayground"
)

type Validator struct {
	HTTPClient *http.Client

	r       *csv.Reader
	nameRow int
	urlRow  int

	name   string
	result *goplayground.RunResult
	err    error
}

func NewValidator(r *csv.Reader, urlRow, nameRow int) *Validator {
	return &Validator{
		r:       r,
		urlRow:  urlRow,
		nameRow: nameRow,
	}
}

func (v *Validator) httpClient() *http.Client {
	if v.HTTPClient == nil {
		return http.DefaultClient
	}
	return v.HTTPClient
}

func (v *Validator) Next() bool {
	v.result, v.name, v.err = nil, "", nil

	record, err := v.r.Read()
	if err == io.EOF {
		return false
	}

	if err != nil {
		v.err = errors.Wrap(err, "cannot get reuslt")
		return false
	}

	if v.urlRow < 0 || len(record) <= v.nameRow {
		v.err = errors.Errorf("invalid nameRow %d", v.nameRow)
		return false
	}

	if v.urlRow < 0 || len(record) <= v.urlRow {
		v.err = errors.Errorf("invalid urlRow %d", v.urlRow)
		return false
	}

	url := record[v.urlRow]
	if url == "" || !strings.HasPrefix(url, goplayground.BaseURL) {
		return true
	}

	src, err := v.getSrc(url)
	if err != nil {
		v.err = errors.Wrap(err, "cannot get src")
		return false
	}

	var cli goplayground.Client
	result, err := cli.Run(src)
	if err != nil {
		v.err = errors.Wrap(err, "cannot run src")
		return false
	}
	v.result = result
	v.name = record[v.nameRow]

	return true
}

func (v *Validator) getSrc(url string) (string, error) {
	resp, err := v.httpClient().Get(url)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return "", err
	}
	return doc.Find("#code").Text(), nil
}

func (v *Validator) Name() string {
	return v.name
}

func (v *Validator) Result() *goplayground.RunResult {
	return v.result
}

func (v *Validator) Err() error {
	return v.err
}
