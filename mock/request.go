package mock

type (
	RequestSvc struct {
		PostFunc func(path string, body any) (httpStatus int, err error)
	}
)

func (svc *RequestSvc) Post(path string, body any) (httpStatus int, err error) {
	return svc.PostFunc(path, body)
}
