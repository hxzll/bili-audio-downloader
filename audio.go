package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	bilibiliVideoInfoAPI = "https://api.bilibili.com/x/web-interface/view?"
	bilibiliDashAPI      = "https://api.bilibili.com/x/player/playurl?"
)

type video struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Bvid    string `json:"bvid"`
		Aid     int    `json:"aid"`
		Cid     int    `json:"cid"`
		Videos  int    `json:"videos"`
		Title   string `json:"title"`
		Pubdate int    `json:"pubdate"`
		Ctime   int    `json:"ctime"`
		Pages   []struct {
			Cid  int    `json:"cid"`
			Page int    `json:"page"`
			Part string `json:"part"`
		} `json:"pages"`
	} `json:"data"`
}

// ref: https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/docs/video/info.md
func getVideoInfo(bvid string, cookie string) (*video, error) {
	apiURL := bilibiliVideoInfoAPI + "bvid=" + bvid
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

	var respData video
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return nil, err
	}

	return &respData, nil
}

type dash struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		CurQuality  int      `json:"quality"`
		Description []string `json:"accept_description"`
		Quality     []int    `json:"accept_quality"`
		Dash        struct {
			Audio []struct {
				ID        int    `json:"id"`
				BaseURL   string `json:"baseUrl"`
				Bandwidth int    `json:"bandwidth"`
				MimeType  string `json:"mimeType"`
				Codecid   int    `json:"codecid"`
				Codecs    string `json:"codecs"`
			} `json:"audio"`
		} `json:"dash"`
		DURLFormat string `json:"format"`
		DURLs      []struct {
			URL  string `json:"url"`
			Size int64  `json:"size"`
		} `json:"durl"`
	} `json:"data"`
}

// call the api which its response body contains video and audio download streams.
//
// ref: https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/docs/video/videostream_url.md
func getDash(aid, cid, quality int, bvid string, cookie string) (*dash, error) {
	apiURL := fmt.Sprintf(
		bilibiliDashAPI+"avid=%d&cid=%d&bvid=%s&qn=%d&type=&otype=json&fourk=1&fnver=0&fnval=2000",
		aid, cid, bvid, quality,
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

	var respData dash
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return nil, err
	}

	return &respData, nil
}

type AudioStream struct {
	Title      string // video tilte
	URL        string
	Bvid       string
	Cid        int
	RefererURL string
	MimeType   string
}

// GetAudioStreams returns audios streams by the `bvid`,
// returns more than one stream if the `bvid` refers more than one video part.
//
// `cookie` is required only if the video is only availble as a login state user.
func GetAudioStreams(bvid string, cookie string) ([]AudioStream, error) {
	if !strings.HasPrefix(bvid, "BV") {
		return nil, fmt.Errorf("invalid bvid: %q", bvid)
	}

	videoInfo, err := getVideoInfo(bvid, cookie)
	if err != nil {
		return nil, err
	}
	if videoInfo.Code != 0 {
		return nil, fmt.Errorf("get video info failed, code=%d, msg=%s", videoInfo.Code, videoInfo.Message)
	}

	var streams []AudioStream

	for _, page := range videoInfo.Data.Pages {
		// 127 for highest quality
		dashInfo, err := getDash(videoInfo.Data.Aid, page.Cid, 127, videoInfo.Data.Bvid, cookie)
		if err != nil {
			return nil, err
		}
		if dashInfo.Code != 0 {
			return nil, fmt.Errorf("get dash info failed, code=%d, msg=%s", dashInfo.Code, dashInfo.Message)
		}
		if len(dashInfo.Data.Dash.Audio) == 0 {
			return nil, fmt.Errorf("no audio found for video %s", bvid)
		}

		// get the max bandwidth audio
		auido := dashInfo.Data.Dash.Audio[0]
		for _, a := range dashInfo.Data.Dash.Audio {
			if a.Bandwidth > auido.Bandwidth {
				auido = a
			}
		}

		title := videoInfo.Data.Title
		if videoInfo.Data.Videos > 1 {
			title += fmt.Sprintf("_p%d_%s", page.Page, page.Part)
		}
		streams = append(streams, AudioStream{
			Title:      title,
			Bvid:       videoInfo.Data.Bvid,
			Cid:        page.Cid,
			URL:        auido.BaseURL,
			RefererURL: "https://www.bilibili.com/video/" + videoInfo.Data.Bvid,
			MimeType:   dashInfo.Data.DURLFormat,
		})
	}

	return streams, nil
}
