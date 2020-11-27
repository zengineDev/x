package httpx

type ControllerContract interface {
	Index(res Response, req Request)
	Get(res Response, req Request)
	Store(res Response, req Request)
	Update(res Response, req Request)
	Delete(res Response, req Request)
}
