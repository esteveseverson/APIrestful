package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Pessoa struct {
	ID   int    `json:"id"`
	Nome string `json:"nome"`
}

// docker run --name restApi -e MYSQL_ROOT_PASSWORD=test123 -d mysql:8
var db *sql.DB

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

	rows, err := db.Query("SELECT id, nome FROM pessoas")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pessoas []Pessoa
	for rows.Next() {
		var p Pessoa
		if err := rows.Scan(&p.ID, &p.Nome); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pessoas = append(pessoas, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pessoas)

}

func getPessoaByID(w http.ResponseWriter, r *http.Request) {

	idStr := regexp.MustCompile(`\d+`).FindString(r.URL.Path)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var p Pessoa
	err = db.QueryRow("SELECT id, nome FROM pessoas WHERE id = ?", id).Scan(&p.ID, &p.Nome)
	if err == sql.ErrNoRows {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(nil)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nil)
}

func createPessoa(w http.ResponseWriter, r *http.Request) {

	var p Pessoa
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO pessoas (nome) VALUES (?)", p.Nome)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	p.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func deletePessoaByID(w http.ResponseWriter, r *http.Request) {

	idStr := regexp.MustCompile(`\d+`).FindString(r.URL.Path)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM pessoas WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Pessoa deletada com sucesso"})
}

func main() {

	var err error
	db, err = sql.Open("mysql", "root:test123@tcp(localhost:3306)/api_db")
	if err != nil {
		fmt.Println("Erro ao conectar ao banco de dados:", err)
		panic(err)
	}
	defer db.Close()

	http.HandleFunc("/", handler)
	fmt.Println("Servidor rodando na porta 8080")
	http.ListenAndServe(":8080", nil)

}
