package entity

type DetailedKomik struct {
	Title       string    `json:"title"`
	AltTitle    string    `json:"altTitle"`
	UpdatedOn   string    `json:"updatedOn"`
	Rating      string    `json:"rating"`
	Status      string    `json:"status"`
	Type        string    `json:"type"`
	Released    string    `json:"released"`
	Author      string    `json:"author"`
	Genre       []Genre   `json:"genre"`
	Description string    `json:"description"`
	Thumbnail   string    `json:"thumbnail"`
	Chapter     []Chapter `json:"chapter"`
}
