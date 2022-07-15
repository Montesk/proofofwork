### Task for Server Engineer

Design and implement “Word of Wisdom” tcp server

- TCP server should be protected from DDOS attacks with the Prof of Work [POW - wiki](https://en.wikipedia.org/wiki/Proof_of_work), 
  the challenge-response protocol should be used
- The choice of the POW algorithm should be explained  
- After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes
- Docker file should be provided both for the server and for the client that solves the POW challenge

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
```json
{ "controller": "prove", "message": { "suggest": "c1419d1224dba805efb4c0397db229db747a56ea|bb6f6c336e94819f99a64b8ab3b03161a298be43|c1419d1224dba805efb4c0397db229db747a56ea" } }
```

Example success response
```json
{
  "action": "prove",
  "success": true,
  "message": "try again"
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
* Client ip address registered in the system after client login, server generates random nonce. Client IP address encrypted in SHA-1 in pattern [date|ip-address|nonce]. All parts are sha-1 encoded
* Server sends to client hash but without nonce part [date|ip-address]
* Client prepares header [date|ip-address|nonce] and try to suggest nonce (incrementally from 0)
* Client try to prove [date|ip-address|client_incr_nonce]