package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"ioa/admin/pkg/errno"

	"github.com/sirupsen/logrus"
)

type Response struct {
	Code int         `json:"code"` // code为0时，表示成功
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"` // code != 0 时的错误信息
}

func ResponseJson(w http.ResponseWriter, data interface{}, err error) {
	code, msg := errno.DecodeErr(err)
	re := Response{
		Code: code,
		Data: data,
		Msg:  msg,
	}

	payload, err := json.Marshal(&re)
	if err != nil {
		logrus.Error(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(payload)
	if err != nil {
		logrus.Error(err)
	}
}

func BindWith(r *http.Request, request interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &request)
}
