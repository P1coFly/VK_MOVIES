package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/P1coFly/vk_movies/internal/config"
	"github.com/P1coFly/vk_movies/internal/models/actor"
	"github.com/P1coFly/vk_movies/internal/models/movie"
	"github.com/P1coFly/vk_movies/internal/storage"
	"github.com/P1coFly/vk_movies/internal/storage/postgresql"
)

// @Summary Получение списка актеров
// @Description Получение списка всех актеров из базы данных
// @Tags Actors
// @ID getActors
// @Produce json
// @Success 200 {array} actor.Actor
// @Failure 500 {object} ErrorResponse
// @Router /api/actors [get]
func ActorsHandler(s *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		actors, err := s.GetActors()
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка при получении списка актеров: %s", err), http.StatusInternalServerError)
			return
		}

		// Отправляем ответ в формате JSON
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(actors); err != nil {
			http.Error(w, fmt.Sprintf("Ошибка при кодировании ответа в JSON: %s", err), http.StatusInternalServerError)
			return
		}

		// Успешный ответ
		w.WriteHeader(http.StatusOK)
	}
}

// @Summary Управление актерами
// @Description Создание, обновление и удаление актеров
// @Tags Actor
// @ID manageActors
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /api/actor [post]
// @Router /api/actor [patch]
// @Router /api/actor [delete]
func ActorHandler(s *postgresql.Storage, cfg *config.Config) http.HandlerFunc {
	return AuthenticatedHandler(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			saveActorHandler(s, w, r)
		case http.MethodPatch:
			updateActorHandler(s, w, r)
		case http.MethodDelete:
			deleteActorHandler(s, w, r)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}, cfg)
}

// @Summary Создание актера
// @Description Создание нового актера в базе данных
// @Tags Actor
// @ID createActor
// @Accept json
// @Produce json
// @Success 201 "Actor created successfully"
// @Failure 400 {object} ErrorResponse
// @Router /api/actor [post]
func saveActorHandler(s *postgresql.Storage, w http.ResponseWriter, r *http.Request) {
	const op = "saveActorHandler"

	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Чтение данных из тела запроса
	var actor actor.Actor
	if err := json.NewDecoder(r.Body).Decode(&actor); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при декодировании JSON: %s", err), http.StatusBadRequest)
		return
	}

	// Сохранение актера в базе данных
	if err := s.SaveActor(actor.Name, actor.Sex, actor.Birthday); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при сохранении актера: %s", err), http.StatusInternalServerError)
		return
	}

	// Успешный ответ
	w.WriteHeader(http.StatusCreated)
}

