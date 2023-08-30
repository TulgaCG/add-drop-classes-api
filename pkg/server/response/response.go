package response

import "errors"

var (
	ErrFailedToAuthenticate    = errors.New("failed to authenticate, check username, password or session")
	ErrFailedToFindAppCtxInCtx = errors.New("failed to find app context in the handler context")
	ErrFailedToFindRolesInCtx  = errors.New("failed to get roles from gin context")
	ErrInvalidParamIDFormat    = errors.New("invalid id format, id must be integer")
	ErrInsufficientPermission  = errors.New("user dont have permission")
	ErrInvalidRequestFormat    = errors.New("invalid request format")
	ErrContentNotFound         = errors.New("failed to find content in db")
)

type Response struct {
	Data  any    `json:"data"`
	Error string `json:"error"`
}

func New(data any, err error) Response {
	if err == nil {
		return Response{
			Data: data,
		}
	}

	return Response{
		Data:  data,
		Error: err.Error(),
	}
}

func WithData(data any) Response {
	return New(data, nil)
}

func WithError(err error) Response {
	return New(nil, err)
}
