package youtube

import (
	sjson "github.com/bitly/go-simplejson"
)

type SearchReslts struct {
	Query      string
	maxResults int
	Videos     []*Video
}

func NewSearchResults(query string, maxResults int, client *Client, body []byte) (*SearchReslts, error) {
	res := &SearchReslts{
		Query:      query,
		maxResults: maxResults,
	}
	err := res.parseSearchResultsInfo(client, body)
	return res, err
}

func (s *SearchReslts) parseSearchResultsInfo(client *Client, body []byte) error {
	videos := make([]*Video, 0, s.maxResults)
	j, err := sjson.NewJson(body)
	if err != nil {
		return err
	}
	j = j.GetPath("contents", "twoColumnSearchResultsRenderer", "primaryContents", "sectionListRenderer", "contents")
	videoDatas, err := j.GetIndex(0).GetPath("itemSectionRenderer", "contents").Array()
	if err != nil {
		return err
	}
	for _, videoData := range videoDatas {
		content, _ := videoData.(map[string]interface{})
		vRawI, ok := content["videoRenderer"]
		if !ok {
			continue
		}
		vRaw, _ := vRawI.(map[string]interface{})
		url, _ := vRaw["videoId"].(string)
		video, err := client.GetVideo(url)
		if err != nil {
			if video != nil {
				continue
			}
			return err
		}
		videos = append(videos, video)
		if len(videos) == s.maxResults {
			break
		}
	}
	s.Videos = videos
	return nil
}
