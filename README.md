bili-audio-downloader
===

个人使用，目前仅支持下载收藏夹中的视频的音频

# 安装

```shell
go install github.com/hxzll/bili-audio-downloader@latest
```

# 使用

```shell
USAGE:
   bili-audio-downloader favlist

OPTIONS:
   --fid value                         favlist id
   --output value, -o value            output directory, where audios files stores in.
   --items value [ --items value ]     items to be downloaded, start from 1, order from newest to ordest. Example: 1,2,3.
   --startBvid value, --startbv value  startBvid downloads the newest videos starting from this video, inluding it.
   --endBvid value, --endbv value      endBvid downloads the newest videos util this video, inluding it.
   --startOid value, --startid value   startOid downloads the newest videos starting from this video, inluding it.
   --endOid value, --endid value       endOid downloads the newest videos util this video, inluding it.
   --cookie value, -c value            the cookie of the bilibili web, used for download login state only data.
   --help, -h                          show help
```
