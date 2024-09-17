# Guessing game sockets demo

This directory contains source for the "guessing game" demo from
lecture.  

This version is a completed version of the "guessing game protocol
that can play the game simply.  However, it is missing several key
features like error handling, support for reading messages that arrive
in multiple `Read` calls, and more--please see the comments for more
details.  For a version that includes better error and message
handling, please see the "Full sockets example".  

## How to run

You can use run the demo in your container repository, or anywhere Go
is installed.  To build the code, run `make`.  

To run the example:
 - In one terminal, run the server:  `./server`
 - In one or more terminals, run the client:  `./client`
 
To make the server easy to run in class, the server is hard-coded to
listen on port 6666.  
