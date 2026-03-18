package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
	"gopkg.in/vansante/go-ffprobe.v2"
)

var dateRegex = regexp.MustCompile(`^([12]\d\d\d)`)

func parseArtist(s string) []string {
	if s == "" {
		return []string{}
	}

	splits := strings.Split(s, ",")

	artists := make([]string, 0, len(splits))
	for _, s := range splits {
		a := strings.TrimSpace(s)

		if a != "" {
			artists = append(artists, a)
		}
	}

	return artists
}

func convertMapKeysToLowercase(m map[string]any) map[string]any {
	res := make(map[string]any)
	for k, v := range m {
		res[strings.ToLower(k)] = v
	}

	return res
}

type ProbeResult struct {
	Tags        ffprobe.Tags
	MediaFormat types.MediaFormat
	Duration    time.Duration
}

// TODO(patrik): This is the same as the ProbeMedia from the service.MediaService
func ProbeMedia(filepath string) (*ProbeResult, error) {
	ctx := context.TODO()

	probe, err := ffprobe.ProbeURL(ctx, filepath)
	if err != nil {
		return nil, err
	}

	var tags ffprobe.Tags
	hasGlobalTags := probe.Format.FormatName != "ogg"

	audioStream := probe.FirstAudioStream()
	if audioStream == nil {
		// TODO(patrik): Better error?
		return nil, errors.New("contains no audio streams")
	}

	if hasGlobalTags {
		tags = probe.Format.TagList
	} else {
		tags = audioStream.TagList
	}

	tags = convertMapKeysToLowercase(tags)

	dur, err := strconv.ParseFloat(audioStream.Duration, 32)
	if err != nil {
		return nil, err
	}

	duration := time.Duration(dur * float64(time.Second))

	mediaFormat := types.MediaFormatUnknown
	switch audioStream.CodecName {
	case "flac":
		mediaFormat = types.MediaFormatFlac
	case "pcm_s16le":
		mediaFormat = types.MediaFormatPcmS16LE
	case "opus":
		mediaFormat = types.MediaFormatOpus
	case "vorbis":
		mediaFormat = types.MediaFormatVorbis
	case "mp3":
		mediaFormat = types.MediaFormatMp3
	case "aac":
		mediaFormat = types.MediaFormatAac
	}

	return &ProbeResult{
		Tags:        tags,
		MediaFormat: mediaFormat,
		Duration:    duration,
	}, nil
}

type TrackInfo struct {
	Name   string
	Artist string
	Number int
	Year   int
}

func getTrackInfo(p string) (TrackInfo, error) {
	probe, err := ProbeMedia(p)
	if err != nil {
		return TrackInfo{}, err
	}

	var res TrackInfo

	res.Name, _ = probe.Tags.GetString("title")
	res.Artist, _ = probe.Tags.GetString("artist")

	if tag, err := probe.Tags.GetString("date"); err == nil {
		match := dateRegex.FindStringSubmatch(tag)
		if len(match) > 0 {
			res.Year, _ = strconv.Atoi(match[1])
		}
	}

	if tag, err := probe.Tags.GetInt("track"); err == nil {
		res.Number = int(tag)
	}

	return res, nil
}

var initCmd = &cobra.Command{
	Use: "init",
}

