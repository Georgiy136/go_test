package models

type Goods struct {
	Id          int    `json:"id"`
	ProjectID   int    `json:"project_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
	Removed     bool   `json:"removed"`
	CreatedAt   string `json:"created_at"`
}

type DataFromRequestGoodsAdd struct {
	ProjectID   int
	Name        string
	Description string
	Priority    int
}

type DataFromRequestGoodsUpdate struct {
	ID          int
	ProjectID   int
	Name        string
	Description *string
}

type DataFromRequestGoodsDelete struct {
	ID        int
	ProjectID int
}
type DataFromRequestGoodsList struct {
	ID        int
	ProjectID int
	Limit     int
	Offset    int
}

type GoodsListDBResponse struct {
	Meta struct {
		Total  int
		Remove int
		Limit  int
		Offset int
	} `json:"meta"`
	Goods []Goods `json:"goods"`
}

type DataFromRequestReprioritizeGood struct {
	ID          int
	NewPriority int
}