// @Summary Обновление актера
// @Description Обновление существующего актера в базе данных
// @Tags Actor
// @ID updateActor
// @Accept json
// @Produce json
// @Param actorID query string true "ID актера"
// @Success 200 "Actor updated successfully"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/actor [patch]
func updateActorHandler(s *postgresql.Storage, w http.ResponseWriter, r *http.Request) {
	const op = "updateActorHandler"

	if r.Method != http.MethodPatch {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Парсинг ID актера из URL
	actorIDStr := r.URL.Query().Get("actorID")
	actorID, err := strconv.ParseInt(actorIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Неверный формат ID актера", http.StatusBadRequest)
		return
	}

	// Чтение данных из тела запроса
	var actor actor.Actor
	if err := json.NewDecoder(r.Body).Decode(&actor); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при декодировании JSON: %s", err), http.StatusBadRequest)
		return
	}

	// Обновление актера в базе данных
	if err := s.UpdateActor(actorID, actor.Name, actor.Sex, actor.Birthday); err != nil {
		if errors.Is(err, storage.ErrActorNotFound) {
			http.Error(w, "Актер не найден", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Ошибка при обновлении актера: %s", err), http.StatusInternalServerError)
		}
		return
	}

	// Успешный ответ
	w.WriteHeader(http.StatusOK)
}

// @Summary Удаление актера
// @Description Удаление актера из базы данных
// @Tags Actor
// @ID deleteActor
// @Param actorID query string true "ID актера"
// @Success 200 "Actor deleted successfully"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/actor [delete]
func deleteActorHandler(s *postgresql.Storage, w http.ResponseWriter, r *http.Request) {
	const op = "deleteActorHandler"

	if r.Method != http.MethodDelete {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Парсинг ID актера из URL
	actorIDStr := r.URL.Query().Get("actorID")
	actorID, err := strconv.ParseInt(actorIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Неверный формат ID актера", http.StatusBadRequest)
		return
	}

	// Удаление актера из базы данных
	if err := s.DeleteActorByID(actorID); err != nil {
		if errors.Is(err, storage.ErrActorNotFound) {
			http.Error(w, "Актер не найден", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Ошибка при удалении актера: %s", err), http.StatusInternalServerError)
		}
		return
	}

	// Успешный ответ
	w.WriteHeader(http.StatusOK)
}

// @Summary Управление фильмами
// @Description Создание, обновление и удаление фильмов
// @Tags Movie
// @ID manageMovies
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /api/movie [post]
// @Router /api/movie [patch]
// @Router /api/movie [delete]
func MovieHandler(s *postgresql.Storage, cfg *config.Config) http.HandlerFunc {
	return AuthenticatedHandler(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			saveMovieHandler(s, w, r)
		case http.MethodPatch:
			updateMovieHandler(s, w, r)
		case http.MethodDelete:
			deleteMovieHandler(s, w, r)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}, cfg)
}

// @Summary Создание фильма
// @Description Создание нового фильма в базе данных
// @Tags Movie
// @ID createMovie
// @Accept json
// @Produce json
// @Success 201 "Movie created successfully"
// @Failure 400 {object} ErrorResponse
// @Router /api/movie [post]
func saveMovieHandler(s *postgresql.Storage, w http.ResponseWriter, r *http.Request) {
	const op = "saveMovieHandler"

	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Чтение данных из тела запроса
	var m movie.Movie
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при декодировании JSON: %s", err), http.StatusBadRequest)
		return
	}

	// Чтение списка ID актеров из параметров запроса
	actorIDs, err := readActorIDsFromRequest(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при чтении ID актеров: %s", err), http.StatusBadRequest)
		return
	}

	// Сохранение фильма в базе данных
	fmt.Println(actorIDs)
	if err := s.SaveMovie(m, actorIDs); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при сохранении фильма: %s", err), http.StatusInternalServerError)
		return
	}

	// Успешный ответ
	w.WriteHeader(http.StatusCreated)
}

