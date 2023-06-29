# Golang API Example

Copyright (c) 2021-2023 Tenebris Technologies Inc.

See the LICENSE file for further information.

This is an example of how to implement a production HTTP/HTTPS server in Golang using github.com/gorilla/mux 
and limit the number of concurrent requests.

The API server itself is implemented as a package so that you can drop it into your own application.

A log wrapper example is also included. It currently writes logs to stdout using fmt.Printf and can be updated to send
log events to another destination (see api/logger.go).

HTTP routes are defined in routes.go.

Three example routes are defined:

/heath is for health checking. If the "DownFile" specified in main.go exists, the health check will fail. This can be
used as a mechanism to remove a node from a load balancer target group prior to maintenance.

/example will respond indicating no id received.

/example/<id> will simply reflect back the specified id.

For assistance developing your project, please contact us via
https://tenebris.com
