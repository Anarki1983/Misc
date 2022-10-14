type middleware func(http.Handler) http.Handler

type Rounter struct {
	middlewareChain []middleware // use middleware
	mux map[string] http.Handler // map handler with route
}

func NewRouter() *Router{
    return &Router{}
}

func (r *Router) Use(m middleware) {
	r.middlewareChain = append(r.middlewareChain, m)
}

func (r *Router) Add(route string, h http.Handler) {
	mergedHandler := h

	// 從最後一個middleware開始呼叫
	for i:=len(r.middlewareChain) - 1; i >= 0; i-- {
		mergedHandler = r.middlewareChain[i](mergedHandler)
	}

	r.mux[route] = mergedHandler
}
