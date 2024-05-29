package entity

type RecommendedKomik struct {
	Title     string `json:"title"`
	Href      string `json:"href"`
	Rating    string `json:"rating"`
	Chapter   string `json:"chapter"`
	Type      string `json:"type"`
	Thumbnail string `json:"thumbnail"`
}
