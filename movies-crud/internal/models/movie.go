package models

type Movie struct {
	ID       string    `json:"id"`
	UID      string    `json:"uid"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
