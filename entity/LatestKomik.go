package entity

type LatestKomik struct {
	Title     string `json:"title"`
	Href      string `json:"href"`
	Thumbnail string `json:"thumbnail"`
	Type      string `json:"type"`
	Chapter   string `json:"chapter"`
	Rating    string `json:"rating"`
}
