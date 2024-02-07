package model

type SpotifyTopArtists struct {
	Items []SpotifyArtist `json:"items"`
}
type SpotifyArtist struct {
	Name   string
	Genres []string
	Images []SpotifyImage
}
type SpotifyImage struct {
	Url    string
	Height int
	Width  int
}

type SearchResult struct {
	Tracks Tracks `json:"tracks"`
}

type ExternalUrls struct {
	Spotify string `json:"spotify"`
}
type Artists struct {
	ExternalUrls ExternalUrls `json:"external_urls"`
	Href         string       `json:"href"`
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Type         string       `json:"type"`
	URI          string       `json:"uri"`
}

type Images struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}
type Album struct {
	AlbumType            string       `json:"album_type"`
	Artists              []Artists    `json:"artists"`
	AvailableMarkets     []string     `json:"available_markets"`
	ExternalUrls         ExternalUrls `json:"external_urls"`
	Href                 string       `json:"href"`
	ID                   string       `json:"id"`
	Images               []Images     `json:"images"`
	Name                 string       `json:"name"`
	ReleaseDate          string       `json:"release_date"`
	ReleaseDatePrecision string       `json:"release_date_precision"`
	TotalTracks          int          `json:"total_tracks"`
	Type                 string       `json:"type"`
	URI                  string       `json:"uri"`
}
type ExternalIds struct {
	Isrc string `json:"isrc"`
}
type Items struct {
	Album            Album        `json:"album"`
	Artists          []Artists    `json:"artists"`
	AvailableMarkets []string     `json:"available_markets"`
	DiscNumber       int          `json:"disc_number"`
	DurationMs       int          `json:"duration_ms"`
	Explicit         bool         `json:"explicit"`
	ExternalIds      ExternalIds  `json:"external_ids"`
	ExternalUrls     ExternalUrls `json:"external_urls"`
	Href             string       `json:"href"`
	ID               string       `json:"id"`
	IsLocal          bool         `json:"is_local"`
	Name             string       `json:"name"`
	Popularity       int          `json:"popularity"`
	PreviewURL       any          `json:"preview_url"`
	TrackNumber      int          `json:"track_number"`
	Type             string       `json:"type"`
	URI              string       `json:"uri"`
}
type Tracks struct {
	Href     string  `json:"href"`
	Items    []Items `json:"items"`
	Limit    int     `json:"limit"`
	Next     string  `json:"next"`
	Offset   int     `json:"offset"`
	Previous any     `json:"previous"`
	Total    int     `json:"total"`
}

// type SearchResult struct {
// 	Tracks struct {
// 		Href  string `json:"href"`
// 		Items []struct {
// 			Album struct {
// 				AlbumType string `json:"album_type"`
// 				Artists   []struct {
// 					ExternalUrls struct {
// 						Spotify string `json:"spotify"`
// 					} `json:"external_urls"`
// 					Href string `json:"href"`
// 					ID   string `json:"id"`
// 					Name string `json:"name"`
// 					Type string `json:"type"`
// 					URI  string `json:"uri"`
// 				} `json:"artists"`
// 				AvailableMarkets []string `json:"available_markets"`
// 				ExternalUrls     struct {
// 					Spotify string `json:"spotify"`
// 				} `json:"external_urls"`
// 				Href   string `json:"href"`
// 				ID     string `json:"id"`
// 				Images []struct {
// 					Height int    `json:"height"`
// 					URL    string `json:"url"`
// 					Width  int    `json:"width"`
// 				} `json:"images"`
// 				Name                 string `json:"name"`
// 				ReleaseDate          string `json:"release_date"`
// 				ReleaseDatePrecision string `json:"release_date_precision"`
// 				TotalTracks          int    `json:"total_tracks"`
// 				Type                 string `json:"type"`
// 				URI                  string `json:"uri"`
// 			} `json:"album"`
// 			Artists []struct {
// 				ExternalUrls struct {
// 					Spotify string `json:"spotify"`
// 				} `json:"external_urls"`
// 				Href string `json:"href"`
// 				ID   string `json:"id"`
// 				Name string `json:"name"`
// 				Type string `json:"type"`
// 				URI  string `json:"uri"`
// 			} `json:"artists"`
// 			AvailableMarkets []string `json:"available_markets"`
// 			DiscNumber       int      `json:"disc_number"`
// 			DurationMs       int      `json:"duration_ms"`
// 			Explicit         bool     `json:"explicit"`
// 			ExternalIds      struct {
// 				Isrc string `json:"isrc"`
// 			} `json:"external_ids"`
// 			ExternalUrls struct {
// 				Spotify string `json:"spotify"`
// 			} `json:"external_urls"`
// 			Href        string `json:"href"`
// 			ID          string `json:"id"`
// 			IsLocal     bool   `json:"is_local"`
// 			Name        string `json:"name"`
// 			Popularity  int    `json:"popularity"`
// 			PreviewURL  any    `json:"preview_url"`
// 			TrackNumber int    `json:"track_number"`
// 			Type        string `json:"type"`
// 			URI         string `json:"uri"`
// 		} `json:"items"`
// 		Limit    int    `json:"limit"`
// 		Next     string `json:"next"`
// 		Offset   int    `json:"offset"`
// 		Previous any    `json:"previous"`
// 		Total    int    `json:"total"`
// 	} `json:"tracks"`
// }
