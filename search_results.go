package youtube

import (
	sjson "github.com/bitly/go-simplejson"
	"strconv"
	"strings"
	"time"
)

type SearchResults struct {
	Query  string
	Videos []*ResultVideo
}

type ResultVideo struct {
	ID       string
	Title    string
	Author   string
	Duration time.Duration
}

func NewSearchResults(query string, body []byte) (*SearchResults, error) {
	res := &SearchResults{
		Query: query,
	}
	err := res.parseSearchResultsInfo(body)
	return res, err
}

func (s *SearchResults) parseSearchResultsInfo(body []byte) error {

	j, err := sjson.NewJson(body)
	if err != nil {
		return err
	}
	j = j.GetPath("contents", "twoColumnSearchResultsRenderer", "primaryContents", "sectionListRenderer", "contents")
	videosData := j.GetIndex(0).GetPath("itemSectionRenderer", "contents")
	videos := make([]*ResultVideo, 0, len(videosData.MustArray()))

	for i := range videosData.MustArray() {
		videoData, ok := videosData.GetIndex(i).CheckGet("videoRenderer")
		if !ok {
			continue
		}
		videoId, _ := videoData.Get("videoId").String()
		title, _ := videoData.GetPath("title", "runs").GetIndex(0).Get("text").String()
		author, _ := videoData.GetPath("ownerText", "runs").GetIndex(0).Get("text").String()
		durationString, _ := videoData.GetPath("lengthText", "simpleText").String()
		duration := convertDuration(durationString)
		resultVideo := &ResultVideo{
			ID:       videoId,
			Title:    title,
			Author:   author,
			Duration: duration,
		}
		videos = append(videos, resultVideo)
	}
	s.Videos = videos
	return nil
}

func convertDuration(durationString string) time.Duration {
	partsTime := strings.Split(durationString, ":")
	lenPartsTime := len(partsTime)
	durationSec, _ := strconv.Atoi(partsTime[0])
	duration := time.Second * time.Duration(durationSec)
	if lenPartsTime > 1 {
		durationMin, _ := strconv.Atoi(partsTime[1])
		duration += time.Minute * time.Duration(durationMin)
	}
	if lenPartsTime > 2 {
		durationHour, _ := strconv.Atoi(partsTime[2])
		duration += time.Hour * time.Duration(durationHour)
	}

	return duration
}
