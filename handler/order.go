package handler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/code043/chi-api/model"
	"github.com/code043/chi-api/repository/order"
	"github.com/google/uuid"
)

type Order struct {
	Repo *order.RedisRepo
}

func (h *Order) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CustomerId uuid.UUID        `json:"customer_id"`
		LineItems  []model.LineItem `json:"line_items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	now := time.Now().UTC()

	order := model.Order{
		OrderId:    rand.Uint64(),
		CustomerId: body.CustomerId,
		LineItems:  body.LineItems,
		CreatedAt:  &now,
	}
	err := h.Repo.Insert(r.Context(), order)
	if err != nil {
		fmt.Println("failed to insert:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(order)
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(res)
	w.WriteHeader(http.StatusCreated)
}
func (o *Order) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List all orders")
}
func (o *Order) GetById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get an order by ID")
}
func (o *Order) UpdatedById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update an order by ID")
}
func (o *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete an order by ID")
}
