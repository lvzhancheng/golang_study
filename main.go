package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func Clien_IP(r *http.Request) string {
	IpAddress := r.Header.Get("X-Real-Ip")
	if IpAddress == "" {
		IpAddress = r.Header.Get("X-Forwarded-For")
	}
	if IpAddress == "" {
		IpAddress = r.RemoteAddr
	}
	return IpAddress
}

func main() {
	mux := mux.NewRouter()
	mux.Handle("/", &myHandler{})
	mux.HandleFunc("/version", version)
	mux.HandleFunc("/healthZ", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Fprintf(w, "200")
		w.WriteHeader(http.StatusOK)
		log.Printf("%s %s %s %d", Clien_IP(r), r.URL, r.Method, 200)
	})
	mux.HandleFunc("/{url:.*}", err)

	server := &http.Server{
		Addr:         ":80",
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}
	log.Println("Starting my httpserver")
	log.Fatal(server.ListenAndServe())
	// 监听信号，优雅退出http服务
	Watch(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(ctx)
	})
}

type myHandler struct {
}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		for _, vv := range v {
			w.Header().Set(k, vv)
		}
	}
	w.Write([]byte("this is lvzhancheng http server"))
	log.Printf("%s %s %s %d", Clien_IP(r), r.URL, r.Method, 200)
}
func version(w http.ResponseWriter, r *http.Request) {
	v, exists := os.LookupEnv("VERSION")
	if exists {
		w.Header().Add("version", v)
		w.Write([]byte("VERSION:" + v))
	} else {
		os.Setenv("VERSION", "0.0.1")
		w.Header().Add("version", "0.0.1")
		w.Write([]byte("VERSION: 0.0.1"))
	}
	log.Printf("%s %s %s %d", Clien_IP(r), r.URL, r.Method, 200)
}
func err(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if len(vars) != 0 {
		w.WriteHeader(404)
		w.Write(([]byte("404 page not found")))
		log.Printf("%s %s %s %d", Clien_IP(r), r.URL, r.Method, 404)
	}
}

// 监听信号
func Watch(fns ...func() error) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	//阻塞
	s := <-ch
	close(ch)
	log.Println("接收到信号", s.String(), "执行关闭函数")
	for i := range fns {
		if err := fns[i](); err != nil {
			log.Println(err)
		}
	}
	log.Println("关闭函数执行完成")
}
