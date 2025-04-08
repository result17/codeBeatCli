package exitcode

import "strconv"

const (
	// Success is used when a heartbeat was sent successfully
	Success = 0
	// ErrGeneric is used for general erros
	ErrGeneric = 1
	// ErrAPI is when API returned an error
	ErrAPI = 102
)

type Err struct {
	Code int
}

func (e Err) Error() string {
	return strconv.Itoa(e.Code)
}
