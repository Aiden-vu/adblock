package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"time"
)

type Cache struct {
	Name string  `json:"name"`
	Ip   string  `json:"ip"`
	Time float64 `json:"time"`
}

func main() {
	blocklistfile, err := os.Open("blocklist.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	var blocklist []string
	var test string
	// test = "00d84987c0.com"
	test = "codersports.com"
	scanner := bufio.NewScanner(blocklistfile)
	for scanner.Scan() {
		blocklist = append(blocklist, scanner.Text())
	}
	if slices.Contains(blocklist, test) {
		println("blocked")
		return
	}

	println(blocklist)

	start := time.Now()

	ip := QueryIp(test)

	elapsed := time.Since(start)

	//things to store: domain name, IP, how long request took
	cache := Cache{
		Name: test,
		Ip:   ip,
		Time: float64(elapsed.Milliseconds()) / 1000.0,
	}

	file, err := os.Create("cache.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(cache); err != nil {
		panic(err)
	}
	// Print the IP address
	fmt.Println(ip)
	fmt.Println(elapsed)
}

func QueryIp(domain_name string) string {
	url := "https://dns.google/resolve?name=" + domain_name + "&type=A"
	method := "GET"
	client := &http.Client{}
	error := "Error sending request:"
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return error
	}
	req.Header.Add("accept", "application/dns-json")

	resp, err := client.Do(req)
	if err != nil {
		return error
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return error
	}

	// Parse the JSON response
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Error parsing JSON response:", err)
		return error
	}

	// Extract the IP address from the first answer
	if data["Answer"] == nil {
		fmt.Println("No IP address found in response")
		return error
	}
	answer := data["Answer"].([]interface{})[0]
	ip := answer.(map[string]interface{})["data"].(string)

	return ip
}
