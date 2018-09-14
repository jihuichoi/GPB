package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

// set the active Avatar imlemetation
var avatars Avatar = UseFileSystemAvatar

// teml represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// net/http 의 HandleFunc 과 유사한 함수. Handler 인터페이스를 통해 사용
// type Handler interface {
// 	ServeHTTP(ResponseWriter, *Request)
// }
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	// ch2: Oauth2 를 통해 provider 로 부터 받아 쿠키에 저장한 사용자 정보를 불러온다.
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	// 템플릿에 request 정보를 전달
	// 자바스크립트 부분에서 {{.Host}} 라는 형식으로 사용할 수 있음
	// t.templ.Execute(w, r)
	// t.templ.Execute(w, nil)

	// ch2: request 을 직접 전달하지 않고, map 을 만들어서 전달
	// 템플릿에서 사용자 정보를 표시하기 위해
	t.templ.Execute(w, data)
}

func main() {

	// 채팅 사이트 주소가 하드코딩됨 (localhost:8080)
	// 이를 커맨드라인에서 -addr 라는 플래그로 처리하도록 경 ./chat -addr=":3000" 이라는 형식으로 실행이 가능해짐
	var addr = flag.String("host", ":8080", "The addr of the application")
	flag.Parse() // parse the flags

	// Oauth2
	// setup gomniauth
	gomniauth.SetSecurityKey("PUT YOUR AUTH KEY HERE")
	gomniauth.WithProviders(
		facebook.New("233530930663961", "c4dc9bf4d7dcc93c8d70f53610470a4d",
			"http://localhost:8080/auth/callback/facebook"),
		github.New("4530f6f362f4798105a3", "60433e0e932ccb55cc36df5f0b8962476d4f6473",
			"http://localhost:8080/auth/callback/github"),
		google.New("767497021571-q4l7entsul2qhjd4fmvt5it0cppe9n0m.apps.googleusercontent.com", "Ow8ry0I_rz6UHtHQIr3ZT6ol",
			"http://localhost:8080/auth/callback/google"),
	)

	// newRoom 함수로 새 룸을 만든다.
	// r := newRoom(UseFileSystemAvatar)
	// r := newRoom(UseGravatar)
	// r := newRoom(UseAuthAvatar)
	r := newRoom()
	// tracer 출력을 stdout 으로 내보냄
	// r.tracer = trace.New(os.Stdout)

	// net/http 기본 핸들러함수 사용
	// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 		w.Write([]byte(`
	// <html>
	// 	<head>
	// 		<title>chat</title>
	// 	</head>
	// 	<body>
	// 		Let's chat!
	// 	</body>
	// </html>
	// `))
	// 	})

	// http.HandleFunc 를 http.Handle 로 교체
	// http.Handle(pattern string, handler Handler)
	// templateHandler struct 에는 ServeHTTP 라는 메서드가 있음. 이는 http.Handle의 인자값인 Handler 인터페이스 조건에 만족
	// templateHandler 의 인스턴스를 생성하고 그 포인터를 전달
	// QST: 그러면 templateHandler 의 ServeHTTP 메서드가 동작하나? 왜?
	// http.Handle 은 Handler 인터페이스를 사용하기 위한 함수. 즉, Handle 이 Handler 인터페이스를 만족하는 struct 를 인수값으로 받아서........
	// 뭐래..
	// http.Handle("/", &templateHandler{filename: "chat.html"})

	// bootstrap 등 static html 부분을 위한 항목
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets/"))))

	// ch2: 주소를 chat 으로 바꾸고, 인증을 위해 MustAuth로 감싼다. 이러면 templateHanlder 는 인증이 되어야만 동작한다.
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"})) // MustAuth 를 통과하지 못하면, /login 으로 이동한다.
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)                                           // 새 룸을 개설
	http.Handle("/upload", &templateHandler{filename: "upload.html"}) // 아바타 사진 업로드
	http.HandleFunc("/uploader", uploaderHandler)
	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))

	// ch3: logout. auth.go 에서 SetCookie 로 저장한 쿠키를 초기화한다.
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	// get the room going
	// 룸을 실행. 무한 루프를 돌면서 상에 따라 select 구문을 실행함
	go r.run()

	// start the web server
	log.Println("String web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
	// 하드코딩된 앱 주소를 flag 로 변경함
	// if err := http.ListenAndServe(":8080", nil); err != nil {
	// 	log.Fatal("ListenAndServe:", err)
	// }

}
