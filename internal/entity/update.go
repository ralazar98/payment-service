package entity

type UpdateBalance struct {
	UserID            int    `json:"userID"`
	Operation         string `json:"operation"`
	ChangingInBalance int    `json:"changingInBalance"`
}
