package kuipAuth

import (
	"ioa/context"

	"errors"
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

type filter struct {
	AuthKey   string
}

func New(arg string) (*filter, error) {
	filter := filter{}
	filter.AuthKey = arg

	return &filter, nil
}

func (f *filter) Name() string {
	return "kuipAuth"
}

func (f *filter) Request(ctx *context.Context) error {
	token, err := request.ParseFromRequest(ctx.Request, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		return []byte(f.AuthKey), nil
	})

	if err != nil {
		log.Println(err)
		ctx.Cancel()
		return err
	}

	if !token.Valid {
		ctx.Cancel()
		log.Println("token valid:", token.Valid)
		return errors.New("token invalid.")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("parse claims error")
		ctx.Cancel()
		return nil
	}

	userId := claims["user_id"].(string)
	ctx.Request.Header.Set("user_id", userId)
	return nil
}

func (f *filter) Response(ctx *context.Context) error {
	return nil
}
