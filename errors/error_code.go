package errors

import (
	"net/http"

	"github.com/HMasataka/apperrors"
)

var (
	BadRequest          = apperrors.ErrorCode{Code: http.StatusBadRequest, Name: "Bad Request"}
	InternalServerError = apperrors.ErrorCode{Code: http.StatusInternalServerError, Name: "Internal Server Error"}
	NotFound            = apperrors.ErrorCode{Code: http.StatusNotFound, Name: "Not Found"}
	Unauthorized        = apperrors.ErrorCode{Code: http.StatusUnauthorized, Name: "Unauthorized"}

	// カスタムCodeは100 < Code < 999のみ許容されている

	AlreadyUsedID = apperrors.ErrorCode{Code: 600, Name: "Already Used ID"}

	// ユーザー関連
	SuspensionUser = apperrors.ErrorCode{Code: 700, Name: "Suspension User"}
	BannedUser     = apperrors.ErrorCode{Code: 701, Name: "Banned User"}
)
