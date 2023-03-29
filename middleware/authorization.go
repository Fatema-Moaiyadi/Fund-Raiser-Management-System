package middleware

import (
	"context"
	"fmt"
	"github.com/fatema-moaiyadi/fund-raiser-system/service"
	systemerrors "github.com/fatema-moaiyadi/fund-raiser-system/system_errors"
	"net/http"
	"strings"
)

//AuthorizeAdmin to check xauth for admin apis
func AuthorizeAdmin(tokenSvc service.TokenService, nextReq http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		token, err := readToken(req)
		if err != nil {
			systemerrors.WriteErrorResponse(res, err)
			return
		}

		payload, err := tokenSvc.VerifyToken(token)
		if err != nil {
			systemerrors.WriteErrorResponse(res, systemerrors.ErrInvalidToken)
			return
		}

		if !payload.IsAdmin {
			systemerrors.WriteErrorResponse(res, systemerrors.ErrForbidden)
			return
		}

		context := context.WithValue(req.Context(), "claims", payload)
		nextReq.ServeHTTP(res, req.WithContext(context))
	}
}

//Authorize to check xauth for its validity
func Authorize(tokenSvc service.TokenService, nextReq http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		token, err := readToken(req)
		if err != nil {
			systemerrors.WriteErrorResponse(res, err)
			return
		}

		payload, err := tokenSvc.VerifyToken(token)
		if err != nil {
			systemerrors.WriteErrorResponse(res, systemerrors.ErrInvalidToken)
			return
		}

		context := context.WithValue(req.Context(), "claims", payload)
		nextReq.ServeHTTP(res, req.WithContext(context))
	}
}

func readToken(req *http.Request) (string, error) {
	header := req.Header.Get("Authorization")

	headerString := fmt.Sprintf("%s", header)
	//if len of header in total is length,
	//we can safely presume that token is incorrect
	if len(headerString) < len("Bearer ") {
		return "", systemerrors.ErrMalformedToken
	}

	token := strings.TrimPrefix(headerString, "Bearer ")

	return token, nil
}
