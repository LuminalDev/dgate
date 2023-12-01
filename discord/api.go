package discord

import (
	"time"

	"github.com/valyala/fasthttp"
)

var (
	clientBuildNumber = "250513"
	clientLocale      = mustGetLocale()
	requestClient     = fasthttp.Client{
		ReadBufferSize:                8192,
		ReadTimeout:                   time.Second * 5,
		WriteTimeout:                  time.Second * 5,
		NoDefaultUserAgentHeader:      true,
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        true,
	}
)
