package movie

type MovieWithReview struct {
	Id                  int
	Adult               bool
	Genres              []Genre
	Title               string
	Language            string
	Overview            string
	Poster              string
	ProductionCompanies []ProductionCompany
	ReleaseDate         string
	Revenue             int
	Runtime             int
	Tagline             string
	Rating              float32
	Votes               int
	UserId              int
	Liked               bool
}
