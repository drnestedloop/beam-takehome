## beam takehome

Welcome! The high level goal for this exercise is to build a client/server application in go-lang that synchronizes files from the client to the server (in one direction). A websocket server/client and basic protocol have been included.

### setup [v]

- install go
- install [air](https://github.com/cosmtrek/air) for hot reloading
- to start client: `make client`
- to start server: `make server`

### overview [v]

This repo includes a simple websocket client/server, with a simple protocol that enables them to send and receive serialized messages back and forth. A demo is included that shows a basic ECHO request that returns the same string that was sent by the client. Currently,
all the client does is infinitely send and receive the ECHO request.

The goal of this exercise to implement a very simple file synchronization protocol using the same protocol. For example, given an input directory, the client should scan that input directory, and serialize the files into messages containing the base64 contents of the file. On the server side, it should be able to handle that message and write the file to disk.

### goals [v]

- Design a SYNC request that can send a file over the websocket

- Implement a basic asynchronous 'file watcher' that takes in an input directory, and detects when any files change. When a file changes the client should send the entire file over the websocket using the SYNC request

- Implement a SYNC handler on the server side that can read in the SYNC request and write the file to disk. It should return a response saying whether this process was successful

## Completed Project
The project was completed according to the above specifications. 
### Changelog
- `/cmd/client/main.go` was updated to read the `SYNCFROM` environment variable (falling back on `../files` if `SYNCFROM` is not set) and watch the directory pointed to by it.
- `/pkg/client/client.go` was updated to include a Sync method for `Client`.
- `/pkg/common/types.go` was updated with abstractions for file operations, along with `SyncRequest` and `SyncResponse` structs.
- `/pkg/fileutils/fileutils.go` was created to house utility functions related to file watching, file serializing, and directory cleaning. 
- `/pkg/server/handlers.go` was updated to include a function for handling sync requests (`HandleSync`).
- `/pkg/server/server.go` was updated with functionality for reading the `SYNCTO` environment variable (see instructions below).

### Instructions for use
 Before running the server, make sure that either the folder `../sync` exists, or the environment `SYNCTO` is set to the path of the folder you wish to store the synchronized files to. Similarly, before running the client, make sure that either the folder `../files` is created, or the environment variable `SYNCFROM` is set to the path of the folder you wish to sync.