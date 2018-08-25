package chat

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

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

	// 템플릿에 request 정보를 전달
	// 자바스크립트 부분에서 {{.Host}} 라는 형식으로 사용할 수 있음
	t.templ.Execute(w, r)
	// t.templ.Execute(w, nil)
}

func main() {

	// 채팅 사이트 주소가 하드코딩됨 (localhost:8080)
	// 이를 커맨드라인에서 -addr 라는 플래그로 처리하도록 변경 ./chat -addr=":3000" 이라는 형식으로 실행이 가능해짐
	var addr = flag.String("addr", ":8080", "The addr of the application")
	flag.Parse() // parse the flags

	// newRoom 함수로 새 룸을 만든다.
	r := newRoom()

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
	http.Handle("/", &templateHandler{filename: "chat.html"})

	// QST: 이 부분은 도대체 무엇을 하나?
	http.Handle("/room", r)

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
