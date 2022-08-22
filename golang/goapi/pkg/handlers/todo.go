package handlers

import (
	"encoding/json"
	"errors"

	"goapi/pkg/common/errs"
	"goapi/pkg/common/logger"
	"goapi/pkg/common/validator"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Todo struct {
	ID        int    `json:"id"`
	Text      string `json:"text" gorm:"column:title" validate:"required" `
	Completed bool   `json:"isCompleted" gorm:"column:is_done"`
}

type todoHandler struct {
	db *gorm.DB
}

func NewTodoHandler(db *gorm.DB) *todoHandler {
	return &todoHandler{
		db: db,
	}
}

func (h todoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	// step 1: แปลง JSON จาก request body เป็น Todo struct
	var todo Todo

	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		handleError(w, errs.NewBadRequestError(err.Error()))
		return
	}
	// step 2: validate
	err = validator.ValidateStruct(todo)
	if err != nil {
		handleError(w, errs.NewBadRequestError(err.Error()))
		return
	}
	// step 3: insert
	tx := h.db.Create(&todo)
	if err := tx.Error; err != nil {
		logger.Error(err.Error())
		handleError(w, errs.NewUnexpectedError(err.Error()))
		return
	}
	// step 4: response
	sendJson(w, http.StatusCreated, todo)
}

func (h todoHandler) ListTodo(w http.ResponseWriter, r *http.Request) {
	// เพิ่มอ่านค่าจาก query params
	query := r.URL.Query()
	wheres := map[string]interface{}{}
	// ถ้าส่ง completed มา ให้ใส่ไปใน map["is_done"] ตามชื่อ column จริงในฐานข้อมูล
	if val, ok := query["completed"]; ok {
		b1, err := strconv.ParseBool(val[0])
		if err != nil {
			handleError(w, errs.NewBadRequestError(err.Error()))
			return
		}
		wheres["is_done"] = b1
	}

	todos := []Todo{}
	// เพิ่ม .Where()
	tx := h.db.Where(wheres).Find(&todos)

	if err := tx.Error; err != nil {
		logger.Error(err.Error())
		handleError(w, errs.NewUnexpectedError(err.Error()))
		return
	}

	sendJson(w, http.StatusOK, todos)
}

func (h todoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	// step 1: get id from path param
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	// step 2: select where id
	todo := Todo{}
	tx := h.db.First(&todo, id)
	if err := tx.Error; err != nil {
		// step 3: handle error not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			handleError(w, errs.NewNotFoundError("todo with given id not found"))
			return
		}
		logger.Error(err.Error())
		handleError(w, errs.NewUnexpectedError(err.Error()))
		return
	}
	// step 4: response
	sendJson(w, http.StatusOK, todo)
}

func (h todoHandler) UpdateTodoStatus(w http.ResponseWriter, r *http.Request) {
	// step 1: get id from path param
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	// step 2: แปลง json body เป็น struct เพื่อเอาค่าสถานะที่ส่งมา
	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// step 3: update only is_done column
	tx := h.db.Model(Todo{ID: id}).Update("is_done", todo.Completed)
	if err := tx.Error; err != nil {
		logger.Error(err.Error())
		handleError(w, errs.NewUnexpectedError(err.Error()))
		return
	}
	// step 4: handle not found error
	if tx.RowsAffected == 0 {
		handleError(w, errs.NewNotFoundError("todo with given id not found"))
		return
	}
	// step 5: response
	w.WriteHeader(http.StatusNoContent)
}

func (h todoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	// step 1: get id from path param
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	// step 2: delete where id
	tx := h.db.Delete(&Todo{}, id)
	if err := tx.Error; err != nil {
		logger.Error(err.Error())
		handleError(w, errs.NewUnexpectedError(err.Error()))
		return
	}
	// step 3: handle not found error
	if tx.RowsAffected <= 0 {
		handleError(w, errs.NewNotFoundError("todo with given id not found"))
		return
	}
	// step 4: response
	w.WriteHeader(http.StatusNoContent)
}

func sendJson(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func handleError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case errs.AppError:
		sendJson(w, e.Code, e)
	case error:
		appErr := errs.AppError{
			Code:    http.StatusInternalServerError,
			Message: e.Error(),
		}
		sendJson(w, appErr.Code, appErr)
	}
}
