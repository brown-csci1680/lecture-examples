/*
 * ipsum_example.c - Send a virtual IP packet
 * This example builds a virtual IP packet and computes
 * the checksum.
 *
 * This example is only designed to show you how to fill in the IP
 * header, compute the checksum, and send the packet: most packet
 * fields are hard-coded, and the interface for sending packets looks
 * nothing like what you should use in the project.
 *
 * Thus, I do NOT recommend copying this code directly, with one
 * exception: you're welcome to just copy ip_sum().  Otherwise, you
 * should just use this example to show you how you can build packets
 * and compute the checksum in your own work.
 * 
 * You can run this program as follows:
 * ./ipsum_example localhost 5000 localhost 5001
 * This will send a virtual IP packet from UDP port 5000 to
 * UDP port 5001.  All fields inside the packet are hard-coded.
 * The payload is the string "Hello world!", similar to how
 * you should be formatting packets with the "send" command.
 * 
 * You can inspect the sent packet in wireshark.  To see that the
 * checksum is correct in Wireshark, right click on the IP header in
 * the middle pane, select "Protocol Preferences" and then "Validate
 * IP checksum".
 */
#include <assert.h>
#include <netdb.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <netinet/ip.h>
#include <sys/types.h>
#include <sys/socket.h>

#define MAX_PACKET_BYTES 1400

uint16_t ip_sum(void *packet, int n);
int create_bind_sock(char *addr, char *port);
void get_addr_by_name(char *addr, char *port,
		      struct sockaddr *sockaddr, socklen_t *socklen);
int send_fake_virtual_ip(int sock, struct sockaddr *udp_addr, socklen_t udp_addr_len,
			 char *payload, size_t payload_len);

// Compute the IP checksum
// This is a modified version of the example in RFC 1071
// https://datatracker.ietf.org/doc/html/rfc1071#section-4.1
uint16_t ip_sum(void *buffer, int len) {
  uint8_t *p = (uint8_t *)buffer;
  uint16_t answer;
  long sum = 0;
  uint16_t odd_byte = 0;

  while (len > 1) {
    uint16_t c = 0;
    c = (p[1] << 8) | p[0];
    
    sum += c;
    p += 2;
    len -= 2;
  }

  if (len == 1) {
    *(uint8_t*)&odd_byte = *p;
    sum += odd_byte;
  }

  sum = (sum >> 16) + (sum & 0xffff);
  sum += (sum >> 16);
  answer = ~sum;
  return answer;

}

// Get an IP address as a 32-bit number from string
uint32_t get_ip_addr(char *ip_str)
{
  struct in_addr addr;
  memset(&addr, 0, sizeof(struct in_addr));

  int rv;
  if ((rv = inet_pton(AF_INET, ip_str, &addr)) < 0) {
    perror("inet_pton");
  }

  // Now, we can get the address in 32-bit form from the sockaddr struct
  // For details, see:  https://beej.us/guide/bgnet/html/#structsockaddrman
  uint32_t addr_as_int = addr.s_addr;

  return addr_as_int;
}

// This function sends a fake virtual IP packet containing payload to
// the UDP address udp_addr using socket sock NOTE: This is just an
// example function which hard-codes all of the values in the virtual
// IP header.  In the project, you should have a function like
// send_ip(...), but it will look VERY different from this!!!  For
// example, you shouldn't be passing in the UDP socket and addr as
// arguments.  Instead, you may want to specify the virtual
// source/dest address, or an interface name
int send_fake_virtual_ip(int sock,
			 struct sockaddr *udp_addr, socklen_t udp_addr_len,
			 char *payload, size_t payload_len)
{
  // The total size of the packet can't be larger than the MTU
  assert(sizeof(struct iphdr) + payload_len < MAX_PACKET_BYTES);

  int bytes_sent;
  char buffer[MAX_PACKET_BYTES];
  
  struct iphdr header;
  int header_size = sizeof(struct iphdr);
  memset(&header, 0, header_size);
  
  header.version = 4;
  header.ihl = 5;      // Header length is always 20 bytes since no IP options
  header.ttl = 16;

  // Protocol hard-coded to 0 here
  // You will change this based on the message type
  header.protocol = 0;

  // Virtual IP source and destination are hard-coded here
  // You will want to do something different!
  header.saddr = get_ip_addr("192.168.0.2");
  header.daddr = get_ip_addr("192.168.0.1");

  uint16_t packet_len = header_size + payload_len;
  header.tot_len = htons(packet_len);
  
  uint16_t checksum = ip_sum(&header, header_size);
  header.check = checksum;

  // Copy the header and payload into one buffer
  memcpy(buffer, &header, header_size);
  memcpy(buffer + header_size, payload, payload_len);

  // Send this over our UDP socket
  if ((bytes_sent = sendto(sock, buffer, header_size + payload_len, 0,
			   udp_addr, udp_addr_len)) < 0) {
    perror("sendto");
    exit(1);
  }

  return bytes_sent;
}

