package apidoc

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"
)

const (
	blockCommentStart  = "/**"
	blockCommentPrefix = " * "
	blockCommentEnd    = " */"
	lineBreak          = "\n"
	jsonIndentString   = "\t"
)

type (
	ApiDefine struct {
		Name           string
		Group          string
		Request        *http.Request
		Headers        []*ApiHeader
		HeaderExample  *Example
		Params         []*ApiParam
		ParamExample   *Example
		Success        []*ApiSuccess
		SuccessExample []*Example
		Errors         []*ApiError
		ErrorExample   []*Example
		//Errors      []*ApiError
	}
)

var (
	bStartLine = []byte(blockCommentStart + lineBreak)
	bPrefix    = []byte(blockCommentPrefix)
	bEndLine   = []byte(blockCommentEnd + lineBreak)
	bSpaceLine = []byte(blockCommentPrefix + lineBreak)
	bbreak     = []byte(lineBreak)
)

//func New2(req *http.Request, name string) ApiDefine {
//	items := strings.Split(req.URL.Path, "/")
//	var packageName string
//	for _, u := range items {
//		if u != "" {
//			packageName = u
//			break
//		}
//	}
//	return ApiDefine{
//		Request: req,
//		Group:   strings.Title(packageName),
//		Name:    name,
//	}
//}

func New(name string) ApiDefine {
	return ApiDefine{
		Name: name,
	}
}

func (ad *ApiDefine) AddParams(vs ...interface{}) *ApiDefine {
	for _, v := range vs {
		rv := reflect.ValueOf(v)
		if rv.Kind() != reflect.Slice {
			ad.addParams(v)
			continue
		}
		// slice は分解して実行
		for i := 0; i < rv.Len(); i++ {
			ad.addParams(rv.Index(i).Interface())
		}
	}
	return ad
}

func (ad *ApiDefine) addParams(v interface{}) *ApiDefine {
	// parameter の type を判定して、id,name,subname を返す
	switch t := v.(type) {
	case *ApiHeader:
		ad.Headers = append(ad.Headers, v.(*ApiHeader))
	case *ApiParam:
		ad.Params = append(ad.Params, v.(*ApiParam))
	case *ApiSuccess:
		ad.Success = append(ad.Success, v.(*ApiSuccess))
	case *ApiError:
		ad.Errors = append(ad.Errors, v.(*ApiError))
	default:
		panic(fmt.Sprintf("failed to add %v type invalid for apiDefine, %v", t))
	}
	return ad
}

func (ad *ApiDefine) AddExample(e *Example) *ApiDefine {
	// parameter の type を判定して、id,name,subname を返す
	switch e.typ {
	case exampleTypeHeader:
		ad.HeaderExample = e
	case exampleTypeParam:
		ad.ParamExample = e
	case exampleTypeSuccess:
		ad.SuccessExample = append(ad.SuccessExample, e)
	case exampleTypeError:
		ad.ErrorExample = append(ad.ErrorExample, e)
	default:
		panic(fmt.Sprintf("failed to add example %v type invalid for apiDefine"))
	}
	return ad
}

func (ad *ApiDefine) WriteBytes(b *bytes.Buffer) {

	// Define
	b.Write(bStartLine)
	b.Write(bPrefix)
	b.Write([]byte("@apiDefine " + ad.Name))
	b.Write(bbreak)

	b.Write(bSpaceLine)

	// Header
	if len(ad.Headers) > 0 {
		for _, header := range ad.Headers {
			WriteRowByte(b, header.Byte())
		}
		b.Write(bSpaceLine)
	}
	if ad.HeaderExample != nil {
		ad.HeaderExample.WriteIndentString(b)
		b.Write(bSpaceLine)
	}

	// Param
	if len(ad.Params) > 0 {
		for _, param := range ad.Params {
			WriteRowByte(b, param.Byte())
		}
		b.Write(bSpaceLine)
	}
	if ad.ParamExample != nil {
		ad.ParamExample.WriteIndentString(b)
		b.Write(bSpaceLine)
	}

	// Success
	if len(ad.Success) > 0 {
		for _, param := range ad.Success {
			WriteRowByte(b, param.Byte())
		}
		b.Write(bSpaceLine)
	}
	if len(ad.SuccessExample) > 0 {
		for _, successExample := range ad.SuccessExample {
			successExample.WriteIndentString(b)
		}
	}

	// Errors
	if len(ad.Errors) > 0 {
		for _, param := range ad.Errors {
			WriteRowByte(b, param.Byte())
		}
		b.Write(bSpaceLine)
	}
	if len(ad.ErrorExample) > 0 {
		for _, errorExample := range ad.ErrorExample {
			errorExample.WriteIndentString(b)
		}
	}

	b.Write(bEndLine)
	b.Write(bbreak)
	b.Write(bbreak)
}

func WriteRow(bb *bytes.Buffer, text string) {
	WriteRowByte(bb, []byte(text))
}

func WriteRowByte(bb *bytes.Buffer, b []byte) {
	bb.Write(bPrefix)
	bb.Write(b)
	bb.Write(bbreak)
}
