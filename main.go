package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

type Pessoa struct {
	ID   int    `json:"id"`
	Nome string `json:"nome"`
}

var pessoas []Pessoa
var nextID int = 1

func handler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET" && r.URL.Path == "/pessoa":
		getPessoas(w, r)
	case r.Method == "GET" && regexp.MustCompile(`^/pessoa/\d+$`).MatchString(r.URL.Path):
		getPessoaByID(w, r)
	case r.Method == "POST" && r.URL.Path == "/pessoa":
		createPessoa(w, r)
	case r.Method == "DELETE" && regexp.MustCompile(`^/pessoa/\d+$`).MatchString(r.URL.Path):
		deletePessoaByID(w, r)
	default:
		http.Error(w, "Metodo não cadastrado", http.StatusNotFound)
	}
}

func getPessoas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pessoas)
}

func getPessoaByID(w http.ResponseWriter, r *http.Request) {
	idStr := regexp.MustCompile(`\d+`).FindString(r.URL.Path)
	id, _ := strconv.Atoi(idStr)

	for _, p := range pessoas {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nil)
}

func createPessoa(w http.ResponseWriter, r *http.Request) {
	var p Pessoa
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	p.ID = nextID
	nextID++
	pessoas = append(pessoas, p)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func deletePessoaByID(w http.ResponseWriter, r *http.Request) {
	idStr := regexp.MustCompile(`\d+`).FindString(r.URL.Path)
	id, _ := strconv.Atoi(idStr)

	newPessoas := pessoas[:0]
	for _, p := range pessoas {
		if p.ID != id {
			newPessoas = append(newPessoas, p)
		}
	}
	pessoas = newPessoas

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pessoas)
}

func main() {

	http.HandleFunc("/", handler)
	fmt.Println("Servidor rodando na porta 8080")
	http.ListenAndServe(":8080", nil)

}
