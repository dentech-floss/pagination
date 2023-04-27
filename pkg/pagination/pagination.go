package pagination

import (
	"errors"
	"strconv"
)

type Page interface {
	Token() string
	Size() int
	Offset() int
	NextToken(resultSize int) *string
	SetTotalCount(totalCount uint64)
	GetTotalCount() uint64
}

type sqlPage struct {
	number     int
	size       int
	totalCount uint64
}

func NewSqlPage(
	token *string,
	size *int,
	defaultSize int,
	maxSize int,
) (Page, error) {
	_number := 1
	if token != nil {
		intToken, err := strconv.Atoi(*token)
		if err != nil {
			return nil, err
		}
		if intToken < 1 {
			return nil, errors.New("page number must be positive")
		}
		_number = intToken
	}

	_size := defaultSize
	if size != nil {
		if *size < 1 {
			return nil, errors.New("page size must be positive")
		}
		if int(*size) > maxSize {
			_size = maxSize
		} else {
			_size = int(*size)
		}
	}

	return &sqlPage{
		number: _number,
		size:   _size,
	}, nil
}

func (p *sqlPage) Token() string {
	return strconv.Itoa(p.number)
}

func (p *sqlPage) Size() int {
	return p.size
}

func (p *sqlPage) Offset() int {
	return (p.number - 1) * p.size
}

func (p *sqlPage) NextToken(resultSize int) *string {
	if resultSize < p.size {
		return nil
	} else if p.totalCount == 0 || p.totalCount > uint64(p.number)*uint64(p.size) {
		nextToken := strconv.Itoa(p.number + 1)
		return &nextToken
	}

	return nil
}

func (p *sqlPage) SetTotalCount(totalCount uint64) {
	p.totalCount = totalCount
}

func (p *sqlPage) GetTotalCount() uint64 {
	return p.totalCount
}
