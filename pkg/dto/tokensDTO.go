package dto

type TokensReqDTO struct {
	TokenId     string `json:"token_id"`
	Token       string `json:"token"`
	RefreshType string `json:"refresh_type"`
	CustomerId  string `json:"customer_id"`
}
