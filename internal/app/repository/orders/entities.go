package orders

import (
	"github.com/jackc/pgtype"
	"github.com/shopspring/decimal"
)

type orderEntity struct {
	Id         uint64              `db:"id"`
	Number     string              `db:"order_number"`
	UserID     string              `db:"user_id"`
	StatusId   int                 `db:"status_id"`
	Accrual    decimal.NullDecimal `db:"accrual"`
	UploadedAt pgtype.Date         `db:"uploaded_at"`
	UpdatedAt  pgtype.Date         `db:"updated_at"`
}