// @Summary Обновление фильма
// @Description Обновление существующего фильма в базе данных
// @Tags Movie
// @ID updateMovie
// @Accept json
// @Produce json
// @Param movieID query string true "ID фильма"
// @Success 200 "Movie updated successfully"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/movie [patch]
func updateMovieHandler(s *postgresql.Storage, w http.ResponseWriter, r *http.Request) {
	const op = "updateMovieHandler"

	if r.Method != http.MethodPatch {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Парсинг ID фильма из параметров запроса
	movieIDStr := r.URL.Query().Get("movieID")
	movieID, err := strconv.ParseInt(movieIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Неверный формат ID фильма", http.StatusBadRequest)
		return
	}

	// Чтение данных из тела запроса
	var m movie.Movie
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при декодировании JSON: %s", err), http.StatusBadRequest)
		return
	}

	// Обновление фильма в базе данных
	if err := s.UpdateMovie(movieID, m.Title, m.Description, m.DateOfIssue, float64(m.Rating)); err != nil {
		if errors.Is(err, storage.ErrMovieNotFound) {
			http.Error(w, "Фильм не найден", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Ошибка при обновлении фильма: %s", err), http.StatusInternalServerError)
		}
		return
	}

	// Успешный ответ
	w.WriteHeader(http.StatusOK)
}

// @Summary Удаление фильма
// @Description Удаление фильма из базы данных
// @Tags Movie
// @ID deleteMovie
// @Param movieID query string true "ID фильма"
// @Success 200 "Movie deleted successfully"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/movie [delete]
func deleteMovieHandler(s *postgresql.Storage, w http.ResponseWriter, r *http.Request) {
	const op = "deleteMovieHandler"

	if r.Method != http.MethodDelete {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Парсинг ID фильма из параметров запроса
	movieIDStr := r.URL.Query().Get("movieID")
	movieID, err := strconv.ParseInt(movieIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Неверный формат ID фильма", http.StatusBadRequest)
		return
	}

	// Удаление фильма из базы данных
	if err := s.DeleteMovieByID(movieID); err != nil {
		if errors.Is(err, storage.ErrMovieNotFound) {
			http.Error(w, "Фильм не найден", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Ошибка при удалении фильма: %s", err), http.StatusInternalServerError)
		}
		return
	}

	// Успешный ответ
	w.WriteHeader(http.StatusOK)
}

// @Summary Поиск фильмов по фрагменту названия
// @Description Поиск фильмов по части названия в базе данных
// @Tags Movies
// @ID findMoviesByTitleFragment
// @Produce json
// @Param titleFragment query string true "Фрагмент названия фильма"
// @Success 200 {array} movie.Movie
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/moviesByTitle [get]
func FindMoviesByTitleFragmentHandler(s *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "findMoviesByTitleFragmentHandler"

		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		// Получение фрагмента названия фильма из параметров запроса
		titleFragment := r.URL.Query().Get("titleFragment")
		if titleFragment == "" {
			http.Error(w, "Необходимо указать фрагмент названия фильма", http.StatusBadRequest)
			return
		}

		// Поиск фильмов по фрагменту названия в хранилище
		movies, err := s.FindMoviesByTitleFragment(titleFragment)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка при поиске фильмов: %s", err), http.StatusInternalServerError)
			return
		}

		// Отправка ответа в формате JSON
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(movies); err != nil {
			http.Error(w, fmt.Sprintf("Ошибка при кодировании JSON: %s", err), http.StatusInternalServerError)
			return
		}
	}
}

// @Summary Поиск фильмов по фрагменту имени актера
// @Description Поиск фильмов по части имени актера в базе данных
// @Tags Movies
// @ID findMoviesByActorNameFragment
// @Produce json
// @Param actorNameFragment query string true "Фрагмент имени актера"
// @Success 200 {array} movie.Movie
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/moviesByActorName [get]
func FindMoviesByActorNameFragmentHandler(s *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "findMoviesByActorNameFragmentHandler"

		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		// Получение фрагмента имени актера из параметров запроса
		actorNameFragment := r.URL.Query().Get("actorNameFragment")
		if actorNameFragment == "" {
			http.Error(w, "Необходимо указать фрагмент имени актера", http.StatusBadRequest)
			return
		}

		// Поиск фильмов по фрагменту имени актера в хранилище
		movies, err := s.FindMoviesByActorNameFragment(actorNameFragment)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка при поиске фильмов: %s", err), http.StatusInternalServerError)
			return
		}

		// Отправка ответа в формате JSON
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(movies); err != nil {
			http.Error(w, fmt.Sprintf("Ошибка при кодировании JSON: %s", err), http.StatusInternalServerError)
			return
		}
	}
}

// ErrorResponse represents an error response structure.
type ErrorResponse struct {
	Error string `json:"error"`
}

// Вспомогательная функция для чтения списка ID актеров из параметров запроса
func readActorIDsFromRequest(r *http.Request) ([]int, error) {
	actorIDsStr := r.URL.Query().Get("actorIDs")
	if actorIDsStr == "" {
		return nil, nil
	}
	actorIDsStrList := strings.Split(actorIDsStr, ",")
	actorIDs := make([]int, len(actorIDsStrList))
	for i, idStr := range actorIDsStrList {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, fmt.Errorf("неверный формат ID актера: %s", idStr)
		}
		actorIDs[i] = id
	}
	return actorIDs, nil
}

// AuthenticatedHandler is a wrapper function to authenticate requests before passing them to the actual handler.
func AuthenticatedHandler(next http.HandlerFunc, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")
		expectedToken := cfg.AuthToken
		if authToken != expectedToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Если аутентификация прошла успешно, передаем запрос дальше.
		next.ServeHTTP(w, r)
	}
}
