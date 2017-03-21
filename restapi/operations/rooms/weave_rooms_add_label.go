package rooms

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// WeaveRoomsAddLabelHandlerFunc turns a function with the right signature into a weave rooms add label handler
type WeaveRoomsAddLabelHandlerFunc func(WeaveRoomsAddLabelParams) middleware.Responder

// Handle executing the request and returning a response
func (fn WeaveRoomsAddLabelHandlerFunc) Handle(params WeaveRoomsAddLabelParams) middleware.Responder {
	return fn(params)
}

// WeaveRoomsAddLabelHandler interface for that can handle valid weave rooms add label params
type WeaveRoomsAddLabelHandler interface {
	Handle(WeaveRoomsAddLabelParams) middleware.Responder
}

// NewWeaveRoomsAddLabel creates a new http.Handler for the weave rooms add label operation
func NewWeaveRoomsAddLabel(ctx *middleware.Context, handler WeaveRoomsAddLabelHandler) *WeaveRoomsAddLabel {
	return &WeaveRoomsAddLabel{Context: ctx, Handler: handler}
}

/*WeaveRoomsAddLabel swagger:route POST /places/{placeId}/rooms/{roomId}/addLabel rooms weaveRoomsAddLabel

Adds a label to the room.

*/
type WeaveRoomsAddLabel struct {
	Context *middleware.Context
	Handler WeaveRoomsAddLabelHandler
}

func (o *WeaveRoomsAddLabel) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)
	var Params = NewWeaveRoomsAddLabelParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
