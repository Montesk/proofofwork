### Task for Server Engineer

Design and implement “Word of Wisdom” tcp server

- TCP server should be protected from DDOS attacks with the Prof of Work [POW - wiki](https://en.wikipedia.org/wiki/Proof_of_work), 
  the challenge-response protocol should be used
- The choice of the POW algorithm should be explained  
- After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes
- Docker file should be provided both for the server and for the client that solves the POW challenge

### Build
`docker build -t proofofwork .`

### Run application
`docker run -p <port>:8001 proofofwork`
e.g.
`docker run -p 9010:8001 proofofwork` runs the application on `9010` of the caller system

full example `docker build -t proofofwork . && docker run -p 9010:8001 proofofwork`

### Send TCP message
`netcat 127.0.0.1 9010` establish connection

then use stubs listed below, protocol mutually json encoded 

### Message stubs
#### Challenge controller
```json
{ "controller": "challenge" }
```

Example response
```json
{ "action":  "challenge", "message": "c1419d1224dba805efb4c0397db229db747a56ea|bb6f6c336e94819f99a64b8ab3b03161a298be43" }
```

#### Prove controller

Prove suggest must append one of the number encoded in SHA-1, in range [1-10]
Use [Online sha-1 encode tool](http://www.sha1-online.com/) for manual testing

```json
{ "controller": "prove", "message": { "suggest": "c1419d1224dba805efb4c0397db229db747a56ea|bb6f6c336e94819f99a64b8ab3b03161a298be43|c1419d1224dba805efb4c0397db229db747a56ea" } }
```

Example success response
```json
{
  "action": "prove",
  "success": true,
  "message": "You could make it a singleton too, but friends don’t let friends create singletons.\" ― Robert Nystrom"
}
```

Example failed response
```json
{
  "action": "prove",
  "success": false,
  "message": "try again"
}
```

### POW Algorithm
* Client ip address registered in the system after client login, server generates random nonce. Client IP address encrypted in SHA-1 in pattern `[date|ip-address|nonce]`. All parts are sha-1 encoded
* Server sends to client hash but without nonce part `[date|ip-address]`
* Client prepares header `[date|ip-address|nonce]` and try to suggest nonce (incrementally from 0)
* Client try to prove `[date|ip-address|client_incr_nonce]`