package httpx

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type JsonResponseBody struct {
	Data interface{} `json:"data"`
}

type ErrorResponseBody struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type Response struct {
	writer http.ResponseWriter
}

func (res Response) Json(body interface{}) {
	resBody := JsonResponseBody{
		Data: body,
	}

	data, err := json.Marshal(resBody)
	_, err = res.writer.Write(data)

	if err != nil {
		log.Error(err)
	}

}

func (res Response) ErrorAsJson(err error) {

	resBody := ErrorResponseBody{
		Message: err.Error(),
		Code:    400,
	}

	res.writer.WriteHeader(http.StatusBadRequest)

	data, err := json.Marshal(resBody)
	_, err = res.writer.Write(data)

	if err != nil {
		log.Error(err)
	}
}
