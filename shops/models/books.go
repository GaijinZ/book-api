package models

type Author struct {
	ID    int    `json:"id"`
	Names string `json:"authors"`
}

type IndustryIdentifier struct {
	Identifier string `json:"identifier"`
}

type GoogleBooksRequest struct {
	Items []struct {
		VolumeInfo struct {
			Name          string               `json:"title"`
			DatePublished string               `json:"publishedDate"`
			ISBN          []IndustryIdentifier `json:"industryIdentifiers"`
			PageCount     int                  `json:"pageCount"`
			Authors       []string             `json:"authors"`
		} `json:"volumeInfo"`
	} `json:"items"`
}
