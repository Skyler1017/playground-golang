package foo

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var swaggerJson = `{
 "swagger": "2.0",
 "info": {
  "contact": {}
 },
 "paths": {
  "/trpc.playground.emptyTest.Hello/SayHello": {
   "post": {
    "parameters": [
     {
      "name": "HelloReq",
      "in": "body",
      "schema": {
       "$ref": "#/definitions/HelloReq"
      }
     }
    ],
    "responses": {
     "200": {
      "description": "",
      "schema": {
       "$ref": "#/definitions/HelloRsp"
      }
     }
    }
   }
  }
 },
 "definitions": {
  "HelloReq": {
   "type": "object",
   "format": "object",
   "properties": {
    "age": {
     "$ref": "#/definitions/google.protobuf.UInt64Value"
    },
    "name": {
     "type": "string"
    }
   }
  },
  "HelloRsp": {
   "type": "object",
   "format": "object",
   "properties": {
    "code": {
     "type": "number"
    },
    "msg": {
     "type": "string"
    }
   }
  },
  "google.protobuf.UInt64Value": {
   "type": "object",
   "format": "object",
   "properties": {
    "value": {
     "type": "number"
    }
   }
  }
 }
}`
var readFromFile = false
var fileName = "foo.json"

func writeToFile() error {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	} else {
		n, _ := f.Seek(0, os.SEEK_END)
		_, err = f.WriteAt([]byte(swaggerJson), n)
		defer f.Close()
	}
	return err
}

func readSwaggerJsoFromFile() string {
	if _, err := os.Stat(fileName); err == nil {
		f, err := os.Open(fileName)
		if err != nil {
			return swaggerJson
		}
		defer f.Close()
		swaggerJsonByte, err := ioutil.ReadAll(f)
		if err != nil {
			return swaggerJson
		}
		return string(swaggerJsonByte)
	}
	writeToFile()
	return swaggerJson
}

func SwaggerJsonHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("content-type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Write([]byte(GetSwaggerJson()))
}

func StartHttpService(uri string, port int, async bool) {
	mux := http.NewServeMux()
	mux.HandleFunc(uri, SwaggerJsonHandler)
	if port == 0 {
		port = 5000
	}
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
	if async {
		go func() {
			defer func() {
				if e := recover(); e != nil {
					fmt.Println(e)
				}
			}()
			server.ListenAndServe()
		}()
	} else {
		server.ListenAndServe()
	}
}

func GetSwaggerJson() string {
	if readFromFile {
		return readSwaggerJsoFromFile()
	}
	return swaggerJson
}
