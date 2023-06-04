package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Product struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
}

var (
	database = make(map[string]Product)
)

func setJSONResp(res http.ResponseWriter, message []byte, httpCode int) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(httpCode)
	res.Write(message)
}

func main() {
	//init db
	database["1"] = Product{ID: "1", Name: "Laptop", Price: 1000, Quantity: 10}
	database["2"] = Product{ID: "2", Name: "Mouse", Price: 10, Quantity: 100}
	database["3"] = Product{ID: "3", Name: "Keyboard", Price: 20, Quantity: 50}

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		message := []byte(`{"message": "Server Is Running"}`)
		setJSONResp(res, message, http.StatusOK)
	})

	http.HandleFunc("/get-products", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			message := []byte(`{"message": "Invalid method"}`)
			setJSONResp(res, message, http.StatusMethodNotAllowed)
			return
		}

		var products []Product

		for _, product := range database {
			products = append(products, product)
		}

		productJSON, err := json.Marshal(&products)
		if err != nil {
			message := []byte(`{"message": "Error while parsing JSON"}`)
			setJSONResp(res, message, http.StatusInternalServerError)
		}

		setJSONResp(res, productJSON, http.StatusOK)
	})

	http.HandleFunc("/add-product", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			message := []byte(`{"message": "Invalid method"}`)
			setJSONResp(res, message, http.StatusMethodNotAllowed)
			return
		}

		var product Product

		payload := req.Body

		defer payload.Close()

		err := json.NewDecoder(payload).Decode(&product)

		if err != nil {
			message := []byte(`{"message": "Error while parsing JSON"}`)
			setJSONResp(res, message, http.StatusInternalServerError)
			return
		}

		database[product.ID] = product
		message := []byte(`{"message": "Product added successfully"}`)

		setJSONResp(res, message, http.StatusOK)

	})

	http.HandleFunc("/get-product", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "GET" {
			message := []byte(`{"message":"Invalid Method"}`)
			setJSONResp(res, message, http.StatusMethodNotAllowed)
			return
		}

		if _, ok := req.URL.Query()["id"]; !ok {
			message := []byte(`{"message":"Invalid ID"}`)
			setJSONResp(res, message, http.StatusBadRequest)
			return
		}

		id := req.URL.Query()["id"][0]

		product, ok := database[id]

		if !ok {
			message := []byte(`{"message":"Product not found"}`)
			setJSONResp(res, message, http.StatusBadRequest)
			return
		}

		productJSON, err := json.Marshal(&product)
		if err != nil {
			message := []byte(`{"message":"Error while parsing JSON"}`)
			setJSONResp(res, message, http.StatusInternalServerError)
			return
		}

		setJSONResp(res, productJSON, http.StatusOK)

	})

	http.HandleFunc("/delete-product", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "DELETE" {
			message := []byte(`{"message":"Invalid Method"}`)
			setJSONResp(res, message, http.StatusMethodNotAllowed)
			return
		}

		if _, ok := req.URL.Query()["id"]; !ok {
			message := []byte(`{"message":"Invalid ID"}`)
			setJSONResp(res, message, http.StatusBadRequest)
			return
		}

		id := req.URL.Query()["id"][0]

		product, ok := database[id]

		if !ok {
			message := []byte(`{"message":"Product not found"}`)
			setJSONResp(res, message, http.StatusBadRequest)
			return
		}

		delete(database, id)

		productJSON, err := json.Marshal(&product)
		if err != nil {
			message := []byte(`{"message":"Error while parsing JSON"}`)
			setJSONResp(res, message, http.StatusInternalServerError)
			return
		}

		setJSONResp(res, productJSON, http.StatusOK)

	})

	http.HandleFunc("/update-product", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "PUT" {
			message := []byte(`{"message":"Invalid Method"}`)
			setJSONResp(res, message, http.StatusMethodNotAllowed)
			return
		}

		id := req.URL.Query()["id"][0]
		product, ok := database[id]
		if !ok {
			message := []byte(`{"message":"Invalid ID"}`)
			setJSONResp(res, message, http.StatusBadRequest)
			return
		}

		var newProduct Product

		payload := req.Body

		defer payload.Close()

		err := json.NewDecoder(payload).Decode(&newProduct)

		if err != nil {
			message := []byte(`{"message": "Error while parsing JSON"}`)
			setJSONResp(res, message, http.StatusInternalServerError)
			return
		}

		product.Name = newProduct.Name
		product.Price = newProduct.Price

		database[product.ID] = product

		productJSON, err := json.Marshal(&product)
		if err != nil {
			message := []byte(`{"message":"Error while parsing JSON"}`)
			setJSONResp(res, message, http.StatusInternalServerError)
			return
		}

		setJSONResp(res, productJSON, http.StatusOK)

	})

	err := http.ListenAndServe("localhost:9001", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
