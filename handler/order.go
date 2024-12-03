package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
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
func (h *Order) List(w http.ResponseWriter, r *http.Request) {
	cursorStr := r.URL.Query().Get("cursor")
	if cursorStr == "" {
		cursorStr = "0"
	}
	const decimal = 10
	const bitSize = 64
	cursor, err := strconv.ParseUint(cursorStr, decimal, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	const size = 50
	res, err := h.Repo.FindAll(r.Context(), order.FindAllPage{
		Offset: uint(cursor),
		Size:   size,
	})
	if err != nil {
		fmt.Println("failed to find all:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var response struct {
		Items []model.Order `json:"items"`
		Next  uint64        `json:"next,omitempty`
	}
	response.Items = res.Orders
	response.Next = res.Cursor

	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println("failed to marhsal: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
func (h *Order) GetById(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	//idParam := chi.URLParam(r, "id")
	idParam := strings.Split(path, "/")[2]

	const base = 10
	const bitSize = 64

	orderID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	o, err := h.Repo.FindById(r.Context(), orderID)
	if errors.Is(err, order.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find by id: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(o); err != nil {
		fmt.Println("failed to marshal: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func (o *Order) UpdatedById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update an order by ID")
}
func (o *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete an order by ID")
}
