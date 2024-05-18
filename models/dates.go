package models

type ArtistConcertDates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}
