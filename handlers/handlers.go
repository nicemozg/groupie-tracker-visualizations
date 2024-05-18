package handler

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/models"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

func GroupieTrackerPageHandler(w http.ResponseWriter, r *http.Request) {
	var tmplPath string
	// Определяем путь к шаблону в зависимости от запрошенного URL
	if r.Method == "GET" {
		if r.URL.Path == "/" {
			tmplPath = filepath.Join("web", "index.html")
		} else {
			// Указываем путь к файлу 404.html
			tmplPath = filepath.Join("web", "404.html")
			// Открываем файл 404.html
			file, err := ioutil.ReadFile(tmplPath)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Printf("Error reading file: %v", err)
				return
			}
			// Устанавливаем заголовок ответа для типа контента
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			// Отправляем статус 404 Not Found и содержимое файла 404.html
			w.WriteHeader(http.StatusNotFound)
			w.Write(file)
			return
		}
		tmpl, err := template.ParseFiles(tmplPath)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error parsing template: %v", err)
			return
		}
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error executing template: %v", err)
			return
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func AlbumListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Отправляем запрос к https://groupietrackers.herokuapp.com/api/artists
		resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
		if err != nil {
			http.Error(w, "Error fetching artists data", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		// Парсим JSON ответ
		var artists []struct {
			ID    int    `json:"id"`
			Image string `json:"image"`
			Name  string `json:"name"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
			http.Error(w, "Error decoding JSON", http.StatusInternalServerError)
			return
		}

		// Кодируем полученные данные в формат JSON и отправляем обратно на клиент
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(artists); err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func ArtistInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Получаем данные из тела запроса
		var requestData struct {
			ID string `json:"id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}

		// Формируем URL для запроса к API артиста и другим адресам
		urlForArtistInfo := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/artists/%s", requestData.ID)
		urlForArtistLocations := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/locations/%s", requestData.ID)
		urlForArtistConcertDates := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/dates/%s", requestData.ID)
		urlForRelations := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/relation/%s", requestData.ID)

		// Отправляем запросы к API артиста и другим адресам
		respArtistInfo, err := http.Get(urlForArtistInfo)
		if err != nil {
			http.Error(w, "Error fetching artist data", http.StatusInternalServerError)
			return
		}
		defer respArtistInfo.Body.Close()

		respArtistLocations, err := http.Get(urlForArtistLocations)
		if err != nil {
			http.Error(w, "Error fetching artist data", http.StatusInternalServerError)
			return
		}
		defer respArtistLocations.Body.Close()

		respArtistConcerts, err := http.Get(urlForArtistConcertDates)
		if err != nil {
			http.Error(w, "Error fetching artist data", http.StatusInternalServerError)
			return
		}
		defer respArtistConcerts.Body.Close()

		respRelations, err := http.Get(urlForRelations)
		if err != nil {
			http.Error(w, "Error fetching artist data", http.StatusInternalServerError)
			return
		}
		defer respRelations.Body.Close()

		// Создаем структуры для данных об артисте, локациях, концертах и отношениях
		var artistInfo models.ArtistInfo
		var artistLocations models.ArtistLocations
		var artistConcertsDates models.ArtistConcertDates
		var artistRelations models.ArtistDatesLocations

		// Декодируем JSON-данные и заполняем структуры
		if err := json.NewDecoder(respArtistInfo.Body).Decode(&artistInfo); err != nil {
			fmt.Println("Error respArtistInfo")
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}
		if err := json.NewDecoder(respArtistLocations.Body).Decode(&artistLocations); err != nil {
			fmt.Println("Error respArtistLocations")
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}
		if err := json.NewDecoder(respArtistConcerts.Body).Decode(&artistConcertsDates); err != nil {
			fmt.Println("Error respArtistConcerts")
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}
		if err := json.NewDecoder(respRelations.Body).Decode(&artistRelations); err != nil {
			fmt.Println("Error respRelations")
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}
		if artistInfo.Name == "" {
			http.Error(w, "Error IDi Nah", http.StatusBadRequest)
			return
		}

		// Собираем все данные в одну структуру
		response := models.ArtistData{
			ArtistInfo:           artistInfo,
			ArtistDates:          artistConcertsDates,
			ArtistLocations:      artistLocations,
			ArtistDatesLocations: artistRelations,
		}

		// Кодируем полученные данные в формат JSON и отправляем обратно клиенту
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
