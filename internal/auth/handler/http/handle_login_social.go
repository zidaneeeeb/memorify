package http

import (
	"context"
	"encoding/json"
	contextlib "hbdtoyou/pkg/context"
	httplib "hbdtoyou/pkg/http"
	"io/ioutil"
	"log"
	"net/http"
)

func (h *authHandler) handleLoginSocial(w http.ResponseWriter, r *http.Request) {
	// add timeout to context
	timeout := h.scopeSettings[ScopeLoginSocial].Timeout
	ctx, cancel := context.WithTimeout(r.Context(), timeout)
	defer cancel()

	var (
		err        error           // stores error in this handler
		source     string          // stores request source
		resBody    []byte          // stores response body to write
		statusCode = http.StatusOK // stores response status code
	)

	// write response
	defer func() {
		// add status code to context for monitoring
		ctx = contextlib.SetHTTPStatusCode(ctx, statusCode)
		*r = *(r.WithContext(ctx))

		// error
		if err != nil {
			log.Printf("[Auth HTTP][handleLoginSocial] Failed to login. Source: %s, Err: %s\n", source, err.Error())
			httplib.WriteErrorResponse(w, statusCode, []string{err.Error()})
			return
		}
		// success
		httplib.WriteResponse(w, resBody, statusCode, httplib.JSONContentTypeDecorator)
	}()

	// prepare channels for main go routine
	resChan := make(chan loginResponseData, 1)
	errChan := make(chan error, 1)

	go func() {
		// get request source
		source, err = httplib.GetSourceFromHeader(r)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- errSourceNotProvided
			return
		}

		// read request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- errBadRequest
			return
		}

		// unmarshall body
		var data loginSocialRequestData
		err = json.Unmarshal(body, &data)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- errBadRequest
			return
		}

		// login
		token, tokenData, err := h.auth.LoginSocial(ctx, data.TokenEmail)
		if err != nil {
			// determine error and status code, by default its internal error
			parsedErr := errInternalServer
			statusCode = http.StatusInternalServerError
			if v, ok := mapHTTPError[err]; ok {
				parsedErr = v
				statusCode = http.StatusBadRequest
			}

			// log the actual error if its internal error
			if statusCode == http.StatusInternalServerError {
				log.Printf("[User HTTP][handleLogin] Internal error from LoginBasic. Err: %s\n", err.Error())
			}

			errChan <- parsedErr
			return
		}

		resChan <- loginResponseData{
			UserID:   tokenData.UserID,
			Fullname: tokenData.Fullname,
			Email:    tokenData.Email,
			Token:    token,
		}
	}()

	// wait and handle main go routine
	select {
	case <-ctx.Done():
		statusCode = http.StatusGatewayTimeout
		err = errRequestTimeout
	case err = <-errChan:
	case resData := <-resChan:
		res := httplib.ResponseEnvelope{
			Data: resData,
		}
		resBody, err = json.Marshal(res)
	}
}
