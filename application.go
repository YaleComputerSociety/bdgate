package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"time"

	"./util"

	"github.com/asaskevich/govalidator"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
)

var redisDB int64 = 4

const redisSchemaVersion = 0

const indexPage = "public/index.html"

const (
	casUrl1 = "https://secure.its.yale.edu/cas/login?service="
	casUrl2 = "https://secure.its.yale.edu/cas/serviceValidate?"
)

var credis redis.Conn

func HandleGetIndex(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving %s to %s...\n", indexPage, r.RemoteAddr)
	http.ServeFile(w, r, indexPage)
}

func genNewId() int64 {
	// Atomically get last used id and increase count.
	id, err := redis.Int64(credis.Do("INCR", "global:lastId"))

	if err != nil {
		panic("Redigo failed to INCR global:lastId.")
	}

	return id
}

func printError(r interface{}) {
	log.Printf("Error: %s.\n", r)

	trace := make([]byte, 1024)
	runtime.Stack(trace, true)
	log.Printf("Stack: %s\n", trace)
}

func HandlePostUrl(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			http.Error(w, "Application failure.", 500)
			printError(r)
		}
	}()

	r.ParseForm()

	if _, ok := r.Form["url"]; ok == false {
		http.Error(w, "Invalid request.", 400)
		return
	}

	url := r.Form["url"][0]

	if !govalidator.IsURL(url) {
		http.Error(w, "Invalid url format.", 400)
		return
	}

	// TODO: check if this url has already been stored?

	// READ ME:

	// that would be collisions (lol, not happening), and the inability
	// to have the same url being inserted more than once in the system.
	var id util.UUId
	for {
		id = *util.NewUUId(genNewId())
		if r, err := credis.Do("GET", "urls:"+id.String()); err != nil {
			panic("ERROR: Failed to Get in redis client.\n")
		} else if r == nil {
			// id not in urlsIds:* yet. Perfect!
			break
		}

		log.Printf("WARN: generated already existing id. Trying again.\n")
	}

	// Store to database.
	skey := fmt.Sprintf("urls:%d", id.Int64())
	if _, err := credis.Do("HMSET", skey,
		"user", "foo",
		"base58", id.Base58(),
		"url", url,
		"timestamp", time.Now().Unix(),
		"version", redisSchemaVersion); err != nil {
		panic("Failed to HSET.\n" + err.Error())
	}

	skey = fmt.Sprintf("shorts:%s", id.Base58())
	if _, err := credis.Do("SET", skey, id.Int64()); err != nil {
		panic("Failed to SET.\n" + err.Error())
	}

	log.Printf("Received url=%s id=%d '%s'", url, id, id.Base58())

	short := "/urls/" + id.Base58()
	message := "Hey! Your short is <a href=\"" + short + "\">" + short + "</a>"
	w.Write([]byte(message))
}

func HandleGetUrl(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			http.Error(w, "Application failure.", 500)
			printError(r)
		}
	}()

	base58, ok := mux.Vars(r)["key"]
	if ok == false {
		http.Error(w, "Failed to find url key.", 404)
		return
	}

	var cbUrl = "http://localhost:5000/urls/" + base58 + "/callback"
	http.Redirect(w, r, casUrl1+url.QueryEscape(cbUrl), 301)
}

func validateTicket(ticket, service string) (bool, error) {
	v := url.Values{}
	v.Set("ticket", ticket)
	v.Set("service", service)
	// Setting ?format=JSON is not working.

	// Reach CAS server to validate user ticket.
	log.Printf("Reaching %s.\n", casUrl2+v.Encode())
	resp, err := http.Get(casUrl2 + v.Encode())
	if err != nil {
		panic("Failed to GET ticket verification endpoint.\n" + err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	fmt.Printf("body: %s\n", body)

	return false, nil
}

func HandleGetUrlCallback(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			http.Error(w, "Application failure.", 500)
			printError(r)
		}
	}()

	ticket := r.URL.Query().Get("ticket")
	if ticket == "" {
		http.Error(w, "Failed to find ticket string.", 400)
		return
	}

	base58, ok := mux.Vars(r)["key"]
	if ok == false {
		http.Error(w, "Failed to find url key.", 404)
		return
	}
	var service = "http://localhost:5000/urls/" + base58 + "/callback"

	isValid, err := validateTicket(ticket, service)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Valid: %v\n", isValid)

	return

	id, err := util.GenIdFromBase58(base58)
	if err != nil {
		http.Error(w, "Invalid url.", 400)
		return
	}

	var p2 struct {
		Timestamp int64 `redis:"timestamp"`
		Version   int
		Url       string `redis:"url"`
		Base58    string
		User      string
	}

	skey := fmt.Sprintf("urls:%d", id.Int64())
	reply, err := redis.Values(credis.Do("HGETALL", skey))
	if err == redis.ErrNil || len(reply) == 0 {
		http.Error(w, "Sorry. No url exists for this.", 404)
		panic("urls:. hash not found for id " + id.String())
	} else if err != nil {
		panic("Failed to HGET.\n" + err.Error())
	}

	if err = redis.ScanStruct(reply, &p2); err != nil {
		panic("Failed to scanstruct urls:" + id.String())
	}

	fmt.Printf("%+v", p2)

	http.Redirect(w, r, p2.Url, 301)
}

func route() *mux.Router {
	r := mux.NewRouter().Schemes("http").Subrouter()

	r.HandleFunc("/", HandleGetIndex).Methods("GET")
	r.HandleFunc("/urls", HandlePostUrl).Methods("POST")
	r.HandleFunc("/urls/{key}", HandleGetUrl).Methods("GET")
	r.HandleFunc("/urls/{key}/callback", HandleGetUrlCallback).Methods("GET")

	//http.HandleFunc("/scheduled", func(w http.ResponseWriter, r *http.Request) {
	//if r.Method == "POST" {
	//log.Printf("Received task %s scheduled at %s\n", r.Header.Get("X-Aws-Sqsd-Taskname"), r.Header.Get("X-Aws-Sqsd-Scheduled-At"))
	//}
	//})

	return r
}

func setupRedis() redis.Conn {
	client, err := redis.Dial("tcp", ":6379", redis.DialDatabase(4))
	if err != nil {
		panic("Failed to start redis.\n" + err.Error())
	}

	log.Println("Redis started.")

	return client
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	credis = setupRedis()
	defer credis.Close()
	f, _ := os.Create("/var/log/golang/golang-server.log")
	defer f.Close()
	// log.SetOutput(f)

	http.Handle("/", route())

	log.Printf("Listening on port %s\n\n", port)
	http.ListenAndServe(":"+port, nil)
}
