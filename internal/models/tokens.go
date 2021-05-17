package models

type Tokens struct {
	TokenId     string `db:"token_id"`
	Token       string `db:"token"`
	RefreshType string `db:"refresh_type"`
	CustomerId  string `db:"customer_id"`
}

type GetTokens struct {
	Token string `db:"token"`
}
