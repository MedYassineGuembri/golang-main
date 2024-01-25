package main

import (
	"errors"
	"estiam/dictionary"
	"estiam/logger"
	"estiam/middleware"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)



func main() {
	dict := dictionary.New("dictionary.json")

	router := mux.NewRouter()
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.AuthMiddleware)

	router.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		actionAdd(w, r, dict)
	}).Methods("POST")

	router.HandleFunc("/define/{word}", func(w http.ResponseWriter, r *http.Request) {
		actionDefine(w, r, dict)
	}).Methods("GET")

	router.HandleFunc("/remove/{word}", func(w http.ResponseWriter, r *http.Request) {
		actionRemove(w, r, dict)
	}).Methods("DELETE")

	// Add a new route for listing all entries
	router.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		actionList(w, dict)
	}).Methods("GET")

	http.Handle("/", router)
	fmt.Println("Server is running on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func actionAdd(w http.ResponseWriter, r *http.Request, d *dictionary.Dictionary) {
	word := r.FormValue("word")
	definition := r.FormValue("definition")

	if err := validateData(word, definition); err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := d.Add(word, definition)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Entrée ajoutée: %s - %s\n", word, definition)
}

func actionDefine(w http.ResponseWriter, r *http.Request, d *dictionary.Dictionary) {
	vars := mux.Vars(r)
	word := vars["word"]
	entry, found, err := d.Get(word)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !found {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Définition de '%s': %s\n", word, entry.Definition)
}

func actionRemove(w http.ResponseWriter, r *http.Request, d *dictionary.Dictionary) {
	vars := mux.Vars(r)
	word := vars["word"]
	err := d.Remove(word)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Entrée supprimée: %s\n", word)
}

func actionList(w http.ResponseWriter, d *dictionary.Dictionary) {
	words, err := d.List()
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, word := range words {
		entry, found, _ := d.Get(word)
		if found {
			fmt.Fprintf(w, "%v: %v\n", word, entry.Definition)
		}
	}
}

func validateData(word, definition string) error {
    if len(word) == 0 || len(definition) == 0 {
        return errors.New("le mot et la définition ne peuvent pas être vides")
    }
    if len(word) > 50 || len(definition) > 200 {
        return errors.New("longueur de mot ou de définition dépassée")
    }
    return nil
}

func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
    logger.Logger.Printf("Error: %s, StatusCode: %d\n", message, statusCode)
    w.WriteHeader(statusCode)
    fmt.Fprintf(w, "Erreur: %s\n", message)
}