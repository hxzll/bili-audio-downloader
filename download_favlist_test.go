package main

import (
	"testing"
)

func TestDownloadFavlist(t *testing.T) {
	const cookie = ""
	err := DownloadFavlist(DownloadFavlistOption{
		Fid:       0,
		Cookie:    cookie,
		OutputDir: "output",
	})
	if err != nil {
		t.Fatalf("download favlist error: %v", err)
	}
}
