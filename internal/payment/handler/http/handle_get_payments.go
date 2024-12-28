package http

import (
	"context"
	"encoding/json"
	"hbdtoyou/internal/payment"
	contextlib "hbdtoyou/pkg/context"
	httplib "hbdtoyou/pkg/http"
	"log"
	"net/http"
)

func (h *paymentsHandler) handleGetPayments(w http.ResponseWriter, r *http.Request) {
	// add timeout to context
	timeout := h.scopeSettings[ScopeGetPayments].Timeout
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
			log.Printf("[Payment HTTP][handleGetPayments] Failed to get payments. Source: %s, Err: %s\n", source, err.Error())
			httplib.WriteErrorResponse(w, statusCode, []string{err.Error()})
			return
		}
		// success
		httplib.WriteResponse(w, resBody, statusCode, httplib.JSONContentTypeDecorator)
	}()

	// prepare channels for main go routine
	resChan := make(chan []payment.Payment, 1)
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
		err = checkAccessToken(ctx, h.auth, token, reqUserID, "handleGetPayments")
		if err != nil {
			statusCode = http.StatusUnauthorized
			errChan <- err
			return
		}

		var payments []payment.Payment
		payments, err = h.payment.GetPayments(ctx)
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
				log.Printf("[Payment HTTP][handleGetPayments] Internal error from GetPayments. Err: %s\n", err.Error())
			}

			errChan <- parsedErr
			return
		}

		resChan <- payments
	}()

	// wait and handle main go routine
	select {
	case <-ctx.Done():
		statusCode = http.StatusGatewayTimeout
		err = errRequestTimeout
	case err = <-errChan:
	case res := <-resChan:
		// format each payment
		payments := make([]paymentHTTP, 0)
		for _, r := range res {
			payments = append(payments, formatPayment(r))
		}

		resBody, err = json.Marshal(httplib.ResponseEnvelope{
			Data: payments,
		})
	}
}
