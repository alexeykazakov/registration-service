package service

import (
	"net/url"

	"github.com/codeready-toolchain/registration-service/pkg/application/service"
	"github.com/codeready-toolchain/registration-service/pkg/application/service/base"
	servicecontext "github.com/codeready-toolchain/registration-service/pkg/application/service/context"
	crterrors "github.com/codeready-toolchain/registration-service/pkg/errors"
	"github.com/codeready-toolchain/registration-service/pkg/link"
)

// ServiceImpl represents the implementation of the link service.
type ServiceImpl struct { // nolint:revive
	base.BaseService
}

type LinkServiceOption func(svc *ServiceImpl)

// NewLinkService creates a service object for generating direct links to different Sandbox components like WebConsole, DevSpaces, RHODS, etc
func NewLinkService(context servicecontext.ServiceContext, opts ...LinkServiceOption) service.LinkService {
	s := &ServiceImpl{
		BaseService: base.NewBaseService(context),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *ServiceImpl) Link(inputURL *url.URL, userID, username string) (*link.Link, error) {
	signup, err := s.Services().SignupService().GetSignup(userID, username)
	if err != nil {
		return nil, err
	}
	if signup == nil || !signup.Status.Ready {
		return nil, crterrors.NewForbiddenError("no active signup found", "")
	}

	var linkType link.LinkType
	for _, lt := range link.LinkTypes {
		if lt.Matches(inputURL.Path) {
			linkType = lt
			break
		}
	}

	if linkType == nil {
		linkType = link.DefaultLink
	}

	u, err := linkType.OutputURL(inputURL, *signup)
	if err != nil {
		return nil, err
	}
	return &link.Link{OutputURL: u}, nil
}
