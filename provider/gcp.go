package provider

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"io"
	"strings"
)

type Gcp struct {
	storage.Client
	filePath string
	ctx      context.Context
}

func init() {
	Constructors[GCP] = &ProviderSpec{
		New:         NewGcp,
		description: "Gcp provider is used to interact with Google's Storage Service",
	}
}

func (p Gcp) Read() (string, error) {
	urlParts := strings.Split(strings.TrimPrefix(p.filePath, "gs://"), "/")
	gsBucket, gsKey := urlParts[0], strings.Join(urlParts[1:], "/")
	bkt := p.Client.Bucket(gsBucket)
	obj := bkt.Object(gsKey)
	r, err := obj.NewReader(p.ctx)
	if err != nil {
		return "", err
	}
	defer r.Close()

	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, r); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func NewGcp(infile string) (Provider, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &Gcp{
		Client:   *client,
		filePath: infile,
		ctx:      ctx,
	}, nil
}
