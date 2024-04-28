package customers

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Risuii/invoice/src/entity"
)

func (c *CustomersRepository) Create(ctx context.Context, data *entity.Customer) error {

	namedStmt, err := c.getNamedStatement(ctx, InsertCustomer)
	if err != nil {
		log.Println("getNamedStatement err: ", err)
		return err
	}

	_, err = namedStmt.ExecContext(ctx, data)
	if err != nil {
		log.Println("create customer err: ", err)
		return err
	}

	redisErr := c.redis.DelWithPattern(ctx, DeleteCustomerRedisKey)
	if redisErr != nil {
		log.Println(redisErr)
	}

	return nil
}

func (c *CustomersRepository) Get(ctx context.Context, id string) (entity.Customer, error) {
	var Customer entity.Customer

	err := c.redis.WithCache(ctx, fmt.Sprintf(GetDetailCustomersRedisKey, id), &Customer, func() (interface{}, error) {
		var customerData entity.Customer
		err := c.masterStmts[GetByID].GetContext(ctx, &customerData, id)
		return customerData, err
	})

	if err != nil {
		log.Println(err)
		return Customer, err
	}

	return Customer, nil
}

func (c *CustomersRepository) Update(ctx context.Context, data *entity.Customer) error {
	var rowsAffected int64

	namedStmt, err := c.getNamedStatement(ctx, UpdateCustomer)
	if err != nil {
		log.Println("get named statement err: ", err)
		return err
	}

	res, err := namedStmt.ExecContext(ctx, data)
	if err != nil {
		log.Println("exec err: ", err)
		return err
	}

	rowsAffected, err = res.RowsAffected()
	if err != nil {
		log.Println("Get rows affected err: ", err)
		return err
	}

	if rowsAffected == 0 {
		log.Println("ID not exist err: ", sql.ErrNoRows)
		return sql.ErrNoRows
	}

	redisErr := c.redis.DelWithPattern(ctx, DeleteCustomerRedisKey)
	if redisErr != nil {
		log.Println(redisErr)
	}

	return nil
}
