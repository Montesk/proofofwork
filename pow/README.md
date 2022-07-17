### Package structure

- main.go - the entrypoint and starting the integration test for prove challenge (network based, required main app started)
- package:client - wrapper for proof of work interface calls
- package:pow - interface of pow and encoding mechanism
- package:server - networked client implementation of pow
- package:service - generate and prove pow with mem-cache