# TCP-IP-in-UDP demo

This example sends an TCP packet inside a virtual IP packet UDP packet as an example of
using Golang's IPv4 header library and the netstack package to parse
TCP packets.  

NOTE:  This example DOES NOT compute the TCP checksum.  An update will
be posted about this later.

The most relevant files are as follows:
 - `cmd/tcp-ip-send/main.go`:  Send a TCP packet inside a virtual IP package
 - `cmd/tcp-ip-recv/main.go`:  Receive a TCP packet inside a virtual
   IP packet
   
Please see the comments inside each file for details.  More
information about this example will be posted in the next 24 hours.
Thanks for your patience!  

If you want to see a similar version of this example in C, see here:
https://github.com/brown-csci1680/lecture-examples/tree/main/tcp-checksum
