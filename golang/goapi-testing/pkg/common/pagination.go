package common

import (
	"math"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	"gorm.io/gorm"
)

type PagingRequest struct {
	Page  int
	Limit int
	Order string
}

// Paginator - Populate pagingRequest object from request query
func Paginator(ctx HContext) PagingRequest {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	order := ctx.DefaultQuery("order", "")

	return PagingRequest{
		Page:  page,
		Limit: limit,
		Order: order,
	}
}

type PagingResult struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	PrevPage  int `json:"prevPage"`
	NextPage  int `json:"nextPage"`
	Count     int `json:"count"`
	TotalPage int `json:"totalPage"`
}

type Pagination struct {
	PagingRequest
	Query   *gorm.DB
	Records interface{}
}

func (p Pagination) Paginate() (*PagingResult, error) {
	return p.PaginateWithAllowedFields(map[string]interface{}{})
}

func (p *Pagination) PaginateWithAllowedFields(allowedFields map[string]interface{}) (*PagingResult, error) {
	page := p.Page
	limit := p.Limit
	order := parseOrder(p.Order, allowedFields)

	ch := make(chan int)
	go p.countRecords(ch)

	offset := (page - 1) * limit
	tx := p.Query.Begin()
	if len(order) > 0 {
		tx = tx.Order(order)
	}

	err := tx.Limit(limit).Offset(offset).Find(p.Records).Error
	if err != nil {
		return nil, err
	}

	tx.Commit()

	count := <-ch
	totalPage := int(math.Ceil(float64(count) / float64(limit)))

	var nextPage int
	if page == totalPage {
		nextPage = totalPage
	} else {
		nextPage = page + 1
	}

	return &PagingResult{
		Page:      page,
		Limit:     limit,
		Count:     count,
		PrevPage:  page - 1,
		NextPage:  nextPage,
		TotalPage: totalPage,
	}, nil
}

func (p *Pagination) countRecords(ch chan int) {
	var count int64
	tx := p.Query.Begin()
	tx.Model(p.Records).Count(&count)
	tx.Commit()
	ch <- int(count)
}

func parseOrder(order string, allowedFields map[string]interface{}) string {
	var orderBy string
	var orderFieldWithSort []string
	var orderFieldSnake string

	if len(allowedFields) == 0 {
		allowedFields = map[string]interface{}{
			"id":         "true",
			"name":       "true",
			"code":       "true",
			"created_at": "true",
			"updated_at": "true",
		}
	}

	orderString := ""

	if len(order) == 0 {
		return orderString
	}

	orders := strings.Split(order, ",")

	for _, orderField := range orders {
		orderField = strings.TrimSpace(orderField)
		orderFieldWithSort = strings.Split(orderField, " ")

		if len(orderFieldWithSort) >= 1 {
			if len(orderFieldWithSort) == 1 {
				orderBy = "DESC"
			} else {
				orderBy = orderFieldWithSort[1]
				orderBy = strings.TrimSpace(orderBy)
				orderBy = strings.ToUpper(orderBy)
			}

			orderField = strings.TrimSpace(orderFieldWithSort[0])
			orderFieldSnake = strcase.ToSnake(orderField)

			if orderBy != "ASC" && orderBy != "DESC" {
				orderBy = "DESC"
			}

			if value, ok := allowedFields[orderFieldSnake]; ok && value == "true" {
				if orderString == "" {
					orderString = orderFieldSnake + " " + orderBy
				} else {
					orderString = orderString + ", " + orderFieldSnake + " " + orderBy
				}
			}
		}
	}

	return orderString
}
