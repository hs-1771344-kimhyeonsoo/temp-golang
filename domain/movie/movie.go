package movie

type Movie struct {
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
}

type Genre struct {
	Id   int
	Name string
}

type ProductionCompany struct {
	Id      int
	Name    string
	Country string
}
