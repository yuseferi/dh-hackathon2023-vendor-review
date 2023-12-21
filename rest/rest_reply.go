package rest

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

const HeaderContentType = "Content-Type"

type JSONResponseError struct {
	Message string `json:"message"`
}

func Error(e error) *JSONResponseError {
	return &JSONResponseError{
		Message: e.Error(),
	}
}

func Reply(w http.ResponseWriter, status int, response any, logger *zap.Logger) {
	setContentTypeJSONAndUTF8(w)
	if response == nil {
		sendReply(w, status, nil, logger)
		return
	}

	marshalledResponse, err := json.Marshal(response)
	if err != nil {
		replyFailed(err, logger)
		responseErr := Error(err)
		marshalledResponseErr, err := json.Marshal(responseErr)
		if err != nil {
			replyFailed(err, logger)
		}
		sendReply(w, http.StatusInternalServerError, marshalledResponseErr, logger)
		return
	}

	sendReply(w, status, marshalledResponse, logger)
}

func replyFailed(err error, logger *zap.Logger) {
	logger.Warn("reply failed", zap.Error(err))
}
func sendReply(w http.ResponseWriter, status int, response []byte, logger *zap.Logger) {
	w.WriteHeader(status)
	if response == nil {
		return
	}

	_, err := w.Write(response)
	if err != nil {
		replyFailed(err, logger)
	}
}

func setContentTypeJSONAndUTF8(w http.ResponseWriter) {
	w.Header().Set(HeaderContentType, "application/json; charset=utf-8")
}
