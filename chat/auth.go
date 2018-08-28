package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// ch2: 쿠키값을 확인하고 값이 없으면, 지정한 페이지로 리다이렉트한다.
	_, err := r.Cookie("auth") // 쿠키값을 가져오지 않고, 쿠키가 있는지 여부만 검사한다.
	if err == http.ErrNoCookie {
		// not authenticated
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		// some other error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// success - call the next handler
	h.next.ServeHTTP(w, r)
}

// 다음 핸들러를 위한 헬퍼 함수. 단순히 authoHanlder 를 wrapping 하는 역할
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// loginHandler handles the third-party login process.
// format: /auth/{action}/{provider}
// loginHanlder 는 http.Handler 를 구현하는 개체를 갖지 않는다. 여기서는 따로 상태(state)를 저장할 필요가 없기 때문이다.
// 따라서 main.go 에서 http.HandleFunc 을 통해 이 함수를 사용한다.
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3] // 이 코드는 나중에 panic 을 일으킬 수 있음. /auth/nonsense 처럼 segs[3]이 없는 경로로 접근하면..
	switch action {
	case "login":
		log.Println("TODO handle login for", provider)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Autho action %s not supported", action)
	}
}
