package movie

import (
	"fmt"
)

type Movie struct {
	Id          int64   `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	DateOfIssue string  `json:"date_of_issue"`
	Rating      float64 `json:"rating"`
}

func New(title, description, dateOfIssue string, rating float64) (*Movie, error) {
	const op = "models.movie.New"

	if len(title) < 1 || len(title) >= 150 {
		return nil, fmt.Errorf("%s: the title length must be greater than 1 and not greater than 150", op)
	}
	if len(description) >= 1000 {
		return nil, fmt.Errorf("%s: the description length must be not greater than 1000", op)
	}
	if rating < 0.0 || rating > 10.0 {
		return nil, fmt.Errorf("%s: the value must be in the range from 0.0 to 10.0", op)
	}

	return &Movie{Title: title, Description: description, DateOfIssue: dateOfIssue, Rating: rating}, nil
}
