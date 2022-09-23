# UDP example

This example demonstrates sending and receiving messages over UDP.
This example consists of two programs:
 - A listener (`cmd/listener/listener.go`), which binds on a port and
   prints out UDP packets it receives
 - A sender (`cmd/sender/sender.go`), which sends UDP packets to a
   given IP:port

Take a look at the source code for each file for an example of how they use UDP functions.  

To build both programs, run `make`.  See each program for its command-line options.  

## Important note

These programs just send and receive ASCII strings and thus print the
output using `Printf`.  In Snowcast, you are sending and receiving
**binary data** from arbitrary files, like mp3 files. Therefore, you
**should not use `Printf` to write data to stdout in Snowcast**.
Instead, use `Write`, which can write arbitrary binary data, eg:
```
os.Stdout.Write(...)
```