var initAlbumCmd = &cobra.Command{
	Use: "album",
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")
		output, _ := cmd.Flags().GetString("output")

		metadata := types.AlbumMetadata{}

		extract := true

		// TODO(patrik): Discard hidden files (starts with .)
		entries, err := os.ReadDir(dir)
		if err != nil {
			slog.Error("failed to read dir", "err", err)
			os.Exit(1)
		}

		var tracks []string
		var images []string

		for _, e := range entries {
			if e.IsDir() {
				continue
			}

			p := path.Join(dir, e.Name())

			ext := path.Ext(p)

			if utils.IsValidTrackExt(ext) {
				tracks = append(tracks, p)
			}

			if utils.IsValidImageExt(ext) {
				images = append(images, p)
			}
		}

		if len(tracks) <= 0 {
			slog.Warn("No tracks found... Quitting")
			return
		}

		p := tracks[0]

		probe, err := ProbeMedia(p)
		if err != nil {
			slog.Error("failed to probe track", "err", err)
			os.Exit(1)
		}

		isSingle := len(tracks) == 1

		metadata.Album.Id = utils.CreateAlbumId()

		if len(images) > 0 {
			// TODO(patrik): Better selection?
			metadata.General.Cover = images[0]
		}

		if !isSingle {
			metadata.Album.Name, _ = probe.Tags.GetString("album")
		} else {
			// NOTE(patrik): If we only have one track then we make the
			// album name the same as the track name
			metadata.Album.Name, _ = probe.Tags.GetString("title")
		}

		if !isSingle {
			if tag, err := probe.Tags.GetString("album_artist"); err == nil {
				metadata.Album.Artists = parseArtist(tag)
			} else {
				if tag, err := probe.Tags.GetString("artist"); err == nil {
					metadata.Album.Artists = parseArtist(tag)
				}
			}
		} else {
			if tag, err := probe.Tags.GetString("artist"); err == nil {
				metadata.Album.Artists = parseArtist(tag)
			}
		}

		if tag, err := probe.Tags.GetString("date"); err == nil {
			match := dateRegex.FindStringSubmatch(tag)
			if len(match) > 0 {
				metadata.General.Year, _ = strconv.ParseInt(match[1], 10, 64)
			}
		}

		for _, p := range tracks {
			filename := path.Base(p)
			fmt.Printf("Found track: %s\n", filename)

			trackInfo, err := getTrackInfo(p)
			if err != nil {
				slog.Error("failed to get track info", "err", err)
				os.Exit(1)
			}

			if trackInfo.Name == "" {
				trackInfo.Name = strings.TrimSuffix(filename, path.Ext(p))
			}

			if trackInfo.Number == 0 || extract {
				trackInfo.Number = utils.ExtractNumber(filename)
			}

			// TODO(patrik): If artist is empty then use album maybe
			artists := parseArtist(trackInfo.Artist)

			metadata.Tracks = append(metadata.Tracks, types.AlbumMetadataTrack{
				Id:      utils.CreateTrackId(),
				File:    filename,
				Name:    trackInfo.Name,
				Number:  int64(trackInfo.Number),
				Year:    0,
				Tags:    []string{},
				Artists: artists,
			})
		}

		data, err := toml.Marshal(&metadata)
		if err != nil {
			slog.Error("failed to marshal metadata", "err", err)
			os.Exit(1)
		}

		err = os.WriteFile(output, data, 0644)
		if err != nil {
			slog.Error("failed to write output", "err", err)
			os.Exit(1)
		}
	},
}

func downloadImage(url, dir, name string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to send http request: %w", err)
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", fmt.Errorf("failed to parse media type: %w", err)
	}

	ext := ""
	switch mediaType {
	case "image/png":
		ext = ".png"
	case "image/jpeg":
		ext = ".jpeg"
	default:
		return "", fmt.Errorf("unsupported media type: %s", mediaType)
	}

	p := path.Join(dir, name+ext)
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to open output file: %w", err)
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to copy http body to file: %w", err)
	}

	return p, nil
}

var initArtistCmd = &cobra.Command{
	Use: "artist",
	Run: func(cmd *cobra.Command, args []string) {
		out, _ := cmd.Flags().GetString("output")
		dirName, _ := cmd.Flags().GetBool("dir-name")

		// TODO(patrik): Add check for artist.toml already exists

		artistName := ""
		coverUrl := ""
		tags := ""

		if dirName {
			stat, _ := os.Stat(out)
			if stat != nil {
				artistName = stat.Name()
			}
		}

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Title("Artist Name").Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return errors.New("cannot be empty")
					}

					return nil
				}).Value(&artistName),
				huh.NewInput().Title("Tags").Value(&tags),
				huh.NewInput().Title("Cover URL").Value(&coverUrl),
			),
		)

		err := form.Run()
		if err != nil {
			slog.Error("failed to run form", "err", err)
			return
		}

		artistName = strings.TrimSpace(artistName)

		// TODO(patrik): Print out overview

		cover := ""

		if coverUrl != "" {
			if p, err := downloadImage(coverUrl, out, "cover"); err == nil {
				cover = p
			} else {
				slog.Error("failed to download cover image", "err", err)
			}
		}

		tagsArr := []string{}

		split := utils.SplitString(tags)
		for _, tag := range split {
			t := strings.TrimSpace(utils.Slug(tag))
			if t != "" {
				tagsArr = append(tagsArr, t)
			}
		}

		metadata := types.ArtistMetadata{
			Id:         utils.CreateArtistId(),
			SearchName: utils.Slug(artistName),
			Name:       artistName,
			Cover:      path.Base(cover),
			Tags:       tagsArr,
		}

		d, err := toml.Marshal(metadata)
		if err != nil {
			slog.Error("failed to marshal artist metadata", "err", err)
			return
		}

		p := path.Join(out, "artist.toml")
		err = os.WriteFile(p, d, 0644)
		if err != nil {
			slog.Error("failed to write artist metadata", "err", err, "path", p)
			return
		}
	},
}

func init() {
	initAlbumCmd.Flags().String("dir", ".", "input directory")
	initAlbumCmd.Flags().StringP("output", "o", "album.toml", "write result to file")

	initArtistCmd.Flags().BoolP("dir-name", "d", false, "take the parent directory name")
	initArtistCmd.Flags().String("output", ".", "output directory")

	initCmd.AddCommand(initAlbumCmd, initArtistCmd)

	rootCmd.AddCommand(initCmd)
}
