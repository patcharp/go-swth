package server

import (
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"strconv"
)

const (
	DefaultQueryPage = 1
	DefaultQuerySize = 20
	MaxQuerySize     = 100
)

type Pagination struct {
	Page uint
	Size uint
}

func GetEchoPagination(c echo.Context) Pagination {
	p := Pagination{
		Page: uint(atoi(c.QueryParam("page"), DefaultQueryPage)),
		Size: uint(atoi(c.QueryParam("size"), DefaultQuerySize)),
	}
	if p.Size > MaxQuerySize {
		p.Size = MaxQuerySize
	}
	return p
}

func GetGinPagination(c gin.Context) Pagination {
	p := Pagination{
		Page: uint(atoi(c.Query("page"), DefaultQueryPage)),
		Size: uint(atoi(c.Query("size"), DefaultQuerySize)),
	}
	if p.Size > MaxQuerySize {
		p.Size = MaxQuerySize
	}
	return p
}

func (p *Pagination) Offset() uint {
	return (p.Page - 1) * p.Size
}

func atoi(s string, v int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return v
	}
	return i
}
