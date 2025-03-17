package main

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/urfave/cli/v3"
)

const appName = "bili-audio-downloader"

var version string

func init() {
	buildInfo, _ := debug.ReadBuildInfo()
	version = buildInfo.Main.Version
}

func main() {
	cmd := &cli.Command{
		Name:    appName,
		Version: version,
		Commands: []*cli.Command{
			{
				Name:    "favlist",
				Aliases: []string{"fav"},
				Usage:   "options for favlist audio download",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:     "fid",
						Usage:    "favlist id",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "output",
						Usage:    "output directory, where audios files stores in.",
						Aliases:  []string{"o"},
						Required: true,
					},
					&cli.IntSliceFlag{
						Name:  "items",
						Usage: "items to be downloaded, start from 1, order from newest to ordest. Example: 1,2,3.",
					},
					&cli.IntFlag{
						Name:        "startBvid",
						Usage:       "startBvid downloads the newest videos starting from this video, inluding it.",
						Aliases:     []string{"startbv"},
						HideDefault: true,
					},
					&cli.IntFlag{
						Name:        "endBvid",
						Usage:       "endBvid downloads the newest videos util this video, inluding it.",
						Aliases:     []string{"endbv"},
						HideDefault: true,
					},
					&cli.IntFlag{
						Name:        "startOid",
						Usage:       "startOid downloads the newest videos starting from this video, inluding it.",
						Aliases:     []string{"startid"},
						HideDefault: true,
					},
					&cli.IntFlag{
						Name:        "endOid",
						Usage:       "endOid downloads the newest videos util this video, inluding it.",
						Aliases:     []string{"endid"},
						HideDefault: true,
					},
					&cli.StringFlag{
						Name:    "cookie",
						Usage:   "the cookie of the bilibili web, used for download login state only data.",
						Aliases: []string{"c"},
					},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					opt := DownloadFavlistOption{
						Fid:       cmd.Int("fid"),
						Items:     cmd.IntSlice("items"),
						StartBvid: cmd.String("startBvid"),
						EndBvid:   cmd.String("endBvid"),
						StartOid:  cmd.Int("startOid"),
						EndOid:    cmd.Int("endOid"),
						Cookie:    cmd.String("cookie"),
						OutputDir: cmd.String("output"),
					}

					if err := DownloadFavlist(opt); err != nil {
						return err
					}

					return nil
				},
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "set program to debug mode.",
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Println(err)
	}
}
