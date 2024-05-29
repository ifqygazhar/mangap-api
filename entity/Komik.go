package entity

type Komik struct {
	Title     string `json:"title"`
	Href      string `json:"href"`
	Genre     string `json:"genre"`
	Year      string `json:"year"`
	Thumbnail string `json:"thumbnail"`
}
