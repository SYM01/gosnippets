package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/valyala/fasthttp"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type fetcher struct {
	method  string
	url     string
	headers []string
	body    []byte

	ignoreSSLErr bool
	showHeader   bool
	showContent  bool
	output       *os.File
}

func (f *fetcher) reqMethod() string {
	if f.method == "" {
		return "GET"
	}
	return f.method
}

func (f *fetcher) client() *fasthttp.Client {
	c := new(fasthttp.Client)
	if f.ignoreSSLErr {
		c.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	return c
}

func (f *fetcher) Run() error {
	request := new(fasthttp.Request)
	request.Header.SetMethod(f.reqMethod())
	for _, item := range f.headers {
		header := strings.SplitN(item, ":", 2)
		request.Header.Add(header[0], strings.TrimSpace(header[1]))
	}
	request.SetRequestURI(f.url)
	request.SetBodyRaw(f.body)

	resp := new(fasthttp.Response)
	if !f.showContent {
		resp.SkipBody = true
	}

	if err := f.client().Do(request, resp); err != nil {
		return err
	}

	resp.WriteTo(f.output)
	return nil
}

func main() {
	var (
		ignoreSSLErr  = flag.Bool("k", false, "Allow insecure server connections when using SSL")
		method        = flag.String("X", "", "Specify request command to use")
		head          = flag.Bool("I", false, "Show document info only")
		includeHeader = flag.Bool("i", false, "Include protocol response headers in the output")
		postData      = flag.String("d", "", "HTTP POST data")
		output        = flag.String("o", "", "Write to file instead of stdout")
		headers       arrayFlags
		err           error
	)
	flag.Var(&headers, "H", "Pass custom header(s) to server")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options...] <url>\n", filepath.Base(os.Args[0]))
		return
	}

	f := &fetcher{
		method:       strings.ToUpper(*method),
		url:          flag.Arg(0),
		headers:      headers,
		ignoreSSLErr: *ignoreSSLErr,
		showHeader:   *includeHeader,
		showContent:  true,
		output:       os.Stdout,
	}

	if *head {
		f.showHeader = true
		f.showContent = false
		if f.method == "" {
			f.method = "HEAD"
		}
	}

	if *postData != "" {
		f.body = []byte(*postData)
		hasContentType := false
		for _, header := range f.headers {
			if strings.HasPrefix(strings.ToLower(header), "content-type:") {
				hasContentType = true
				break
			}
		}

		if !hasContentType {
			f.headers = append(f.headers, "Content-Type: application/x-www-form-urlencoded")
		}
	}

	if *output != "" {
		f.output, err = os.OpenFile(*output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
		if err != nil {
			log.Fatal(err)
		}
		defer f.output.Close()
	}

	if err = f.Run(); err != nil {
		log.Fatal(err)
	}
}
