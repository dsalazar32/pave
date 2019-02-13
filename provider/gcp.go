package provider

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
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
	gsBucket, gsKey, err := parseCloudStorageUrl(p.filePath)
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

func (p Gcp) Write(s string) error {
	gsBucket, gsKey, err := parseCloudStorageUrl(p.filePath)
	if err != nil {
		return err
	}

	bkt := p.Client.Bucket(gsBucket)
	obj := bkt.Object(gsKey)
	w := obj.NewWriter(p.ctx)
	if _, err := fmt.Fprintf(w, s); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return nil
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
