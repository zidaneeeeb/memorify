package http

// import (
// 	"context"
// 	"encoding/json"
// 	contextlib "hbdtoyou/pkg/context"
// 	httplib "hbdtoyou/pkg/http"
// 	"log"
// 	"net/http"
// )

// func (h *paymentHandler) handleDeletePaymentByID(w http.ResponseWriter, r *http.Request, paymentID string) {
// 	// add timeout to context
// 	timeout := h.scopeSettings[ScopeDeletePayment].Timeout
// 	ctx, cancel := context.WithTimeout(r.Context(), timeout)
// 	defer cancel()

// 	var (
// 		err        error           // stores error in this handler
// 		source     string          // stores request source
// 		resBody    []byte          // stores response body to write
// 		statusCode = http.StatusOK // stores response status code
// 	)

// 	// write response
// 	defer func() {
// 		// error
// 		if err != nil {
// 			log.Printf("[Payment HTTP][handleDeletePayment] Failed to delete payment by ID. payment ID: %d, Source: %s, Err: %s\n", paymentID, source, err.Error())
// 			httplib.WriteErrorResponse(w, statusCode, []string{err.Error()})
// 			return
// 		}
// 		// success
// 		httplib.WriteResponse(w, resBody, statusCode, httplib.JSONContentTypeDecorator)
// 	}()

// 	// prepare channels for main go routine
// 	resChan := make(chan string, 1)
// 	errChan := make(chan error, 1)

// 	go func() {
// 		// get request source
// 		source, err = httplib.GetSourceFromHeader(r)
// 		if err != nil {
// 			statusCode = http.StatusBadRequest
// 			errChan <- errSourceNotProvided
// 			return
// 		}
// 		ctx = contextlib.SetSource(ctx, source)

// 		// get user ID
// 		reqUserID, err := httplib.GetUserIDFromHeader(r)
// 		if err != nil {
// 			statusCode = http.StatusBadRequest
// 			errChan <- errInvalidUserID
// 			return
// 		}
// 		ctx = contextlib.SetUserID(ctx, reqUserID)

// 		// get token from header
// 		token, err := httplib.GetBearerTokenFromHeader(r)
// 		if err != nil {
// 			statusCode = http.StatusBadRequest
// 			errChan <- errInvalidToken
// 			return
// 		}

// 		// check access token
// 		err = checkAccessToken(ctx, h.auth, token, reqUserID, "handleDeletePayment")
// 		if err != nil {
// 			statusCode = http.StatusUnauthorized
// 			errChan <- err
// 			return
// 		}

// 		// TODO: add authorization flow with roles

// 		var paymentID string
// 		err = h.payment.DeletePaymentByID(ctx, paymentID)
// 		if err != nil {
// 			// determine error and status code, by default its internal error
// 			parsedErr := errInternalServer
// 			statusCode = http.StatusInternalServerError
// 			if v, ok := mapHTTPError[err]; ok {
// 				parsedErr = v
// 				statusCode = http.StatusBadRequest
// 			}

// 			// log the actual error if its internal error
// 			if statusCode == http.StatusInternalServerError {
// 				log.Printf("[Payment HTTP][handleDeletePayment] Internal error from DeletePayment. Err: %s\n", err.Error())
// 			}

// 			errChan <- parsedErr
// 			return
// 		}

// 		resChan <- paymentID
// 	}()

// 	// wait and handle main go routine
// 	select {
// 	case <-ctx.Done():
// 		statusCode = http.StatusGatewayTimeout
// 		err = errRequestTimeout
// 	case err = <-errChan:
// 	case paymentID := <-resChan:
// 		resBody, err = json.Marshal(httplib.ResponseEnvelope{
// 			Data: paymentID,
// 		})
// 	}
// }
