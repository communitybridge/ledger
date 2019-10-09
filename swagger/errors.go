package swagger

import (
	"github.com/communitybridge/ledger/gen/restapi/operations/health"

	"github.com/sirupsen/logrus"

	"github.com/communitybridge/ledger/gen/models"
	"github.com/go-openapi/runtime/middleware"
)

type codedResponse interface {
	Code() string
}

// ErrorResponse wraps the error in the api standard models.ErrorResponse object
func ErrorResponse(err error) *models.ErrorResponse {
	cd := ""
	if e, ok := err.(codedResponse); ok {
		cd = e.Code()
	}

	e := models.ErrorResponse{
		Code:    cd,
		Message: err.Error(),
	}
	return &e
}

// HealthErrorHandler handles error resp from calls to the health endpoint
func HealthErrorHandler(label string, err error) middleware.Responder {
	logrus.WithError(err).Error(label)

	return health.NewGetHealthBadRequest()

}
