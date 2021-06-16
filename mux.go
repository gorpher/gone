package gone

import "net/http"

var (
	mux     = http.NewServeMux()
	methods = map[string]map[string]http.HandlerFunc{}
)

func POST(pattern string, handlerFunc http.HandlerFunc) {
	handle(pattern, http.MethodPost, handlerFunc)
}

func HEAD(pattern string, handlerFunc http.HandlerFunc) {
	handle(pattern, http.MethodHead, handlerFunc)
}

func PUT(pattern string, handlerFunc http.HandlerFunc) {
	handle(pattern, http.MethodPut, handlerFunc)
}

func PATCH(pattern string, handlerFunc http.HandlerFunc) {
	handle(pattern, http.MethodPatch, handlerFunc)
}

func DELETE(pattern string, handlerFunc http.HandlerFunc) {
	handle(pattern, http.MethodDelete, handlerFunc)
}

func CONNECT(pattern string, handlerFunc http.HandlerFunc) {
	handle(pattern, http.MethodConnect, handlerFunc)
}

func OPTIONS(pattern string, handlerFunc http.HandlerFunc) {
	handle(pattern, http.MethodOptions, handlerFunc)
}

func TRACE(pattern string, handlerFunc http.HandlerFunc) {
	handle(pattern, http.MethodTrace, handlerFunc)
}

func GET(pattern string, handlerFunc http.HandlerFunc) {
	handle(pattern, http.MethodGet, handlerFunc)
}

func Handle(pattern string, handlerFunc http.HandlerFunc) {
	handle(pattern, "*", handlerFunc)
}

func handle(pattern string, method string, handlerFunc http.HandlerFunc) {
	var first = false
	if _, ok := methods[pattern]; !ok {
		methods[pattern] = map[string]http.HandlerFunc{
			method: handlerFunc,
		}
		first = true
	}
	if _, ok := methods[pattern][method]; !ok {
		methods[pattern][method] = handlerFunc
		first = true
	}
	if first {
		mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			h := methods[pattern][r.Method]
			if h == nil {
				h = methods[pattern]["*"]
			}
			if h != nil {
				h(w, r)
				return
			}
			http.NotFound(w, r)
		})
	}
}

func NewServeMux() *http.ServeMux {
	return mux
}
