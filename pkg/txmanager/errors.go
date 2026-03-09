package txmanager

import "fmt"

type TxError struct {
	op  string
	err error
}

func (e *TxError) Error() string {
	return fmt.Sprintf("transaction manager: %s: %s", e.op, e.err.Error())
}

func (e *TxError) Unwrap() error {
	return e.err
}

func WrapCommitError(err error) *TxError {
	return &TxError{
		op:  "commit",
		err: err,
	}
}

func WrapRollbackError(err error) *TxError {
	return &TxError{
		op:  "rollback",
		err: err,
	}
}

func WrapBeginError(err error) *TxError {
	return &TxError{
		op:  "begin",
		err: err,
	}
}

func WrapTransactionClosureError(err error) *TxError {
	return &TxError{
		op:  "closure",
		err: err,
	}
}

func WrapRetryExceeded(err error) *TxError {
	return &TxError{
		op:  "retry exceeded",
		err: err,
	}
}
