package runtime

import (
	"reflect"
	"testing"

	"github.com/gogim1/goscript/file"
	"github.com/stretchr/testify/assert"
)

func TestTypeCheck(t *testing.T) {
	assert.Nil(t, typeCheck(file.SourceLocation{}, []Value{}, []reflect.Type{}))
	assert.Nil(t, typeCheck(file.SourceLocation{}, nil, nil))
	assert.Nil(t, typeCheck(file.SourceLocation{}, []Value{NewVoid()}, []reflect.Type{VoidType}))
	assert.Nil(t, typeCheck(file.SourceLocation{}, []Value{NewString("str1"), NewString("str2")}, []reflect.Type{StringType, StringType}))

	assert.NotNil(t, typeCheck(file.SourceLocation{}, []Value{NewVoid()}, []reflect.Type{}))
	assert.NotNil(t, typeCheck(file.SourceLocation{}, []Value{NewVoid()}, nil))
	assert.NotNil(t, typeCheck(file.SourceLocation{}, []Value{}, []reflect.Type{VoidType}))
	assert.NotNil(t, typeCheck(file.SourceLocation{}, nil, []reflect.Type{VoidType}))
	assert.NotNil(t, typeCheck(file.SourceLocation{}, []Value{NewVoid()}, []reflect.Type{StringType}))

	assert.Nil(t, typeCheck(file.SourceLocation{}, []Value{NewVoid()}, []reflect.Type{ValueType}))

}
