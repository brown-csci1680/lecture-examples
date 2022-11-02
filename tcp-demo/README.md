# TCP-IP-in-UDP demo

This example sends an TCP packet inside a virtual IP packet UDP packet
as an example of using Golang's IPv4 header library and the netstack
package to parse TCP packets.  

The most relevant files are as follows:
 - `cmd/tcp-ip-send/main.go`:  Send a TCP packet inside a virtual IP package
 - `cmd/tcp-ip-recv/main.go`:  Receive a TCP packet inside a virtual
   IP packet
 - `pkg/iptcp_utils/iptcp_utils.go`:  Utility functions, including
   checksum functions
 
**Please see the comments inside each file for details.  There are a
   lot of comments to help explain how things work, please read them!**  

If you want to see a similar version of this example in C, see here:
https://github.com/brown-csci1680/lecture-examples/tree/main/tcp-checksum
