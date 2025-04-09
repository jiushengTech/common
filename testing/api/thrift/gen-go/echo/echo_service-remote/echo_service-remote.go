// Code generated by Thrift Compiler (0.19.0). DO NOT EDIT.

package main

import (
	"context"
	"flag"
	"fmt"
	thrift "github.com/apache/thrift/lib/go/thrift"
	"github.com/jiushengTech/common/testing/api/thrift/gen-go/echo"
	"math"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var _ = echo.GoUnusedProtection__

func Usage() {
	fmt.Fprintln(os.Stderr, "Usage of ", os.Args[0], " [-h host:port] [-u url] [-f[ramed]] function [arg1 [arg2...]]:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "\nFunctions:")
	fmt.Fprintln(os.Stderr, "  Response Echo(Request req)")
	fmt.Fprintln(os.Stderr, "  void VisitOneway(Request req)")
	fmt.Fprintln(os.Stderr)
	os.Exit(0)
}

type httpHeaders map[string]string

func (h httpHeaders) String() string {
	var m map[string]string = h
	return fmt.Sprintf("%s", m)
}

func (h httpHeaders) Set(value string) error {
	parts := strings.Split(value, ": ")
	if len(parts) != 2 {
		return fmt.Errorf("header should be of format 'Key: Value'")
	}
	h[parts[0]] = parts[1]
	return nil
}

func main() {
	flag.Usage = Usage
	var host string
	var port int
	var protocol string
	var urlString string
	var framed bool
	var useHttp bool
	headers := make(httpHeaders)
	var parsedUrl *url.URL
	var trans thrift.TTransport
	_ = strconv.Atoi
	_ = math.Abs
	flag.Usage = Usage
	flag.StringVar(&host, "h", "localhost", "Specify host and port")
	flag.IntVar(&port, "p", 9090, "Specify port")
	flag.StringVar(&protocol, "P", "binary", "Specify the protocol (binary, compact, simplejson, json)")
	flag.StringVar(&urlString, "u", "", "Specify the url")
	flag.BoolVar(&framed, "framed", false, "Use framed transport")
	flag.BoolVar(&useHttp, "http", false, "Use http")
	flag.Var(headers, "H", "Headers to set on the http(s) request (e.g. -H \"Key: Value\")")
	flag.Parse()

	if len(urlString) > 0 {
		var err error
		parsedUrl, err = url.Parse(urlString)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
			flag.Usage()
		}
		host = parsedUrl.Host
		useHttp = len(parsedUrl.Scheme) <= 0 || parsedUrl.Scheme == "http" || parsedUrl.Scheme == "https"
	} else if useHttp {
		_, err := url.Parse(fmt.Sprint("http://", host, ":", port))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing URL: ", err)
			flag.Usage()
		}
	}

	cmd := flag.Arg(0)
	var err error
	var cfg *thrift.TConfiguration = nil
	if useHttp {
		trans, err = thrift.NewTHttpClient(parsedUrl.String())
		if len(headers) > 0 {
			httptrans := trans.(*thrift.THttpClient)
			for key, value := range headers {
				httptrans.SetHeader(key, value)
			}
		}
	} else {
		portStr := fmt.Sprint(port)
		if strings.Contains(host, ":") {
			host, portStr, err = net.SplitHostPort(host)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error with host:", err)
				os.Exit(1)
			}
		}
		trans = thrift.NewTSocketConf(net.JoinHostPort(host, portStr), cfg)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error resolving address:", err)
			os.Exit(1)
		}
		if framed {
			trans = thrift.NewTFramedTransportConf(trans, cfg)
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating transport", err)
		os.Exit(1)
	}
	defer trans.Close()
	var protocolFactory thrift.TProtocolFactory
	switch protocol {
	case "compact":
		protocolFactory = thrift.NewTCompactProtocolFactoryConf(cfg)
		break
	case "simplejson":
		protocolFactory = thrift.NewTSimpleJSONProtocolFactoryConf(cfg)
		break
	case "json":
		protocolFactory = thrift.NewTJSONProtocolFactory()
		break
	case "binary", "":
		protocolFactory = thrift.NewTBinaryProtocolFactoryConf(cfg)
		break
	default:
		fmt.Fprintln(os.Stderr, "Invalid protocol specified: ", protocol)
		Usage()
		os.Exit(1)
	}
	iprot := protocolFactory.GetProtocol(trans)
	oprot := protocolFactory.GetProtocol(trans)
	client := echo.NewEchoServiceClient(thrift.NewTStandardClient(iprot, oprot))
	if err := trans.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening socket to ", host, ":", port, " ", err)
		os.Exit(1)
	}

	switch cmd {
	case "Echo":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "Echo requires 1 args")
			flag.Usage()
		}
		arg9 := flag.Arg(1)
		mbTrans10 := thrift.NewTMemoryBufferLen(len(arg9))
		defer mbTrans10.Close()
		_, err11 := mbTrans10.WriteString(arg9)
		if err11 != nil {
			Usage()
			return
		}
		factory12 := thrift.NewTJSONProtocolFactory()
		jsProt13 := factory12.GetProtocol(mbTrans10)
		argvalue0 := echo.NewRequest()
		err14 := argvalue0.Read(context.Background(), jsProt13)
		if err14 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.Echo(context.Background(), value0))
		fmt.Print("\n")
		break
	case "VisitOneway":
		if flag.NArg()-1 != 1 {
			fmt.Fprintln(os.Stderr, "VisitOneway requires 1 args")
			flag.Usage()
		}
		arg15 := flag.Arg(1)
		mbTrans16 := thrift.NewTMemoryBufferLen(len(arg15))
		defer mbTrans16.Close()
		_, err17 := mbTrans16.WriteString(arg15)
		if err17 != nil {
			Usage()
			return
		}
		factory18 := thrift.NewTJSONProtocolFactory()
		jsProt19 := factory18.GetProtocol(mbTrans16)
		argvalue0 := echo.NewRequest()
		err20 := argvalue0.Read(context.Background(), jsProt19)
		if err20 != nil {
			Usage()
			return
		}
		value0 := argvalue0
		fmt.Print(client.VisitOneway(context.Background(), value0))
		fmt.Print("\n")
		break
	case "":
		Usage()
		break
	default:
		fmt.Fprintln(os.Stderr, "Invalid function ", cmd)
	}
}
