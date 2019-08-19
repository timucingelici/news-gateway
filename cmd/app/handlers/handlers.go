package handlers

import (
	"encoding/json"
	"errors"
	"github.com/badoux/checkmail"
	"github.com/go-chi/chi"
	"github.com/timucingelici/news-gateway/pkg/store"
	"log"
	"net/http"
	"sort"
	"time"
)

type RestHandler struct {
	store store.DataStore
}

func SetupRoutes(dataStore store.DataStore) http.Handler {

	rh := &RestHandler{store: dataStore}
	r := chi.NewRouter()

	r.Get("/news", rh.GetAllNews)
	r.Get("/news/{newsID}", rh.GetNews)
	r.Post("/news/{newsID}/share", rh.ShareNews)

	return r

}

func (rh *RestHandler) GetAllNews(w http.ResponseWriter, r *http.Request) {

	var newsList []News

	resp, err := rh.store.GetAllWithValues("*")

	if err != nil {
		log.Println("Failed to get news from the store. Err : ", err)
		rh.send500(w)
		return
	}

	// Check if any filtering is requested
	provider := r.URL.Query().Get("provider")
	category := r.URL.Query().Get("category")

	for key, v := range resp {

		var n News

		if err := json.Unmarshal([]byte(v), &n); err != nil {
			log.Printf("Failed to marshall the news for key %s. Err : %s\n", key, err)
		}

		if provider != "" && n.Provider != provider {
			continue
		}

		if category != "" && n.Category != category {
			continue
		}

		n.ID = key
		newsList = append(newsList, n)
	}

	sort.Sort(byNewsDate(newsList))

	rh.sendResponse(newsList, w, http.StatusOK)
}

func (rh *RestHandler) GetNews(w http.ResponseWriter, r *http.Request) {

	var news News

	newsID := chi.URLParam(r, "newsID")

	if newsID == "" {
		rh.send404(w)
		return
	}

	resp, err := rh.store.Get(newsID)

	if err != nil {
		log.Println("Failed to get news from the store. Err : ", err)
		rh.send500(w)
		return
	}

	if resp == "" {
		log.Println("Failed to find news with this ID : ", newsID)
		rh.send404(w)
		return
	}

	if err := json.Unmarshal([]byte(resp), &news); err != nil {
		log.Printf("Failed to marshall the news for key %s. Err : %s\n", newsID, err)
	}

	news.ID = newsID
	rh.sendResponse(news, w, http.StatusOK)
}

func (rh *RestHandler) ShareNews(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		rh.send400(w)
		return
	}

	email := r.PostFormValue("email")

	if email == "" {
		rh.sendResponse(&errResponse{"email field must be sent to share this article"}, w, 400)
		return
	}

	if err := checkmail.ValidateFormat(email); err != nil {
		rh.sendResponse(&errResponse{"email field must be a valid email"}, w, 400)
		return
	}

	newsID := chi.URLParam(r, "newsID")

	if newsID == "" {
		rh.send404(w)
		return
	}

	resp, err := rh.store.Get(newsID)

	if err != nil {
		log.Println("Failed to get news from the store. Err : ", err)
		rh.send500(w)
		return
	}

	if resp == "" {
		log.Println("Failed to find news with this ID : ", newsID)
		rh.send404(w)
		return
	}

	// Assume some magic happens and email has been sent.

	rh.sendResponse(&successResponse{"thanks for sharing someone else's email without asking them"}, w, http.StatusOK)
}

// Helper functions

func (rh *RestHandler) send400(w http.ResponseWriter) {
	err := &errResponse{errors.New("unable to process your request").Error()}
	rh.sendResponse(err, w, http.StatusBadRequest)
}
func (rh *RestHandler) send404(w http.ResponseWriter) {
	err := &errResponse{errors.New("the resource you're looking for does not exist").Error()}
	rh.sendResponse(err, w, http.StatusNotFound)
}

func (rh *RestHandler) send500(w http.ResponseWriter) {
	err := &errResponse{errors.New("oops. something went wrong").Error()}
	rh.sendResponse(err, w, http.StatusInternalServerError)
}

func (rh *RestHandler) sendResponse(response interface{}, w http.ResponseWriter, responseCode int) {

	bytes, err := json.Marshal(response)

	if err != nil {
		log.Println("Failed to encode the response. Err: ", err)
		w.WriteHeader(http.StatusInternalServerError)

		_, err := w.Write([]byte(err.Error()))

		if err != nil {
			log.Println("Failed to return error response. Err: ", err)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseCode)

	_, err = w.Write(bytes)

	if err != nil {
		log.Println("Failed to encode the response. Err: ", err)
		w.WriteHeader(http.StatusInternalServerError)

		_, err := w.Write([]byte(err.Error()))

		if err != nil {
			log.Println("Failed to return error response")
		}
	}
}

// Types

type News struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Description string    `json:"description"`
	DateTime    time.Time `json:"datetime"`
	Provider    string    `json:"provider"`
	Category    string    `json:"category"`
	Thumbnail   string    `json:"thumbnail"`
}

// Custom sort for the news.
// I cheated here a bit because sort.Reverse acting weird.
type byNewsDate []News

func (n byNewsDate) Len() int           { return len(n) }
func (n byNewsDate) Less(i, j int) bool { return n[i].DateTime.Unix() > n[j].DateTime.Unix() }
func (n byNewsDate) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }

type errResponse struct {
	Error string `json:"error"`
}

type successResponse struct {
	Message string `json:"message"`
}
