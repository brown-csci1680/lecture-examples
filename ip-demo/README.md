# IP-in-UDP demo

This example sends an IP packet inside a UDP packet as an example of
using Golang's IPv4 header library.  It also computes the IP checksum
using an external library.  

The most relevant files are as follows:
 - `cmd/udp-ip-send/main.go`:  Send an IP packet inside a UDP packet
   and compute the checksum
 - `cmd/udp-ip-recv/main.go`:  Receive an IP packet inside a UDP
   packet (with no other validation or checking)
   
Please see the comments inside each file for details.  More
information about this example will be posted in the next 24 hours.
Thanks for your patience!  

If you want to see a similar version of this example in C, see here:
https://github.com/brown-csci1680/lecture-examples/tree/main/ip-checksum
