package main

import (
	"fmt"
	"net/http"
	"regexp"
)

type Pessoa struct {
	ID   string
	Nome string
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET" && r.URL.Path == "/pessoa":
		fmt.Println("Entrou no metodo GET/pessoa")
		w.Write([]byte("Entrou no metodo GET/pessoa"))
	case r.Method == "GET" && regexp.MustCompile(`^/pessoa/\d+$`).MatchString(r.URL.Path):
		fmt.Println("Entrou no metodo GET/pessoa/{id}")
		w.Write([]byte("Entrou no metodo GET/pessoa/{id}"))
	case r.Method == "POST" && r.URL.Path == "/pessoa":
		fmt.Println("Entrou no metodo POST/pessoa")
		w.Write([]byte("Entrou no metodo POST/pessoa"))
	case r.Method == "DELETE" && regexp.MustCompile(`^/pessoa/\d+$`).MatchString(r.URL.Path):
		fmt.Println("Entrou no metodo DELETE/pessoa/{id}")
		w.Write([]byte("Entrou no metodo DELETE/pessoa/{id}"))
	default:
		fmt.Println("Metodo não cadastrado")
		w.Write([]byte("Metodo não cadastrado"))
	}
}

func main() {

	/*httpObj := NewHttpInterface()

	http.Handle("/foo", httpObj)*/

	http.HandleFunc("/", handler)
	fmt.Println("Servidor rodando na porta 8080")
	http.ListenAndServe(":8080", nil)

}

/*type HttpStructImplementation struct {
}

func NewHttpInterface() http.Handler {
	return &HttpStructImplementation{}
}

func (h *HttpStructImplementation) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("123"))
}*/
