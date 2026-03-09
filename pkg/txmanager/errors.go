package txmanager

import "fmt"

type TxError struct {
	Op  string
	Err error
}

func (e *TxError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("transaction manager: %s", e.Op)
	}
	return fmt.Sprintf("transaction manager: %s: %s", e.Op, e.Err.Error())
}

func (e *TxError) Unwrap() error {
	return e.Err
}

func WrapCommitError(err error) *TxError {
	return &TxError{
		Op:  "commit",
		Err: err,
	}
}

func WrapRollbackError(err error) *TxError {
	return &TxError{
		Op:  "rollback",
		Err: err,
	}
}

func WrapBeginError(err error) *TxError {
	return &TxError{
		Op:  "begin",
		Err: err,
	}
}

func WrapTransactionClosureError(err error) *TxError {
	return &TxError{
		Op:  "closure",
		Err: err,
	}
}

func WrapRetryExceeded(err error) *TxError {
	return &TxError{
		Op:  "retry exceeded",
		Err: err,
	}
}
