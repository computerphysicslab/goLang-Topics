package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/patrickmn/go-cache"
)

// A cache to avoid repeated API calls to github
var myCache *cache.Cache

func logfileCreate() {
	logFile, err := os.OpenFile("logfile.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(logFile)
	// t := time.Now()
	// currentTime = t.Format("2020-01-02 15:04:05")
	log.Println("Log created/open")
}

func cacheInit() {
	// Load serialized cache from file if exists
	b, err := ioutil.ReadFile("cachePersistent.dat")
	if err != nil {
		// panic(err)
		myCache = cache.New(5*time.Minute, 10*time.Minute)
	}

	// Deserialize
	decodedMap := make(map[string]cache.Item, 500)
	d := gob.NewDecoder(bytes.NewBuffer(b))
	err = d.Decode(&decodedMap)
	if err != nil {
		panic(err)
	}

	myCache = cache.NewFrom(5*time.Minute, 10*time.Minute, decodedMap)
}

func getAPIGitHubJSON(path string) string {
	// Get API result from cache if available
	b, found := myCache.Get(path)
	if found {
		bodyString := b.(string)
		log.Println("cached", path, bodyString)
		return bodyString
	}

	// To avoid API quota limit
	githubDelay := 720 * time.Millisecond // 5000 QPH = 83.3 QPM = 1.38 QPS = 720ms/query
	time.Sleep(githubDelay)

	req, err := http.NewRequest("GET", "https://api.github.com"+path, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Accept", "application/vnd.github.mercy-preview+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("fromAPI", path, string(body))
	myCache.Set(path, string(body), cache.NoExpiration) // Store API result in cache

	return (string(body))
}

func getTopicsFromRepository(repo string) []string {
	// Querying the API
	jsonResult := getAPIGitHubJSON("/repos/" + repo + "/topics")

	// Content structure
	var result map[string][]string

	// Unmarshal or Decode the JSON to the map
	err := json.Unmarshal([]byte(jsonResult), &result)
	if err != nil {
		log.Fatalln(err)
	}

	return result["names"]
}

func getReposWithTopic(topic string) map[string]interface{} {
	// Querying the API
	jsonResult := getAPIGitHubJSON("/search/repositories?q=topic:" + topic + "&sort=stars&order=desc")

	// Result structure
	var result map[string]interface{}

	// Unmarshal or Decode the JSON to the interface
	err := json.Unmarshal([]byte(jsonResult), &result)
	if err != nil {
		log.Fatalln(err)
	}

	return result
}

// freqSort is used to sort a map indexed w/ strings by its integer value in descending order
// kv is the struct of (key, value) to output of the sort function
type kv struct {
	Key   string
	Value int
}

func freqSort(values map[string]int) []kv {
	var ss []kv
	for k, v := range values {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	return ss
}

func main() {
	logfileCreate()
	cacheInit()

	reposFound := getReposWithTopic("traefik")
	fmt.Printf("Repositories found: %v\n\n", reposFound["total_count"])

	items := reposFound["items"].([]interface{})

	// Create map to store topic frequencies
	topicFreq := make(map[string]int)

	for key, result := range items[:] {
		item := result.(map[string]interface{})
		// fullName := result["full_name"].(string)
		// fmt.Println("Reading Value for Key :", key)
		// Reading each value by its key
		fmt.Println(key, " ", item["id"],
			" || ", item["full_name"],
			" || ", item["stargazers_count"],
			" || ", item["description"])

		topics := getTopicsFromRepository(item["full_name"].(string))

		// Loop for every topic in the array
		for _, topic := range topics {
			// Increase topicFreq counter
			topicFreq[topic]++
		}
	}

	topicBest := freqSort(topicFreq)[:10]

	fmt.Printf("\nRESULT: %v\n", topicBest)

	// fmt.Printf("DEBUG: %v\n", myCache.Items())

	// Store cache into persistent file
	// Serialize cache
	b := new(bytes.Buffer)
	e := gob.NewEncoder(b)
	// Encoding the map
	err := e.Encode(myCache.Items())
	if err != nil {
		panic(err)
	}

	// Save serialized cache into file
	err = ioutil.WriteFile("cachePersistent.dat", b.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}
