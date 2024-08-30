package handler

import "io"

type HandlerInterface interface {
	Copy(io.Reader)

	Paste([]int) io.Reader
}
