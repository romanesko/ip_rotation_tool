package main

import (
	"bufio"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

type conf struct {
	Port     string   `yaml:"port"`
	Endpoint string   `yaml:"endpoint"`
	IPs      []string `yaml:"ips"`
}

func (c *conf) getConf() *conf {
	yamlFile, err := os.ReadFile("conf.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func fileOp(ip string) error {

	if net.ParseIP(ip) == nil {
		return errors.New("invalid ip")
	}

	filename := "/etc/hosts"

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return err

	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "fapi.binance.com") {
			lines = append(lines, fmt.Sprintf("%s\t%s", ip, "fapi.binance.com"))
		} else if strings.Contains(line, "api.binance.com") {
			lines = append(lines, fmt.Sprintf("%s\t%s", ip, "api.binance.com"))
		} else {
			lines = append(lines, line)
		}
	}
	file.Close()

	file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		writer.WriteString(line + "\n")
	}
	writer.Flush()
	file.Close()
	return nil
}

var current = ""

func main() {

	var c conf
	c.getConf()

	fmt.Println("Listening on :" + c.Port + c.Endpoint)

	http.HandleFunc(c.Endpoint, func(w http.ResponseWriter, r *http.Request) {
		selected := r.URL.Query().Get("ip")

		if selected != "" {
			fileWriteERr := fileOp(selected)

			if fileWriteERr != nil {
				fmt.Fprintf(w, "Error: %s", fileWriteERr)
				return
			}
			current = selected
			http.Redirect(w, r, c.Endpoint, http.StatusSeeOther)

		} else {
			selected = current
		}

		fmt.Fprintf(w, `<html><style>.active {font-weight:bold} a {color: blue;}</style><pre>`)
		for _, ip := range c.IPs {
			cls := ""
			if selected == ip {
				cls = "active"
			}
			fmt.Fprintf(w, "<a href=\"?ip=%s\" class=\"%s\">%s</a>\n", ip, cls, ip)
		}
		fmt.Fprintf(w, `</pre></html>`)

	})

	log.Fatal(http.ListenAndServe(":"+c.Port, nil))
}
