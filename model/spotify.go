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
	Height string
	Width  string
}
