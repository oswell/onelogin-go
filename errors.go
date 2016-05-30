package onelogin

import (
    "errors"
)

func ErrorOcurred(err error)(error) {
    logger.Errorf("An error occurred, %s", err.Error())
    return errors.New("An error ocurred.")
}
