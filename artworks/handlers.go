package artworks

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jcleira/golang-custom-handler/handler"
)

// ConfigureHandlers is meant to be called by the server.go main routine.
// It will configure the artworks package handlers, currently it configures the
// database connection the json schemas for validation and the router.
//
// r: The HTTP server *mux.Router to be configured.
// db: The database connection to use.
//
// Returns nothing.
func ConfigureHandlers(r *mux.Router, db *sql.DB) {
	artworksClient := &Client{
		DB: db,
	}

	r.Handle("/artworks", GetArtworksHandler(artworksClient)).Methods("GET")
	r.Handle("/artworks", AddArtworkHandler(artworksClient)).Methods("PUT", "OPTIONS")
	r.Handle("/artworks/{id:[0-9]+}", GetArtworkHandler(artworksClient)).Methods("GET")
	r.Handle("/artworks/{id:[0-9]+}", UpdateArtworkHandler(artworksClient)).Methods("PUT")
	r.Handle("/artworks/{id:[0-9]+}", DeleteArtworkHandler(artworksClient)).Methods("DELETE")
}

// GetArtworksHandler provides a HTTP endpoint to fetch all the Artworks.
//
// Response example:
// [{
//   ID: 1,
//   Rei: '#elle',
//   CreatedAt: 1489140631,
//   Pro: 'Ayuntamiento de Mahón',
//   Ubi: 'Desconocido',
//   ...
//   Tecnica: '',
//	},
//   ID: 2,
//   Rei: '#elle2',
//   CreatedAt: 1489140631,
//   Pro: 'Ayuntamiento de Mahón',
//   Ubi: 'Desconocido',
//   ...
//   Tecnica: '',
// }]
//
//
// artworksClient : The Artworks client either real or fake that implements the
//		  						 ArtworksController interface, a fake artworks client is used
//      						 for testing purposes.
//
// Returns a CustomHander ready to be added to a HTTP server / router.
func GetArtworksHandler(artworksClient ArtworksController) handler.CustomHandler {
	return func(w http.ResponseWriter, r *http.Request) *handler.HTTPError {
		artworks, err := artworksClient.GetArtworks()
		if err != nil {
			return &handler.HTTPError{err, http.StatusInternalServerError}
		}

		json.NewEncoder(w).Encode(artworks)
		return nil
	}
}

// AddArtworkHandler provides a HTTP endpoint to insert an Artwork information.
//
// artworksClient : The Artworks client either real or fake that implements the
//		  						 ArtworksController interface, a fake artworks client is used
//      						 for testing purposes.
//
// Returns a CustomHandler ready to be added to a HTTP server / router.
func AddArtworkHandler(artworksClient ArtworksController) handler.CustomHandler {
	return func(w http.ResponseWriter, r *http.Request) *handler.HTTPError {
		var artwork Artwork

		if err := json.NewDecoder(r.Body).Decode(&artwork); err != nil {
			return &handler.HTTPError{err, http.StatusBadRequest}
		}
		defer r.Body.Close()

		artwork.CreatedAt = time.Now().Unix()

		if err := artworksClient.AddUpdateArtwork("INSERT", &artwork); err != nil {
			return &handler.HTTPError{err, http.StatusInternalServerError}
		}

		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(artwork)
		return nil
	}
}

// GetArtworkHandler provides a HTTP endpoint to fetch a single Artwork.
//
// Response example:
// {
//   ID: 1,
//   Rei: '#elle',
//   CreatedAt: 1489140631,
//   Pro: 'Ayuntamiento de Mahón',
//   Ubi: 'Desconocido',
//   ...
//   Tecnica: '',
// }
//
// artworksClient : The Artworks client either real or fake that implements the
//		  						 ArtworksController interface, a fake artworks client is used
//      						 for testing purposes.
//
// Returns a CustomHander ready to be added to a HTTP server / router.
func GetArtworkHandler(artworksClient ArtworksController) handler.CustomHandler {
	return func(w http.ResponseWriter, r *http.Request) *handler.HTTPError {
		urlID, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			return &handler.HTTPError{
				errors.New("Unable to update Artwork URL ID mismatch body artwork ID"),
				http.StatusBadRequest,
			}
		}

		artwork, err := artworksClient.GetArtwork(urlID)
		if err != nil {
			return &handler.HTTPError{err, http.StatusInternalServerError}
		}

		json.NewEncoder(w).Encode(artwork)
		return nil
	}
}

// UpdateArtworkHandler provides a HTTP endpoint to update an Artwork
// information.
//
// artworksClient : The Artworks client either real or fake that implements the
//		  						 ArtworksController interface, a fake artworks client is used
//      						 for testing purposes.
//
// Returns a CustomHandler ready to be added to a HTTP server / router.
func UpdateArtworkHandler(artworksClient ArtworksController) handler.CustomHandler {
	return func(w http.ResponseWriter, r *http.Request) *handler.HTTPError {
		var artwork Artwork

		if err := json.NewDecoder(r.Body).Decode(&artwork); err != nil {
			return &handler.HTTPError{err, http.StatusBadRequest}
		}
		defer r.Body.Close()

		if urlID, _ := strconv.Atoi(mux.Vars(r)["id"]); urlID != artwork.ID {
			return &handler.HTTPError{
				errors.New("Unable to update Artwork URL ID mismatch body artwork ID"),
				http.StatusBadRequest,
			}
		}

		if err := artworksClient.AddUpdateArtwork("UPDATE", &artwork); err != nil {
			return &handler.HTTPError{err, http.StatusInternalServerError}
		}

		w.WriteHeader(http.StatusNoContent)
		return nil
	}
}

// DeleteArtworkHandler provides a HTTP endpoint to delete Artwork information
// by the given ID.
//
// artworksClient : The Artworks client either real or fake that implements the
//		  						 ArtworksController interface, a fake artworks client is used
//      						 for testing purposes.
//
// Returns a CustomHandler ready to be added to a HTTP server / router.
func DeleteArtworkHandler(artworksClient ArtworksController) handler.CustomHandler {
	return func(w http.ResponseWriter, r *http.Request) *handler.HTTPError {
		urlID, _ := strconv.Atoi(mux.Vars(r)["id"])

		if err := artworksClient.DeleteArtwork(urlID); err != nil {
			return &handler.HTTPError{err, http.StatusInternalServerError}
		}

		w.WriteHeader(http.StatusNoContent)
		return nil
	}
}
