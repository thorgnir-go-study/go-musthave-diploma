package orders

//
//import (
//	"context"
//	"github.com/jackc/pgtype"
//	shopspring "github.com/jackc/pgtype/ext/shopspring-numeric"
//	_ "github.com/jackc/pgx/v4/stdlib"
//	"github.com/jmoiron/sqlx"
//)
//
//type ordersPostgresRepository struct {
//	DB *sqlx.DB
//}
//
//func NewOrdersPostgresRepository(_ context.Context, dsn string) (*ordersPostgresRepository, error) {
//	db, err := sqlx.Open("pgx", dsn)
//	//db.Conn.ConnInfo().RegisterDataType(pgtype.DataType{
//		Value: &shopspring.Numeric{},
//		Name:  "numeric",
//		OID:   pgtype.NumericOID,
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	//if err = prepareStatements(db); err != nil {
//	//	return nil, err
//	//}
//
//	return &ordersPostgresRepository{DB: db}, nil
//}
