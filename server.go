/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : server.go

* Purpose : Server ip check api for bird

* Creation Date : 04-04-2014

* Last Modified : Fri 04 Apr 2014 10:15:58 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/kiyor/parsebird/lib"
	"github.com/kiyor/subnettool"
	"net"
	"strings"
)

func main() {
	routes := getIP()
	j, _ := json.MarshalIndent(routes, "", "    ")
	fmt.Println(string(j))
	runHttp(routes)
}

func getIP() []parsebird.Route {
	var routes []parsebird.Route
	ips, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, v := range ips {
		ipstr := strings.Split(v.String(), "/")[0]
		ip := net.ParseIP(ipstr)
		if ip.To4() != nil {
			token := subnettool.ParseIPInt(ip.To4())
			if token[0] == 127 || token[0] == 10 {
				continue
			}
			var r parsebird.Route
			r.Ip = ip
			routes = append(routes, r)
		}
	}
	return routes
}

func runHttp(routes []parsebird.Route) {
	m := martini.Classic()
	m.Get("/serverinfo/ip/:ip", func(params martini.Params) string {
		return fmt.Sprintf("%v", checkExist(routes, params["ip"]))
	})
	m.Run()
}

func checkExist(routes []parsebird.Route, str string) bool {
	for _, v := range routes {
		if v.Ip.String() == str {
			return true
		}
	}
	return false
}
