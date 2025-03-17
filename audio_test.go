package main

import (
	"encoding/json"
	"testing"
)

func Test_getVideoInfo(t *testing.T) {
	video, err := getVideoInfo("BVxxx", "")
	if err != nil {
		t.Fatalf("get video error: %v", err)
	}
	if video.Code != 0 {
		t.Fatalf("get video failed with code %d, msg: %s", video.Code, video.Message)
	}

	videoJson, _ := json.MarshalIndent(video, "", "\t")
	t.Logf("video:\n%s", videoJson)
}

func Test_getDash(t *testing.T) {
	const cookie = ""

	video, err := getVideoInfo("BVxxx", "")
	if err != nil {
		t.Fatalf("get video error: %v", err)
	}
	if video.Code != 0 {
		t.Fatalf("get video failed with code %d, msg: %s", video.Code, video.Message)
	}

	dash, err := getDash(video.Data.Aid, video.Data.Cid, 127, video.Data.Bvid, cookie)
	if err != nil {
		t.Fatalf("get dash error: %v", err)
	}
	dashJson, _ := json.MarshalIndent(dash, "", "\t")
	t.Logf("dash:\n%s", dashJson)
}
