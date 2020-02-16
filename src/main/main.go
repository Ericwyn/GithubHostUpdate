package main

import (
	"ajax"
	"ezefile"
	"ezestr"
	"fmt"
	"regexp"
	"strings"
)

var githubHosts = map[string]string{}

func main() {
	readGithubHostList()
	hostBac, hostNew := readHost()
	getGithubHost()
	for k, v := range githubHosts {
		if v != "" {
			hostNew = append(hostNew, v+" "+k+"\n")
		}
	}
	backupHost(hostBac)
	//fmt.Println(hostNew)
	writeHost(hostNew)
}

func getGithubHost() {
	for k := range githubHosts {
		ip := getOneHost(k)
		if ip != "" {
			fmt.Println("get host:", k, ",ip:", ip)
			githubHosts[k] = ip
		}
	}
}

func getOneHost(host string) string {
	var ip = ""
	ajax.Send(ajax.Request{
		Method: ajax.GET,
		Url:    "https://www.ipaddress.com/search/" + host,
		Header: map[string]string{
			"Referer":    "https://www.ipaddress.com/",
			"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36",
		},
		Success: func(response *ajax.Response) {
			//body := response.Body
			// 正则提取 ： <th>IP Address</th><td><ul class="comma-separated"><li>.+?</li>
			ip = getHostFromHtmlResp(response.Body)
		},
		Fail: func(status int, errMsg string) {
			fmt.Println(errMsg)
		},
	})
	return ip
}

var htmlResReg, _ = regexp.Compile(`<th>IP Address</th><td><ul class="comma-separated"><li>.+?</li>`)

func getHostFromHtmlResp(html string) string {
	ip := htmlResReg.FindString(html)
	//fmt.Println(ip)
	if ip != "" {
		ip = strings.Replace(ip, "</li>", "", -1)
		ip = strings.Replace(ip, `<th>IP Address</th><td><ul class="comma-separated"><li>`, "", -1)
		if ezestr.Continues(ip, ",") {
			return strings.Split(ip, ",")[0]
		}
		return ip
	}
	return ""
}

func readGithubHostList() {
	ezefile.ReadLine("hosts.list", func(line string) {
		line = ezestr.Remove(line, []string{" ", "\n"})
		githubHosts[line] = ""
	})
}

func readHost() ([]string, []string) {
	hostsNew := make([]string, 0)
	hostsBac := make([]string, 0)
	ezefile.ReadLine("/etc/hosts", func(line string) {
		//fmt.Print(line)
		hostsBac = append(hostsBac, line)
		appenFlag := true
		for k := range githubHosts {
			// 如果这一行包含需要设置的域名 “ github.com” 并且不是以 “#” 开头，那么这一行就要跳过
			if ezestr.Continues(line, " "+k) && !ezestr.StartWith(line, "#") {
				appenFlag = false
			}
		}
		if appenFlag {
			hostsNew = append(hostsNew, line)
		}
	})
	return hostsBac, hostsNew
}

func writeHost(hosts []string) {
	ezefile.Write("/etc/hosts", hosts)
}

func backupHost(hosts []string) {
	ezefile.Write("/etc/hosts.bac", hosts)
}
