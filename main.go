package main

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang_study/pkg/metrics"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

var logFile *os.File

func init() {
	viper.SetConfigFile("config.toml")
	erro := viper.ReadInConfig()
	if erro != nil { // 读取配置信息失败
		panic(fmt.Errorf("fatal error config file: %s", erro))
	}
	logrus.SetLevel(logrus.Level(viper.GetInt("log.level")))
	logFile, err := os.OpenFile(viper.GetString("log.path"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("Failed to open log file")
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	logrus.SetOutput(mw)
}

func closeLogFile() {
	if logFile != nil {
		logrus.Infoln("closing log file")
		logFile.Close()
	}
}

func randInt(max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max)
}
func delay() {
	delayTime := randInt(2000)
	time.Sleep(time.Millisecond * time.Duration(delayTime))
	logrus.Infoln("请求延时", delayTime, "ms完成。")
}
func main() {
	defer closeLogFile()
	metrics.Register()
	mux := mux.NewRouter()
	mux.Handle("/", &myHandler{})
	mux.HandleFunc("/version", version)
	mux.HandleFunc("/healthZ", healthZ)
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/{url:.*}", err)

	server := &http.Server{
		Addr:         ":" + viper.GetString("http_server.port"),
		WriteTimeout: time.Second * viper.GetDuration("http_server.timeout"),
		Handler:      mux,
	}
	go func() {
		logrus.Infoln("HTTP服务启动", "http://localhost"+server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Errorln(err)
			os.Exit(0)
		}
		logrus.Infoln("HTTP服务关闭请求")
	}()
	// 监听信号，优雅退出http服务
	Watch(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(ctx)
	})
	logrus.Infoln("程序退出")
}

type myHandler struct {
}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	delay()
	for k, v := range r.Header {
		for _, vv := range v {
			w.Header().Set(k, vv)
		}
	}
	w.Write([]byte("this is lvzhancheng http server"))
	logrus.Infoln(Clien_IP(r), r.URL, r.Method, 200)
}
func version(w http.ResponseWriter, r *http.Request) {
	delay()
	v, exists := os.LookupEnv("VERSION")
	if exists {
		w.Header().Add("version", v)
		w.Write([]byte("VERSION:" + v))
	} else {
		os.Setenv("VERSION", "0.0.1")
		w.Header().Add("version", "0.0.1")
		w.Write([]byte("VERSION: 0.0.1"))
	}
	logrus.Infoln(Clien_IP(r), r.URL, r.Method, 200)
}
func healthZ(w http.ResponseWriter, r *http.Request) {
	delay()
	w.WriteHeader(http.StatusOK)
	logrus.Infoln(Clien_IP(r), r.URL, r.Method, 200)
}
func err(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if len(vars) != 0 {
		w.WriteHeader(404)
		w.Write(([]byte("404 page not found")))
		logrus.Infoln(Clien_IP(r), r.URL, r.Method, 404)
	}
}

// 监听信号
func Watch(fns ...func() error) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	//阻塞
	s := <-ch
	close(ch)
	logrus.Infoln("接收到信号", s.String(), "执行关闭函数")
	for i := range fns {
		if err := fns[i](); err != nil {
			logrus.Errorln(err)
		}
	}
	logrus.Infoln("关闭函数执行完成")
}
