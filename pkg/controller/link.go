package controller

import (
	"errors"
	"net/http"

	"github.com/codeready-toolchain/registration-service/pkg/application"
	"github.com/codeready-toolchain/registration-service/pkg/context"
	crterrors "github.com/codeready-toolchain/registration-service/pkg/errors"
	"github.com/codeready-toolchain/registration-service/pkg/log"
	"github.com/gin-gonic/gin"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// Link implements the link endpoint, which provides direct links to different Sandbox components like WebConsole, DevSpaces, RHODS, etc.
type Link struct {
	app application.Application
}

// NewLink returns a new Link instance.
func NewLink(app application.Application) *Link {
	return &Link{
		app: app,
	}
}

// GetHandler returns a direct link based on the incoming request
func (l *Link) GetHandler(ctx *gin.Context) {
	userID := ctx.GetString(context.SubKey)
	username := ctx.GetString(context.UsernameKey)

	link, err := l.app.LinkService().Link(ctx.Request.URL, userID, username)
	e := &apierrors.StatusError{}
	if errors.As(err, &e) {
		crterrors.AbortWithError(ctx, int(e.Status().Code), err, "error generating link")
		return
	}
	if err != nil {
		log.Error(ctx, err, "error generating link")
		crterrors.AbortWithError(ctx, http.StatusInternalServerError, err, "error generating link")
		return
	}
	log.Infof(ctx, "Generated direct link. Input URL: %s. Output URL: %s", ctx.Request.URL.String(), link.OutputURL)

	ctx.JSON(http.StatusOK, link)
}
