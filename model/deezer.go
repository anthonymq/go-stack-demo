package model

import "strconv"

type DeezerSearchTrackResults struct {
	Data []struct {
		ID                    int    `json:"id"`
		Readable              bool   `json:"readable"`
		Title                 string `json:"title"`
		TitleShort            string `json:"title_short"`
		TitleVersion          string `json:"title_version"`
		Link                  string `json:"link"`
		Duration              int    `json:"duration"`
		Rank                  int    `json:"rank"`
		ExplicitLyrics        bool   `json:"explicit_lyrics"`
		ExplicitContentLyrics int    `json:"explicit_content_lyrics"`
		ExplicitContentCover  int    `json:"explicit_content_cover"`
		Preview               string `json:"preview"`
		Md5Image              string `json:"md5_image"`
		Artist                struct {
			ID            int    `json:"id"`
			Name          string `json:"name"`
			Link          string `json:"link"`
			Picture       string `json:"picture"`
			PictureSmall  string `json:"picture_small"`
			PictureMedium string `json:"picture_medium"`
			PictureBig    string `json:"picture_big"`
			PictureXl     string `json:"picture_xl"`
			Tracklist     string `json:"tracklist"`
			Type          string `json:"type"`
		} `json:"artist"`
		Album struct {
			ID          int    `json:"id"`
			Title       string `json:"title"`
			Cover       string `json:"cover"`
			CoverSmall  string `json:"cover_small"`
			CoverMedium string `json:"cover_medium"`
			CoverBig    string `json:"cover_big"`
			CoverXl     string `json:"cover_xl"`
			Md5Image    string `json:"md5_image"`
			Tracklist   string `json:"tracklist"`
			Type        string `json:"type"`
		} `json:"album"`
		Type string `json:"type"`
	} `json:"data"`
	Total int    `json:"total"`
	Next  string `json:"next"`
}

func (m DeezerSearchTrackResults) ToSearchResults() []SearchResult {
	var results []SearchResult
	for _, track := range m.Data {
		results = append(results, SearchResult{
			Id:      strconv.Itoa(track.ID),
			Title:   track.Title,
			Artists: []string{track.Artist.Name},
			Cover:   track.Album.Cover,
		})
	}
	return results
}

type DeezerGetPlaylists struct {
	Data []struct {
		ID            int    `json:"id"`
		Title         string `json:"title"`
		Duration      int    `json:"duration"`
		Public        bool   `json:"public"`
		IsLovedTrack  bool   `json:"is_loved_track"`
		Collaborative bool   `json:"collaborative"`
		NbTracks      int    `json:"nb_tracks"`
		Fans          int    `json:"fans"`
		Link          string `json:"link"`
		Picture       string `json:"picture"`
		PictureSmall  string `json:"picture_small"`
		PictureMedium string `json:"picture_medium"`
		PictureBig    string `json:"picture_big"`
		PictureXl     string `json:"picture_xl"`
		Checksum      string `json:"checksum"`
		Tracklist     string `json:"tracklist"`
		CreationDate  string `json:"creation_date"`
		Md5Image      string `json:"md5_image"`
		PictureType   string `json:"picture_type"`
		TimeAdd       int    `json:"time_add"`
		TimeMod       int    `json:"time_mod"`
		Creator       struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Tracklist string `json:"tracklist"`
			Type      string `json:"type"`
		} `json:"creator"`
		Type string `json:"type"`
	} `json:"data"`
	Total int `json:"total"`
}
