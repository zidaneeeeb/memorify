package http

import (
	"context"
	"encoding/json"
	"hbdtoyou/internal/content"
	contextlib "hbdtoyou/pkg/context"
	httplib "hbdtoyou/pkg/http"
	"log"
	"net/http"
)

func (h *contentsHandler) handleGetContents(w http.ResponseWriter, r *http.Request) {
	// add timeout to context
	timeout := h.scopeSettings[ScopeGetContents].Timeout
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
		// error
		if err != nil {
			log.Printf("[Content HTTP][handleGetContents] Failed to get contents. Source: %s, Err: %s\n", source, err.Error())
			httplib.WriteErrorResponse(w, statusCode, []string{err.Error()})
			return
		}
		// success
		httplib.WriteResponse(w, resBody, statusCode, httplib.JSONContentTypeDecorator)
	}()

	// prepare channels for main go routine
	resChan := make(chan []content.Content, 1)
	errChan := make(chan error, 1)

	go func() {
		// get request source
		source, err = httplib.GetSourceFromHeader(r)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- errSourceNotProvided
			return
		}
		ctx = contextlib.SetSource(ctx, source)

		// get user ID
		reqUserID, err := httplib.GetUserIDFromHeader(r)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- errInvalidUserID
			return
		}
		ctx = contextlib.SetUserID(ctx, reqUserID)

		// get token from header
		token, err := httplib.GetBearerTokenFromHeader(r)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- errInvalidToken
			return
		}

		// check access token
		err = checkAccessToken(ctx, h.auth, token, reqUserID, "handleGetContents")
		if err != nil {
			statusCode = http.StatusUnauthorized
			errChan <- err
			return
		}

		var filter content.GetContentsFilter
		filter, err = h.parseHandleGetContentsQuery(r)
		if err != nil {
			statusCode = http.StatusBadRequest
			errChan <- errBadRequest
			return
		}

		var contents []content.Content
		contents, err = h.content.GetContents(ctx, filter)
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
				log.Printf("[Finance HTTP][handleGetContents] Internal error from GetContents. Err: %s\n", err.Error())
			}

			errChan <- parsedErr
			return
		}

		resChan <- contents
	}()

	// wait and handle main go routine
	select {
	case <-ctx.Done():
		statusCode = http.StatusGatewayTimeout
		err = errRequestTimeout
	case err = <-errChan:
	case res := <-resChan:
		// format each content
		contents := make([]contentHTTP, 0)
		for _, r := range res {
			contents = append(contents, formatContent(r))
		}

		resBody, err = json.Marshal(httplib.ResponseEnvelope{
			Data: contents,
		})
	}
}