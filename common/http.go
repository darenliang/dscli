package common

import (
	"net/http"
	"time"
)

// HttpClient is used to download file chunks from Discord
var HttpClient = http.Client{Timeout: 60 * time.Second}
