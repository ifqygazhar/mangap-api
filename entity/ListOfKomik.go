package entity

type ListOfKomik struct {
	Title     string `json:"title"`
	Chapter   string `json:"chapter"`
	Type      string `json:"type"`
	Thumbnail string `json:"thumbnail"`
	Rating    string `json:"rating"`
	Href      string `json:"href"`
	TotalData int    `json:"total_data,omitempty"`
}
