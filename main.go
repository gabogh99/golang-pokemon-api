package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"pokemon-api/database"
)

func getAllPokemons(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(database.PokemonDbAsValueArray())
}

func addPokemon(w http.ResponseWriter, r *http.Request) {
	var newPokemon database.Pokemon

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &newPokemon)

	if _, ok := database.PokemonDb[newPokemon.ID]; ok {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	database.PokemonDb[newPokemon.ID] = newPokemon
	w.WriteHeader(http.StatusOK)
}

func handleRequests() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "10005"
	}
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.Use(commonMiddleware)
	myRouter.HandleFunc("/pokemons", getAllPokemons).Methods("GET")
	myRouter.HandleFunc("/pokemons", addPokemon).Methods("POST")
	log.Fatal(http.ListenAndServe(":"+port, myRouter))
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func main() {
	fmt.Println("Pokemon Rest API")
	handleRequests()
}