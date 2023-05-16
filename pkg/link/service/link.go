package service

import (
	"net/url"

	"github.com/codeready-toolchain/registration-service/pkg/link"
	"github.com/codeready-toolchain/registration-service/pkg/signup"

	"github.com/gin-gonic/gin"
)

func GenerateLink(inputURL string, ctx *gin.Context) (string, error) {
	// TODO Turn it into a Service
	// TODO Get Signup

	s := signup.Signup{}

	u, err := url.Parse(inputURL)
	if err != nil {
		return "", err
	}

	for _, lt := range link.LinkTypes {
		if lt.Matches(u.Path) {
			return lt.OutputURL(ctx, s)
		}
	}

	return link.DefaultLink.OutputURL(ctx, s)
}