int main(int argc, char **argv)
{
    int sock;

    if (argc != 5) {
	printf("%s:  sender <udp_src_addr> <udp_src_port> <udp_dst_addr> <udp_dst_port>\n",
	       argv[0]);
	exit(1);
    }
    // Parse server port from command line
    char *udp_src_addr = argv[1];
    char *udp_src_port = argv[2];
    char *udp_dst_addr = argv[3];
    char *udp_dst_port = argv[4];

    // Create a socket bound on this UDP port
    sock = create_bind_sock(udp_src_addr, udp_src_port);

    // Resolve the destination address from command line
    struct sockaddr dest_addr;
    socklen_t dest_len = sizeof(struct sockaddr);
    get_addr_by_name(udp_dst_addr, udp_dst_port,
		     &dest_addr, &dest_len);

    // Message we will send
    char *message = "Hello world!";
    size_t message_length = strlen(message);
    
    int bytes_sent = send_fake_virtual_ip(sock, &dest_addr, dest_len,
					  message, message_length);
    printf("Sent %d bytes\n", bytes_sent);

    close(sock);
    return 0;
}

// Create and bind a UDP socket on addr:port
int create_bind_sock(char *addr, char *port)
{
  int rv, sock;
  
  struct addrinfo hints, *res, *servinfo;
  memset(&hints, 0, sizeof (struct addrinfo));
  hints.ai_family = AF_INET;
  hints.ai_socktype = SOCK_DGRAM;
  hints.ai_flags = AI_PASSIVE;

  if ((rv = getaddrinfo(addr, port, &hints, &servinfo)) != 0) {
    perror("getaddrinfo");
    exit(1);
  }

  for (res = servinfo; res != NULL; res = res->ai_next) {
    if((sock = socket(res->ai_family, res->ai_socktype, res->ai_protocol)) < 0) {
      continue;
    }

    // Bind to the address and port
    if((rv = bind(sock, res->ai_addr, res->ai_addrlen)) < 0) {
      close(sock);
      continue;
    }

    break;
  }

  if (res == NULL) {
    printf("Could not bind to socket\n");
    exit(1);
  }

  return sock;
}  

// Get a sockaddr suitable for UDP's sendto() from addr:port
// When called, socklen must point to a value that contains the
// size of sockaddr.
void get_addr_by_name(char *addr, char *port,
		      struct sockaddr *sockaddr, socklen_t *socklen)
{
  int rv;
  struct addrinfo hints, *servinfo;
  memset(&hints, 0, sizeof (struct addrinfo));
  hints.ai_family = AF_INET;
  hints.ai_socktype = SOCK_DGRAM;
  hints.ai_flags = AI_PASSIVE;

  if ((rv = getaddrinfo(addr, port, &hints, &servinfo)) != 0) {
    perror("getaddrinfo");
    exit(1);
  }

  // NOTE: getaddrinfo actually returns a linked list of results,
  // which should each be tested to find the correct one.  For the
  // purposes of this example (and, likey, the project), this is
  // sufficient, but in practice better erroror handling shold happen
  // here!
  if (servinfo == NULL) {
    printf("Could not determine address\n");
    exit(1);
  }
  
  // Provided socklen must be large enough to store the result
  assert(*socklen >= servinfo->ai_addrlen);

  // Copy the frist result into sockaddr and socklen
  memcpy(sockaddr, servinfo->ai_addr, servinfo->ai_addrlen);
  *socklen = servinfo->ai_addrlen;
}

