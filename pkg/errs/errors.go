package errs

import "strings"

type BoreError struct {
	Message       string
	OriginalError error
}

var (
	ErrInvalidConstructorArg = New("invalid arguments provided to New function")
	ErrStoragePathRequired   = New("storage path is required")
	ErrFailedToConnectToDB   = New("failed to connect to database")
	ErrFailedToCloseDB       = New("failed to close database connection")
	ErrFailedToRemoveDataDir = New("failed to remove data directory")
	ErrCollectionNotFound    = New("collection not found")
)

func New(message string) *BoreError {
	return &BoreError{
		Message:       message,
		OriginalError: nil,
	}
}

func Wrap(err error, message string) *BoreError {
	return &BoreError{
		Message:       message,
		OriginalError: err,
	}
}

func (e *BoreError) WithError(err error) *BoreError {
	return &BoreError{Message: e.Message, OriginalError: err}
}

func (e *BoreError) Error() string {
	if e.OriginalError == nil {
		return e.Message
	}

	return strings.TrimSpace(e.Message) + ": " + e.OriginalError.Error()
}

func (e *BoreError) String() string {
	return e.Error()
}

var _ error = (*BoreError)(nil)
