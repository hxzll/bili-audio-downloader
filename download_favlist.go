package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"
	"slices"
)

const (
	bilibiliFavlistAPI = "https://api.bilibili.com/x/v3/fav/resource/ids?"
)

type favlist struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []struct {
		Id   int64  `json:"id"`
		Type int    `json:"type"`
		Bvid string `json:"bvid"`
	} `json:"data"`
}

func getFavlist(fid int64, cookie string) (*favlist, error) {
	apiURL := fmt.Sprintf(
		bilibiliFavlistAPI+"media_id=%d&platform=web", fid,
	)
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build http request: %v", err)
	}
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var respData favlist
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return nil, err
	}

	return &respData, nil
}

type DownloadFavlistOption struct {
	// Fid is the favlist id.
	Fid int64
	// Items set the items to be downloaded, start from 1, order from newest to ordest.
	Items []int64
	// StartBvid downloads the newest videos starting from this video, inluding it.
	StartBvid string
	// EndBvid downloads the newest videos util this video, inluding it.
	EndBvid string
	// StartOid downloads the newest videos starting from this video, inluding it.
	StartOid int64
	// EndOid downloads the newest videos util this video, inluding it.
	EndOid    int64
	Cookie    string
	OutputDir string
}

func DownloadFavlist(opt DownloadFavlistOption) error {
	if opt.Fid == 0 {
		return errors.New("favlist id is required")
	}
	if opt.OutputDir == "" {
		return errors.New("output dir is not set")
	}

	favlist, err := getFavlist(opt.Fid, opt.Cookie)
	if err != nil {
		return err
	}
	if favlist.Code != 0 {
		return fmt.Errorf("get favlist failed, code=%d, msg=%s", favlist.Code, favlist.Message)
	}
	if len(favlist.Data) == 0 {
		fmt.Println("favlist is empty")
		return nil
	}

	var bvids []string
	started := true
	if opt.StartBvid != "" {
		started = false
	}
	if opt.StartOid != 0 {
		started = false
	}
	for i, media := range favlist.Data {
		if !started && opt.StartBvid != "" && media.Bvid == opt.StartBvid {
			started = true
		}
		if !started && opt.StartOid != 0 && media.Id == opt.StartOid {
			started = true
		}
		if !started {
			continue
		}

		if len(opt.Items) > 0 && !slices.Contains(opt.Items, int64(i+1)) {
			continue
		}

		bvids = append(bvids, media.Bvid)

		if opt.EndBvid != "" && media.Bvid == opt.EndBvid {
			break
		}
		if opt.EndOid != 0 && media.Id == opt.EndOid {
			break
		}
	}

	for _, bvid := range bvids {
		audios, err := GetAudioStreams(bvid, opt.Cookie)
		if err != nil {
			return err
		}

		for _, a := range audios {
			filename := a.Title
			outputPath := path.Join(opt.OutputDir, SanitizeFilename(filename)+".m4a")

			if err := Download(a.URL, a.RefererURL, outputPath); err != nil {
				return err
			}
		}
	}

	fmt.Println("download success.")
	return nil
}
