package lib

import (
	"bytes"
	"fmt"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
)

type Cusctx struct {
	C      echo.Context
	M      *mgo.Session
	B      *bytes.Buffer
	DBNAME string
}

func NewContext(c echo.Context, m *mgo.Session, b *bytes.Buffer) *Cusctx {
	var cc Cusctx
	cc.C = c
	cc.M = m
	cc.B = b
	return &cc
}

func (cc *Cusctx) Logf(format string, data ...interface{}) {
	if cc.B == nil {
		fmt.Printf("L::"+format+":::", data...)
		return
	}
	cc.B.WriteString(fmt.Sprintf("L:"+format, data...))
	cc.B.WriteString(":::")
}

func (cc *Cusctx) Errf(format string, data ...interface{}) {
	if cc.B == nil {
		fmt.Printf("E:"+format+":::", data...)
		return
	}
	cc.B.WriteString(fmt.Sprintf("E:"+format, data...))
	cc.B.WriteString(":::")
}
