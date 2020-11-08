package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/wesleywillians/go-rabbitmq/queue"
)

type Order struct {
	Coupon   string
	CcNumber string
}

type Result struct {
	Status string
}

// Primeira func a ser executada
func init() {
	// err irá receber o erro caso ocorra
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/process", process)
	http.ListenAndServe(":9090", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/home.html"))
	t.Execute(w, Result{})
}

func process(w http.ResponseWriter, r *http.Request) {

	coupon := r.PostFormValue("coupon")
	ccNumber := r.PostFormValue("cc-number")

	// Cria um objeto (struct)
	order := Order{
		Coupon:   coupon,
		CcNumber: ccNumber,
	}

	// Converte order para JSON
	jsonOrder, err := json.Marshal(order)
	if err != nil {
		log.Fatal("Error parsing to json")
	}

	rabbitMQ := queue.NewRabbitMQ()
	// Abre um canal e faz ligação com habbitMQ
	ch := rabbitMQ.Connect()
	// Fecha canal
	defer ch.Close()

	// Converte em formato de texto, contentType, exange,
	err = rabbitMQ.Notify(string(jsonOrder), "application/json", "orders_ex", "")
	if err != nil {
		// OBS- FATAL FAZ PARAR O SISTEMA
		log.Fatal("Error sending message to the queue")
	}

	t := template.Must(template.ParseFiles("templates/process.html"))
	t.Execute(w, "")
}
