package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"html"
	"log"
	"net/http"
)

const redisUrl = "redis:6379"
const counterKey = "Counter"

func dialRedis() redis.Conn {
	c, err := redis.Dial("tcp", redisUrl)
	if err != nil {
		log.Fatal("Error connecting to redis: ", err)
	}

	log.Println("Connected to redis")
	return c
}

func main() {
	//create a connection
	c := dialRedis()
	defer c.Close()

	http.HandleFunc("/count", handleCount(c))

	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Incoming request", r.URL)

		fmt.Fprintf(w, "Hello world, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the main page of %q", html.EscapeString(r.URL.String()))
	})

	log.Println("Listening for web connection")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleCount(c redis.Conn) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Incoming request", r.URL)

		//get the value from redis
		counter, err := redis.Int(c.Do("GET", counterKey))
		if err != nil {
			//failed to get, perhaps because never init, so init
			counter = 0
		}

		//increment and set the value to redis
		counter++
		_, err = c.Do("SET", counterKey, counter)
		if err != nil {
			fmt.Fprintf(w, "Failed to set value: %q", err.Error())
			return
		}

		//response with counter
		fmt.Fprintln(w, `	<html>
									<body>
										<h1>Counter (redis + web app)</h1>
										Count:`, counter, `<br>
										v0.0.4
									</body>
								</html>`)
	}
}
