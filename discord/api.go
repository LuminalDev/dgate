package discord

import (
	"github.com/valyala/fasthttp"
	"time"
)

var (
	clientBuildNumber = mustGetLatestBuild()
	requestClient     = fasthttp.Client{
		ReadBufferSize:                8192,
		ReadTimeout:                   time.Second * 5,
		WriteTimeout:                  time.Second * 5,
		NoDefaultUserAgentHeader:      true,
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        true,
	}
)
