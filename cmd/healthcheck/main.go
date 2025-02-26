/*
Healthcheck is a lightweight program designed to perform health checks by sending an HTTP request to the local host on a specified port.
This tool is particularly useful in environments where a simple health probe is needed, such as in Docker containers.
The default probe URL is http://localhost:3000/liveness, with the option to change only the port number.

Usage:

	healthcheck [flags]

The flags are:

	-port <1-65535>
		Specify the port number to which the healthcheck request will be sent.

Return values (exit codes):

	0
		The request was successful (HTTP 200 or HTTP 204).

	> 0
		The request was not successful (connection error or unexpected HTTP status code).
*/

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	var port = flag.Int("port", 3000, "HTTP port for healthcheck")
	flag.Parse()

	response, err := http.Get(fmt.Sprintf("http://localhost:%d/liveness", *port))

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	} else if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNoContent {
		response.Body.Close()
		fmt.Fprintln(os.Stderr, "healthcheck request failed with status: ", response.Status)
		os.Exit(1)
	}

	response.Body.Close()
	os.Exit(0)
}
