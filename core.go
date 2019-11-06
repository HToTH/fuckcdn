package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func Getip() {
}
func CheckPortOpen(ip string) bool {

	_, err := net.DialTimeout("tcp", ip, time.Second*1)
	if err != nil {
		return false
	}
	return true
}
func FindWord(ip, value string) bool {
	return strings.Contains(ip, value)
}
func (data *Data) GetHttpResponse(ip string, c chan string) string {
	if CheckPortOpen(ip) == false {
		c <- "0"
		return "0"
	}
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		//fmt.Println("error dialing", err.Error())
		c <- "0"
		return "0"
	}
	msg := "GET / HTTP/1.1\r\n"
	msg += "Host:" + data.domain + "\r\n"
	msg += "DNT: 1\r\n"
	msg += "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.87 Safari/537.36\r\n"
	msg += "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3\r\n"
	// msg += "Accept-Encoding: gzip, deflate\r\n"
	msg += "Accept-Language: zh-CN,zh;q=0.9,en;q=0.8\r\n"
	msg += "Connection: close\r\n"
	msg += "\r\n\r\n"
	defer conn.Close()
	_, err = io.WriteString(conn, msg)
	if err != nil {
		fmt.Println("write string failed", err)
		return "0"
	}

	buf := make([]byte, 4096)
	count, err := conn.Read(buf)

	if err != nil {
		//fmt.Println("read string failed", err)
		c <- "0"
		return "0"
	}
	if FindWord(string(buf[0:count]), data.value) {
		tmp := ":"
		if FindWord(string(buf[0:count]), "CF-Cache-Status") {
			tmp = ":cloudFlare Cdn"
		}
		ip = ip + tmp
		c <- ip
		return ip
	} else {
		c <- "0"
		return "0"
	}
}

//将10进制IP转成正常IP格式
func TenToNormalIp(ip int64) string {
	s2 := strconv.FormatInt(ip, 2) //10 yo 2
	s3 := string(s2)               //tostring
	s := ""                        //填充字符串
	if len(s3) < 32 {
		for i := 1; i < 32-len(s3)+1; i++ {
			s = s + "0"
		}
	}
	s3 = s + s3
	//分割二进制变成10进制int64
	b1, _ := strconv.ParseUint(s3[0:8], 2, 8)
	b2, _ := strconv.ParseUint(s3[8:16], 2, 8)
	b3, _ := strconv.ParseUint(s3[16:24], 2, 8)
	b4, _ := strconv.ParseUint(s3[24:32], 2, 8)
	//10进制int64转string
	d1 := strconv.FormatUint(b1, 10)
	d2 := strconv.FormatUint(b2, 10)
	d3 := strconv.FormatUint(b3, 10)
	d4 := strconv.FormatUint(b4, 10)
	ips := d1 + "." + d2 + "." + d3 + "." + d4
	return ips
}

func (data *Data) Start(port int64, i int, c chan string, thread int, end int) {
	p_temp := strconv.FormatInt(port, 10)
	d := i
	for ; i < d+thread; i++ {
		ip := TenToNormalIp(int64(i))
		ips := ip + ":" + p_temp
		go data.GetHttpResponse(ips, c)
	}
	if i < end {
		time.Sleep(time.Millisecond * 1000)
		data.Start(port, i, c, thread, end)
	}
}

func (data *Data) ReciveMessage(c chan string) {
	i := 0.0
	d := 4278124286.0
	for {
		select {
		case res := <-c:
			if res != "0" {
				dd := strings.Split(res, ":")
				f, err := os.OpenFile("fuckcdn.log", os.O_APPEND|os.O_RDWR, 0777)
				if err != nil {
					fmt.Printf("open err%s", err)
					return
				}
				s := fmt.Sprintf("ip:%s,port:%s,domain:%s,关键字:%s,备注:%s\n", dd[0], dd[1], data.domain, data.value, dd[2])
				f.WriteString(s)
				f.Close()
				fmt.Printf("ip:%s,port:%s,domain:%s,关键字:%s,备注:%s\n", dd[0], dd[1], data.domain, data.value, dd[2])
			} else {
				fmt.Printf("\r%.8f%", i/d*100)
			}
		case <-time.After(time.Second * 1):
			fmt.Printf("\r%.8f%", i/d*100)
		}
		i++
	}
}

func NormalIpToten(ip string) int {
	ips := strings.Split(ip, ".")
	if len(ips) == 4 {
		b1, err := strconv.Atoi(ips[0])
		b2, err := strconv.Atoi(ips[1])
		b3, err := strconv.Atoi(ips[2])
		b4, err := strconv.Atoi(ips[3])
		if err != nil || checkNum(b1, 3) == -1 || checkNum(b2, 2) == -1 || checkNum(b3, 1) == -1 || checkNum(b4, 0) == -1 {
			fmt.Print("IP地址错误")
			return -1
		}
		return checkNum(b1, 3) + checkNum(b2, 2) + checkNum(b3, 1) + checkNum(b4, 0)
	} else {
		fmt.Print("IP地址错误")
		return -1
	}
}
func checkNum(i, d int) int {
	cc := 1 * i
	if i >= 0 && i <= 255 {
		for c := 0; c < d; c++ {
			cc = 256 * cc
		}
		return cc
	} else {
		return -1
	}
}
