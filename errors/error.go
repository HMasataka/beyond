package errors

import (
	"errors"

	"database/sql"

	"github.com/HMasataka/apperrors"
	"github.com/go-sql-driver/mysql"
)

const mysqlDuplicatedErrorCode = 1062

func IsDuplicated(err error) bool {
	var mySQLError *mysql.MySQLError
	return errors.As(err, &mySQLError) && mySQLError.Number == mysqlDuplicatedErrorCode
}

func IsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows) || apperrors.StatusCode(err) == NotFound.Code
}
