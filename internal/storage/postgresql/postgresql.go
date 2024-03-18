package postgresql

import (
	"database/sql"
	"fmt"

	actor "github.com/P1coFly/vk_movies/internal/models/actor"
	movie "github.com/P1coFly/vk_movies/internal/models/movie"
	"github.com/P1coFly/vk_movies/internal/storage"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(URLPath string) (*Storage, error) {
	const op = "storage.postgresql.New"

	connStr := fmt.Sprintf("host=%s user=api_service password=12345678 dbname=VK_MOVIES sslmode=disable", URLPath)
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveActor(name, sex, birthday string) error {
	const op = "storage.postgresql.SaveActor"
	_, err := s.db.Exec(`INSERT INTO public."ACTORS" (name, sex, birthday) values ($1, $2, $3)`,
		name, sex, birthday)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) DeleteActorByID(actorID int64) error {
	const op = "storage.postgresql.DeleteActorByID"
	result, err := s.db.Exec(`DELETE FROM public."ACTORS" WHERE id = $1`, actorID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Проверка на количество удаленных записей
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrActorNotFound)
	}

	return nil
}

func (s *Storage) GetActors() ([]actor.Actor, error) {
	const op = "storage.postgresql.GetActors"
	actorsArr := []actor.Actor{}

	rows, err := s.db.Query(`SELECT 
		A.id AS actor_id,
    	A.name AS actor_name,
    	A.sex AS actor_sex,
    	A.birthday AS actor_birthday,
    	COALESCE(STRING_AGG(M.title, ', '), '') AS films
	FROM 
    	public."ACTORS" AS A
	LEFT JOIN
    	public."ACTORS_MOVIES" AS AM ON A.id = AM.actor_id
	LEFT JOIN
    	public."MOVIES" AS M ON AM.movie_id = M.id
	GROUP BY 
    	A.id, A.name, A.sex, A.birthday;`)

	if err != nil {
		return actorsArr, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		a := actor.Actor{}
		err := rows.Scan(&a.Id, &a.Name, &a.Sex, &a.Birthday, &a.Films)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		actorsArr = append(actorsArr, a)
	}

	return actorsArr, nil
}

func (s *Storage) UpdateActor(actorID int64, newName, newSex, newBirthday string) error {
	const op = "storage.postgresql.UpdateActor"

	// Проверяем, что актер с указанным идентификатором существует
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM public."ACTORS" WHERE id = $1`, actorID).Scan(&count)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if count == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrActorNotFound)
	}

	if newName != "" && newSex != "" && newBirthday != "" {
		_, err := s.db.Exec(`UPDATE public."ACTORS" SET name = $1, sex = $2, birthday = $3 WHERE id = $4`,
			newName, newSex, newBirthday, actorID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	}

	var name, sex, birthday string
	err = s.db.QueryRow(`SELECT name, sex, birthday FROM public."ACTORS" WHERE id = $1`, actorID).Scan(&name, &sex, &birthday)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if newName != "" {
		name = newName
	}
	if newSex != "" {
		sex = newSex
	}
	if newBirthday != "" {
		birthday = newBirthday
	}

	_, err = s.db.Exec(`UPDATE public."ACTORS" SET name = $1, sex = $2, birthday = $3 WHERE id = $4`,
		name, sex, birthday, actorID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) SaveMovie(m movie.Movie, actorIDs []int) error {
	const op = "storage.postgresql.SaveMovie"
	var movieID int

	s.db.QueryRow(`INSERT INTO public."MOVIES" (title, description, date_of_issue, rating) VALUES ($1, $2, $3, $4) returning id`,
		m.Title, m.Description, m.DateOfIssue, m.Rating).Scan(&movieID)
	fmt.Println(movieID)
	for _, actorID := range actorIDs {
		_, err := s.db.Exec(`INSERT INTO public."ACTORS_MOVIES" (actor_id, movie_id) VALUES ($1, $2)`,
			actorID, movieID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	return nil
}

func (s *Storage) DeleteMovieByID(actorID int64) error {
	const op = "storage.postgresql.DeleteMovieByID"
	result, err := s.db.Exec(`DELETE FROM public."MOVIES" WHERE id = $1`, actorID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Проверка на количество удаленных записей
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrActorNotFound)
	}

	return nil
}

func (s *Storage) FindMoviesByTitleFragment(titleFragment string) ([]movie.Movie, error) {
	const op = "storage.postgresql.FindMoviesByTitleFragment"
	rows, err := s.db.Query(`SELECT id, title, description, date_of_issue, rating FROM public."MOVIES" WHERE title ILIKE '%' || $1 || '%'`, titleFragment)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var movies []movie.Movie
	for rows.Next() {
		var movie movie.Movie
		if err := rows.Scan(&movie.Id, &movie.Title, &movie.Description, &movie.DateOfIssue, &movie.Rating); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		movies = append(movies, movie)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return movies, nil
}

func (s *Storage) FindMoviesByActorNameFragment(actorNameFragment string) ([]movie.Movie, error) {
	const op = "storage.postgresql.FindMoviesByActorNameFragment"
	rows, err := s.db.Query(`
		SELECT m.id, m.title, m.description, m.date_of_issue, m.rating
		FROM public."MOVIES" m
		JOIN public."ACTORS_MOVIES" am ON m.id = am.movie_id
		JOIN public."ACTORS" a ON am.actor_id = a.id
		WHERE a.name ILIKE '%' || $1 || '%'
	`, actorNameFragment)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var movies []movie.Movie
	for rows.Next() {
		var movie movie.Movie
		if err := rows.Scan(&movie.Id, &movie.Title, &movie.Description, &movie.DateOfIssue, &movie.Rating); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		movies = append(movies, movie)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return movies, nil
}

func (s *Storage) UpdateMovie(movieID int64, newTitle, newDescription string, newDateOfIssue string, newRating float64) error {
	const op = "storage.postgresql.UpdateMovie"

	// Проверяем, что фильм с указанным идентификатором существует
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM public."MOVIES" WHERE id = $1`, movieID).Scan(&count)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if count == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrMovieNotFound)
	}

	if newTitle != "" && newDescription != "" && newDateOfIssue != "" && newRating != 0 {
		_, err := s.db.Exec(`UPDATE public."MOVIES" SET title = $1, description = $2, date_of_issue = $3, rating = $4 WHERE id = $5`,
			newTitle, newDescription, newDateOfIssue, newRating, movieID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	}

	var title, description, dateOfIssue string
	var rating float64
	err = s.db.QueryRow(`SELECT title, description, date_of_issue, rating FROM public."MOVIES" WHERE id = $1`, movieID).Scan(&title, &description, &dateOfIssue, &rating)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if newTitle != "" {
		title = newTitle
	}
	if newDescription != "" {
		description = newDescription
	}
	if newDateOfIssue != "" {
		dateOfIssue = newDateOfIssue
	}
	if newRating != 0 {
		rating = newRating
	}

	_, err = s.db.Exec(`UPDATE public."MOVIES" SET title = $1, description = $2, date_of_issue = $3, rating = $4 WHERE id = $5`,
		title, description, dateOfIssue, rating, movieID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetSortedMovies(column, order string) ([]movie.Movie, error) {
	const op = "storage.postgresql.GetSortedMovies"

	// Проверяем, что указанный столбец существует
	validColumns := map[string]bool{
		"id":            true,
		"title":         true,
		"rating":        true,
		"date_of_issue": true,
	}

	if !validColumns[column] {
		return nil, fmt.Errorf("%s: invalid column to sort", op)
	}

	// Проверяем, что указанный порядок сортировки корректен
	validOrders := map[string]bool{
		"ASC":  true,
		"DESC": true,
	}

	if !validOrders[order] {
		return nil, fmt.Errorf("%s: incorrect sort order", op)
	}

	// Формируем запрос с учетом указанных параметров сортировки
	query := fmt.Sprintf("SELECT id, title, description, date_of_issue, rating FROM public.\"MOVIES\" ORDER BY %s %s", column, order)

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var movies []movie.Movie
	for rows.Next() {
		var movie movie.Movie
		if err := rows.Scan(&movie.Id, &movie.Title, &movie.Description, &movie.DateOfIssue, &movie.Rating); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		movies = append(movies, movie)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return movies, nil
}
