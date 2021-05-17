package repository

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/medivh13/koalatest/internal/models"
	"github.com/medivh13/koalatest/internal/repository"
	"github.com/medivh13/koalatest/pkg/common/crypto"
	btbErrors "github.com/medivh13/koalatest/pkg/errors"
)

const (
	Register        = `INSERT INTO customers (customer_id, customer_name, email, phone_number, dob, sex, salt, password, create_date) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	CekEmail        = `SELECT customer_id from customers where email = $1`
	CekPhone        = `SELECT customer_id from customers where phone_number = $1`
	CreateToken     = `INSERT INTO tokens (token_id, token, refresh_type, customer_id) VALUES ($1, $2, $3,$4,$5)`
	GetToken        = `SELECT token from tokens where customer_id = $1`
	CekEmailOrPhone = `SELECT customer_id, salt from customers where email = $1 OR phone_number = $2`
	CekPass         = `SELECT customer_id from customers where password = $1`
)

var statement PreparedStatement

type PreparedStatement struct {
	register        *sqlx.Stmt
	cekEmail        *sqlx.Stmt
	cekPhone        *sqlx.Stmt
	cekEmailOrPhone *sqlx.Stmt
	cekPass         *sqlx.Stmt
	getToken        *sqlx.Stmt
}

type MySQLRepo struct {
	Conn *sqlx.DB
}

func NewMySQLRepo(Conn *sqlx.DB) repository.Repository {

	repo := &MySQLRepo{Conn}
	InitPreparedStatement(repo)
	return repo
}

//Preparex query for database queries
func (m *MySQLRepo) Preparex(query string) *sqlx.Stmt {
	statement, err := m.Conn.Preparex(query)
	if err != nil {
		log.Fatalf("Failed to preparex query: %s. Error: %s", query, err.Error())
	}

	return statement
}

func InitPreparedStatement(m *MySQLRepo) {
	statement = PreparedStatement{
		cekEmail:        m.Preparex(CekEmail),
		cekPhone:        m.Preparex(CekPhone),
		cekEmailOrPhone: m.Preparex(CekEmailOrPhone),
		cekPass:         m.Preparex(CekPass),
		getToken:        m.Preparex(GetToken),
	}
}

func (m *MySQLRepo) Register(data *models.Customers) error {
	custID := []models.ExistingCustomer{}

	err := statement.cekEmail.Select(&custID, data.Email)
	if err != nil {
		log.Println("Failed Query to get Existing Email : ", err.Error())
		return fmt.Errorf(btbErrors.ErrorDB)
	}

	if len(custID) > 0 {
		return fmt.Errorf(btbErrors.ErrorExistingData)
	}

	err = statement.cekPhone.Select(&custID, data.PhoneNumber)
	if err != nil {
		log.Println("Failed Query to get Existing Phone : %s", err.Error())
		return fmt.Errorf(btbErrors.ErrorDB)
	}

	if len(custID) > 0 {
		return fmt.Errorf(btbErrors.ErrorExistingData)
	}
	tx, err := m.Conn.Beginx()
	if err != nil {
		log.Println("Failed to start database transaction. Error: %s", err.Error())
		return fmt.Errorf(btbErrors.ErrorDB)
	}

	_, err = tx.Exec(Register, data.CustomerId, data.CustomerName, data.Email, data.PhoneNumber, data.Dob, data.Sex, data.Salt, data.Password, data.CreatedDate)
	if err != nil {
		log.Println(err.Error())
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(CreateToken, data.CustomerId, crypto.EncodeSHA256HMAC(data.Salt, data.CustomerId), data.Salt, data.CustomerId)
	if err != nil {
		log.Println(err.Error())
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Failed to commit database transaction. Error: %s", err.Error())
	}

	return nil
}

func (m *MySQLRepo) GetToken(phoneOrEmail, pass string) ([]*models.GetTokens, error) {
	cust := []models.ExistingCustomerWithSalt{}
	custID := []models.ExistingCustomer{}

	var data []*models.GetTokens

	err := statement.cekEmailOrPhone.Select(&cust, phoneOrEmail)
	if err != nil {
		log.Println("Failed Query to get Existing Email or Phone : ", err.Error())
		return nil, fmt.Errorf(btbErrors.ErrorDB)
	}

	if len(cust) < 1 {
		return nil, fmt.Errorf(btbErrors.ErrorDataNotFound)
	}

	err = statement.cekPass.Select(&custID, crypto.EncodeSHA256HMAC(cust[0].Salt, pass))
	if err != nil {
		log.Println("Failed Query check Password : ", err.Error())
		return nil, fmt.Errorf(btbErrors.ErrorDB)
	}

	if len(custID) < 1 {
		return nil, fmt.Errorf(btbErrors.ErrorDataNotFound)
	}

	err = statement.getToken.Select(&data, custID[0].CustomerId)

	if len(data) == 0 {
		return nil, fmt.Errorf(btbErrors.ErrorDataNotFound)
	}

	return data, nil
}
