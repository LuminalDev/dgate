package discord

import (
	"fmt"
	"regexp"

	"github.com/valyala/fasthttp"
)

var (
	JS_FILE_REGEX    = regexp.MustCompile(`assets/+([a-z0-9.]+)\.js`)
	BUILD_INFO_REGEX = regexp.MustCompile(`Build Number: "\)\.concat\("([0-9]{4,8})"`)
)

func getLatestBuild() (string, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.Header.SetMethod(fasthttp.MethodGet)
	req.SetRequestURI("https://discord.com/app")

	if err := requestClient.Do(req, resp); err != nil {
		return "", err
	}

	matches := JS_FILE_REGEX.FindAllStringSubmatch(string(resp.Body()), -1)
	fmt.Println(string(resp.Body()))
	asset := matches[len(matches)-10][1]

	req.Header.SetMethod(fasthttp.MethodGet)
	req.SetRequestURI(fmt.Sprintf("https://discord.com/assets/%s.js", asset))

	if err := requestClient.Do(req, resp); err != nil {
		return "", err
	}

	match := BUILD_INFO_REGEX.FindStringSubmatch(string(resp.Body()))
	return match[1], nil
}

func mustGetLatestBuild() string {
	if build, err := getLatestBuild(); err != nil {
		panic(err)
	} else {
		return build
	}
}
