package model

func (m SpotifySearchResult) ToSearchResult() []SearchResult {
	var searchResults []SearchResult
	for _, track := range m.Tracks.Items {
		searchResults = append(searchResults, SearchResult{
			Title:   track.Name,
			Artists: []string{track.Artists[0].Name},
			Cover:   track.Album.Images[0].URL,
		})
	}
	return searchResults
}

type SpotifyTopArtists struct {
	Items []SpotifyArtist `json:"items"`
}
type SpotifyArtist struct {
	Name   string
	Genres []string
	Images []Image
}

type SpotifySearchResult struct {
	Tracks Tracks `json:"tracks"`
}

type ExternalUrls struct {
	Spotify string `json:"spotify"`
}
type Artist struct {
	// ExternalUrls ExternalUrls `json:"external_urls"`
	// Href         string       `json:"href"`
	// ID           string       `json:"id"`
	Name string `json:"name"`
	// Type         string       `json:"type"`
	// URI          string       `json:"uri"`
}

type Image struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}
type Album struct {
	// AlbumType            string       `json:"album_type"`
	Artists []Artist `json:"artists"`
	// AvailableMarkets     []string     `json:"available_markets"`
	// ExternalUrls         ExternalUrls `json:"external_urls"`
	// Href                 string       `json:"href"`
	// ID                   string       `json:"id"`
	Images []Image `json:"images"`
	Name   string  `json:"name"`
	// ReleaseDate          string       `json:"release_date"`
	// ReleaseDatePrecision string       `json:"release_date_precision"`
	// TotalTracks          int          `json:"total_tracks"`
	// Type                 string       `json:"type"`
	// URI                  string       `json:"uri"`
}
type ExternalIds struct {
	Isrc string `json:"isrc"`
}
type Item struct {
	Album   Album    `json:"album"`
	Artists []Artist `json:"artists"`
	// AvailableMarkets []string     `json:"available_markets"`
	// DiscNumber       int          `json:"disc_number"`
	// DurationMs       int          `json:"duration_ms"`
	// Explicit         bool         `json:"explicit"`
	// ExternalIds      ExternalIds  `json:"external_ids"`
	// ExternalUrls     ExternalUrls `json:"external_urls"`
	// Href             string       `json:"href"`
	// ID               string       `json:"id"`
	// IsLocal          bool         `json:"is_local"`
	Name string `json:"name"`
	// Popularity       int          `json:"popularity"`
	// PreviewURL       any          `json:"preview_url"`
	// TrackNumber      int          `json:"track_number"`
	// Type             string       `json:"type"`
	// URI              string       `json:"uri"`
}
type Tracks struct {
	// Href     string `json:"href"`
	Items []Item `json:"items"`
	// Limit    int    `json:"limit"`
	// Next     string `json:"next"`
	// Offset   int    `json:"offset"`
	// Previous any    `json:"previous"`
	// Total    int    `json:"total"`
}
