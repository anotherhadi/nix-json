package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const commitsURL = "https://api.github.com/repos/nix-community/nur-search/commits?page=1&per_page=1"

func GetLatestRelease(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, commitsURL, nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("github request failed: %w", err)
	}
	defer resp.Body.Close()

	commits := []struct {
		Sha string `json:"sha"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&commits)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal github response: %w", err)
	}

	if len(commits) < 1 {
		return "", fmt.Errorf("unexpected result from github: %w", err)
	}
	return commits[0].Sha, nil
}

const packagesURL = "https://raw.githubusercontent.com/nix-community/nur-search/%s/data/packages.json"

func DownloadRelease(ctx context.Context, release string) (io.ReadCloser, error) {
	apiurl := fmt.Sprintf(packagesURL, release)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiurl, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github request failed: %w", err)
	}

	return PackagesWrapper(resp.Body), nil
}

type packagesWrapper struct {
	pkgs io.Closer
	wrap io.Reader
}

// newOptionsWrapper translates a set of packages into the
// format of the indexer
//
// set of packages:
//
//	{
//		"pkg1": { ... },
//		"pkg2": { ... }
//	}
//
// what indexer expects:
//
//	{
//		"packages": {
//		  "pkg1": { ... },
//		  "pkg1": { ... }
//		}
//	}
func PackagesWrapper(rd io.ReadCloser) *packagesWrapper {
	mrd := io.MultiReader(
		bytes.NewBufferString(`{"packages":`),
		rd,
		bytes.NewBufferString(`}`),
	)
	return &packagesWrapper{
		pkgs: rd,
		wrap: mrd,
	}
}

func (w *packagesWrapper) Read(p []byte) (n int, err error) {
	return w.wrap.Read(p)
}

func (w *packagesWrapper) Close() error {
	return w.pkgs.Close()
}

func main() {
	ctx := context.Background()

	latest, err := GetLatestRelease(ctx)
	if err != nil {
		fmt.Println("Error fetching latest release:", err)
		return
	}
	fmt.Println("Latest release:", latest)

	pkgs, err := DownloadRelease(ctx, latest)
	if err != nil {
		fmt.Println("Error downloading release:", err)
		return
	}
	defer pkgs.Close()

	outputFile, err := os.Create("nur.json")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, pkgs)
	if err != nil {
		panic(err)
	}
	fmt.Println("Downloaded successfully to nur.json")
}
