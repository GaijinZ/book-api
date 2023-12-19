package models

type Author struct {
	ID    int    `json:"id"`
	Names string `json:"authors"`
}

type ISBN struct {
	Identifier string `json:"identifier"`
}

type BooksRequest struct {
	Items []struct {
		VolumeInfo struct {
			Name          string   `json:"title"`
			DatePublished string   `json:"publishedDate"`
			ISBN          []ISBN   `json:"isbn"`
			PageCount     int      `json:"pageCount"`
			Authors       []string `json:"authors"`
		} `json:"volumeInfo"`
	} `json:"items"`
}
