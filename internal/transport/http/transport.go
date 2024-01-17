package http

import (
	"encoding/json"
	"net/http"
"log"
	"github.com/gorilla/mux"
	"albumproject/internal/album"
	"albumproject/internal/endpoint"  
)

// makeHTTPHandlers configures HTTP handlers for the album service
// func MakeHTTPHandlers(router *mux.Router, service album.AlbumService) {
// 	router.HandleFunc("/addalbums", createAlbumHandler(service)).Methods("POST")
// 	router.HandleFunc("/getbyIdalbums/{id}", getAlbumHandler(service)).Methods("GET")
// 	router.HandleFunc("/updatealbums/{id}", updateAlbumHandler(service)).Methods("PUT")
// 	router.HandleFunc("/deletealbums/{id}", deleteAlbumHandler(service)).Methods("DELETE")
// 	router.HandleFunc("/getSongsWithCapitalTitles", getSongsWithCapitalTitlesHandler(service)).Methods("GET")
//     router.HandleFunc("/albums/search", searchAlbumsHandler(service)).Methods("GET")
//}
func MakeHTTPHandlers(router *mux.Router, service album.AlbumService) {
    router.HandleFunc(endpoint.CreateAlbumEndpoint, createAlbumHandler(service)).Methods("POST")
    router.HandleFunc(endpoint.GetAlbumEndpoint, getAlbumHandler(service)).Methods("GET")
    router.HandleFunc(endpoint.UpdateAlbumEndpoint, updateAlbumHandler(service)).Methods("PUT")
    router.HandleFunc(endpoint.DeleteAlbumEndpoint, deleteAlbumHandler(service)).Methods("DELETE")
    router.HandleFunc(endpoint.SearchAlbumsEndpoint, searchAlbumsHandler(service)).Methods("GET")
}

func createAlbumHandler(service album.AlbumService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var album album.Album
		err := json.NewDecoder(r.Body).Decode(&album)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err = service.CreateAlbum(r.Context(), album)
		if err != nil {
			http.Error(w, "Failed to create album", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
func getAlbumHandler(service album.AlbumService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        albumID := vars["id"]

        log.Printf("Attempting http to retrieve album with ID: %s", albumID)

        // Call your service to get the album by ID
        album, err := service.GetAlbum(r.Context(), albumID)
        if err != nil {
            log.Printf("Album with ID %s not found", albumID)
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }

        // Serialize the album to JSON and send it in the response
        json.NewEncoder(w).Encode(album)
    }
}
	
func updateAlbumHandler(service album.AlbumService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		albumID := vars["id"]

		var update album.Album
		err := json.NewDecoder(r.Body).Decode(&update)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err = service.UpdateAlbum(r.Context(), albumID, update)
		if err != nil {
			http.Error(w, "Failed to update album", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func deleteAlbumHandler(service album.AlbumService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		albumID := vars["id"]

		err := service.DeleteAlbum(r.Context(), albumID)
		if err != nil {
			http.Error(w, "Failed to delete album", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
func getSongsWithCapitalTitlesHandler(service album.AlbumService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Attempting to retrieve songs with capital titles")

		// Call your service to get songs with capital titles
		songs, err := service.GetSongsWithCapitalTitles(r.Context())
		if err != nil {
			log.Printf("Error retrieving songs with capital titles: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Serialize the songs to JSON
		response, err := json.Marshal(songs)
		if err != nil {
			log.Printf("Error serializing songs to JSON: %s", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Set the Content-Type header to application/json
		w.Header().Set("Content-Type", "application/json")

		// Write the JSON response
		w.Write(response)
	}
}
func searchAlbumsHandler(service album.AlbumService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        searchTerm := r.URL.Query().Get("q")
        if searchTerm == "" {
            http.Error(w, "Search term is required", http.StatusBadRequest)
            return
        }

        albums, err := service.SearchAlbums(r.Context(), searchTerm)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(albums)
    }
}