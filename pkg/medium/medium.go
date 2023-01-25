package medium

import "github.com/go-resty/resty/v2"

type Medium struct {
	Matrix MatrixInterface
}

func NewMedium() *Medium {
	restyClient := resty.New()

	return &Medium{
		Matrix: NewMatrixClient(restyClient),
	}
}
