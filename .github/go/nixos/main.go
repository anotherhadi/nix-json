package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const prefix = "nixos/unstable/"

func GetLatestRelease(ctx context.Context) (string, error) {
	s3client := s3.NewFromConfig(aws.Config{
		Region: "eu-west-1",
	})

	// The `startAfter` is a marker for S3 to start iterating from. Just use the latest
	// at the moment of writing nixpkgs release to never iterate from the beginning
	startAfter := "nixos/unstable/nixos-25.05beta751650.64e75cd44acf"
	var latest types.Object
	input := &s3.ListObjectsV2Input{
		Bucket:     aws.String("nix-releases"),
		Prefix:     aws.String(prefix),
		Delimiter:  aws.String("/"),
		StartAfter: aws.String(startAfter),
	}
	p := s3.NewListObjectsV2Paginator(s3client, input)
	for p.HasMorePages() {
		page, err := p.NextPage(ctx)
		if err != nil {
			return startAfter, fmt.Errorf("get next page: %w", err)
		}
		for _, obj := range page.Contents {
			latest = obj
		}
	}

	if latest.Key == nil {
		return startAfter, nil
	}
	return *latest.Key, nil
}

type readCloser struct {
	io.Reader
	io.Closer
}

func DownloadRelease(ctx context.Context, release string) (io.ReadCloser, error) {
	release = strings.TrimPrefix(release, prefix)
	fileURL, _ := url.JoinPath(
		"https://releases.nixos.org/nixos/unstable",
		release,
		"options.json.br",
	)

	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, fmt.Errorf("fetch packages: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("expected http 200, but got %d", resp.StatusCode)
	}

	br := brotli.NewReader(resp.Body)
	return &readCloser{
		Reader: br,
		Closer: resp.Body,
	}, nil
}

func main() {
	ctx := context.Background()

	release, err := GetLatestRelease(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("Latest release:", release)

	rc, err := DownloadRelease(ctx, release)
	if err != nil {
		panic(err)
	}
	defer rc.Close()

	outputFile, err := os.Create("nixos.json")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, rc)
	if err != nil {
		panic(err)
	}
	fmt.Println("Downloaded successfully to nixos.json")
}
