package driver

import (
	"github.com/HMasataka/apperrors"
	errorcode "github.com/HMasataka/beyond/errors"
)

func wrapError(err error) error {
	if err == nil {
		return nil
	}

	if errorcode.IsNotFound(err) {
		return apperrors.Wrap(errorcode.NotFound, err)
	}

	if errorcode.IsDuplicated(err) {
		return apperrors.Wrap(errorcode.AlreadyUsedID, err)
	}

	return err
}
