package provider

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"strings"
)

type Aws struct {
	*s3manager.Downloader
	filePath string
}

func init() {
	Constructors[AWS] = &ProviderSpec{
		New:         NewAws,
		description: "Aws provider is used to interact with Amazon's Simple Cloud Provider Service.",
	}
}

func (p Aws) Read() (string, error) {
	buf := aws.NewWriteAtBuffer([]byte{})
	urlParts := strings.Split(strings.TrimPrefix(p.filePath, "s3://"), "/")
	s3Bucket, s3Key := urlParts[0], strings.Join(urlParts[1:], "/")
	_, err := p.Downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(s3Key),
	})
	if err != nil {
		return "", fmt.Errorf("failed to download file: %v", err)
	}
	return string(buf.Bytes()), nil
}

func NewAws(infile string) (Provider, error) {
	return &Aws{
		Downloader: s3manager.NewDownloader(session.Must(session.NewSession())),
		filePath:   infile,
	}, nil
}
