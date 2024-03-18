package postgresql_test

import (
	"database/sql"
	"testing"

	"github.com/P1coFly/vk_movies/internal/models/actor"
	"github.com/P1coFly/vk_movies/internal/models/movie"
	"github.com/P1coFly/vk_movies/internal/storage/postgresql"
	"github.com/stretchr/testify/assert"
)

func deleteLastActor(storage *postgresql.Storage) error {
	actors, err := storage.GetActors()
	if err != nil {
		return err
	}
	if len(actors) == 0 {
		return nil
	}
	lastActor := actors[len(actors)-1]
	return storage.DeleteActorByID(lastActor.Id)
}

func deleteLastMovie(storage *postgresql.Storage) error {
	movies, err := storage.GetSortedMovies("id", "ASC")
	if err != nil {
		return err
	}
	if len(movies) == 0 {
		return nil
	}
	lastMovie := movies[len(movies)-1]
	return storage.DeleteMovieByID(lastMovie.Id)
}

func TestActor(t *testing.T) {
	storage, err := postgresql.New("localhost")
	if err != nil {
		t.Fatal("Error initializing storage:", err)
	}

	testActor := actor.Actor{Name: "TestActor", Sex: "M", Birthday: "2000-01-01"}
	err = storage.SaveActor(testActor.Name, testActor.Sex, testActor.Birthday)
	if err != nil {
		t.Fatal("Error saving actor:", err)
	}

	actors, err := storage.GetActors()
	if err != nil {
		t.Fatal("Error retrieving actors from database:", err)
	}

	var found bool
	for _, a := range actors {
		if a.Name == testActor.Name {
			found = true
			break
		}
	}
	err = deleteLastActor(storage)
	if err != nil && err != sql.ErrNoRows {
		t.Fatal("Error deleting last actor before test:", err)
	}

	assert.True(t, found, "Test actor not found in database")
}

func TestMovie(t *testing.T) {
	storage, err := postgresql.New("localhost")
	if err != nil {
		t.Fatal("Error initializing storage:", err)
	}

	testMovie := movie.Movie{Title: "TestMovie", Description: "TestDescription", DateOfIssue: "2000-01-01", Rating: 7.5}
	err = storage.SaveMovie(testMovie, []int{1, 2})
	if err != nil {
		t.Fatal("Error saving movie:", err)
	}

	movies, err := storage.GetSortedMovies("title", "ASC")
	if err != nil {
		t.Fatal("Error retrieving movies from database:", err)
	}

	var found bool
	for _, m := range movies {
		if m.Title == testMovie.Title {
			found = true
			break
		}
	}

	err = deleteLastMovie(storage)
	if err != nil && err != sql.ErrNoRows {
		t.Fatal("Error deleting last movie before test:", err)
	}

	assert.True(t, found, "Test movie not found in database")
}

func TestUpdateActor(t *testing.T) {
	storage, err := postgresql.New("localhost")
	if err != nil {
		t.Fatal("Error initializing storage:", err)
	}

	// Сохранение актера для обновления
	testActor := actor.Actor{Name: "TestActor", Sex: "M", Birthday: "2000-01-01"}
	err = storage.SaveActor(testActor.Name, testActor.Sex, testActor.Birthday)
	if err != nil {
		t.Fatal("Error saving actor for update:", err)
	}

	// Получение списка актеров для обновления
	actorsBeforeUpdate, err := storage.GetActors()
	if err != nil {
		t.Fatal("Error retrieving actors before update:", err)
	}

	// Выбор актера для обновления (последний добавленный)
	actorToUpdate := actorsBeforeUpdate[len(actorsBeforeUpdate)-1]

	// Обновление данных актера
	newName := "UpdatedName"
	newSex := "F"
	newBirthday := "1990-05-05T00:00:00Z"
	err = storage.UpdateActor(actorToUpdate.Id, newName, newSex, newBirthday)
	if err != nil {
		t.Fatal("Error updating actor:", err)
	}

	// Получение списка актеров после обновления
	actorsAfterUpdate, err := storage.GetActors()
	if err != nil {
		t.Fatal("Error retrieving actors after update:", err)
	}

	// Проверка обновленных данных актера
	var found bool
	a := actorsAfterUpdate[len(actorsBeforeUpdate)-1]
	if a.Id == actorToUpdate.Id && a.Name == newName && a.Sex == newSex && a.Birthday == newBirthday {
		found = true
	}

	assert.True(t, found, "Updated actor not found in database")

	// Удаление созданного актера после теста
	err = storage.DeleteActorByID(actorToUpdate.Id)
	if err != nil {
		t.Fatal("Error deleting actor after test:", err)
	}
}

func TestUpdateMovie(t *testing.T) {
	storage, err := postgresql.New("localhost")
	if err != nil {
		t.Fatal("Error initializing storage:", err)
	}

	// Сохранение фильма для обновления
	testMovie := movie.Movie{Title: "TestMovie", Description: "TestDescription", DateOfIssue: "2000-01-01", Rating: 7.5}
	err = storage.SaveMovie(testMovie, []int{1, 2})
	if err != nil {
		t.Fatal("Error saving movie for update:", err)
	}

	// Получение списка фильмов для обновления
	moviesBeforeUpdate, err := storage.GetSortedMovies("id", "ASC")
	if err != nil {
		t.Fatal("Error retrieving movies before update:", err)
	}

	// Выбор фильма для обновления (последний добавленный)
	movieToUpdate := moviesBeforeUpdate[len(moviesBeforeUpdate)-1]

	// Обновление данных фильма
	newTitle := "UpdatedTitle"
	newDescription := "UpdatedDescription"
	newDateOfIssue := "2020-12-31T00:00:00Z"
	var newRating float64 = 8.0
	err = storage.UpdateMovie(movieToUpdate.Id, newTitle, newDescription, newDateOfIssue, newRating)
	if err != nil {
		t.Fatal("Error updating movie:", err)
	}

	// Получение списка фильмов после обновления
	moviesAfterUpdate, err := storage.GetSortedMovies("id", "ASC")
	if err != nil {
		t.Fatal("Error retrieving movies after update:", err)
	}

	// Проверка обновленных данных фильма
	var found bool
	m := moviesAfterUpdate[len(moviesBeforeUpdate)-1]

	if m.Id == movieToUpdate.Id && m.Title == newTitle && m.Description == newDescription && m.DateOfIssue == newDateOfIssue && m.Rating == newRating {
		found = true
	}

	assert.True(t, found, "Updated movie not found in database")

	// Удаление созданного фильма после теста
	err = storage.DeleteMovieByID(movieToUpdate.Id)
	if err != nil {
		t.Fatal("Error deleting movie after test:", err)
	}
}
