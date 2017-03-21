package labels

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// WeaveLabelsGetHandlerFunc turns a function with the right signature into a weave labels get handler
type WeaveLabelsGetHandlerFunc func(WeaveLabelsGetParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn WeaveLabelsGetHandlerFunc) Handle(params WeaveLabelsGetParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// WeaveLabelsGetHandler interface for that can handle valid weave labels get params
type WeaveLabelsGetHandler interface {
	Handle(WeaveLabelsGetParams, interface{}) middleware.Responder
}

// NewWeaveLabelsGet creates a new http.Handler for the weave labels get operation
func NewWeaveLabelsGet(ctx *middleware.Context, handler WeaveLabelsGetHandler) *WeaveLabelsGet {
	return &WeaveLabelsGet{Context: ctx, Handler: handler}
}

/*WeaveLabelsGet swagger:route GET /places/{placeId}/labels/{labelId} labels weaveLabelsGet

Get a label.

*/
type WeaveLabelsGet struct {
	Context *middleware.Context
	Handler WeaveLabelsGetHandler
}

func (o *WeaveLabelsGet) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)
	var Params = NewWeaveLabelsGetParams()

	uprinc, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	var principal interface{}
	if uprinc != nil {
		principal = uprinc
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
