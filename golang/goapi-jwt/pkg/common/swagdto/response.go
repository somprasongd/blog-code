package swagdto

type PagingResult struct {
	Page      int `json:"page" example:"1"`
	Limit     int `json:"limit" example:"10"`
	PrevPage  int `json:"prevPage" example:"0"`
	NextPage  int `json:"nextPage" example:"2"`
	Count     int `json:"count" example:"20"`
	TotalPage int `json:"totalPage" example:"2"`
}

type Response struct {
	Status    int         `json:"status" example:"200"`
	Data      interface{} `json:"data,omitempty"`
	RequestId string      `json:"requestId" example:"3b6272b9-1ef1-45e0"`
}

type ResponseWithPage struct {
	Status     int          `json:"status" example:"200"`
	Data       interface{}  `json:"data,omitempty"`
	Pagination PagingResult `json:"_pagination,omitempty"`
	RequestId  string       `json:"requestId" example:"3b6272b9-1ef1-45e0"`
}

type ResponseWithPageAndFilter struct {
	Status     int          `json:"status" example:"200"`
	Data       interface{}  `json:"data,omitempty"`
	Pagination PagingResult `json:"_pagination,omitempty"`
	Filters    interface{}  `json:"_filters,omitempty" example:"{status:'done'}"`
	RequestId  string       `json:"requestId" example:"3b6272b9-1ef1-45e0"`
}
