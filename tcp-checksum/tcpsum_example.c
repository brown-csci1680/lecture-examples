/*
 * tcpsum_example.c - Send a virtual TCP packet
 * This example builds a TCP packet encapsulated in a virtual
 * IP packet
 *
 * This example is only designed to show you how to fill in the IP and
 * TCP headers, compute the checksum, and send the packet: most packet
 * fields are hard-coded, and the interface for sending packets looks
 * nothing like what you should use in the project.
 *
 * Thus, I do NOT recommend copying this code directly, with one
 * exception: you're welcome to just copy ip_sum() and
 * compute_tcp_checksum().  Otherwise, you should just use this
 * example to show you how you can build packets and compute the
 * checksum in your own work.  Note: compute_tcp_checksum() is pretty
 * naive--you can modify it to avoid extra copying if you want better
 * performance.
 *
 * You can run this program as follows:
 * ./tcpsum_example <message>
 * This will send a virtual TCP packet from UDP port 5000 to
 * UDP port 5001.  All fields inside the packet are hard-coded.
 * The packet is an ACK starting at sequence number 1, with
 * the payload given as the first agument.  For example, running:
 * ./tcpsum_example "Hello world!"
 * will send a packet with "Hello world!" as the payload.
 *
 * You can inspect the sent packet in wireshark.  To see that the
 * checksum is correct in Wireshark, right click on the TCP header in
 * the middle pane, select "Protocol Preferences" and then "Validate
 * TCP checksum".
 */
#include <assert.h>
#include <netdb.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <netinet/ip.h>
#include <netinet/tcp.h>
#include <sys/types.h>
#include <sys/socket.h>

#define MAX_PACKET_BYTES 1400

//uint16_t ip_sum(void *packet, int n);
int create_bind_sock(char *addr, char *port);
void get_addr_by_name(char *addr, char *port,
		      struct sockaddr *sockaddr, socklen_t *socklen);
int send_fake_virtual_ip(int sock, struct sockaddr *udp_addr, socklen_t udp_addr_len,
			 uint32_t virtual_ip_src, uint32_t virtual_ip_dst,
			 uint8_t protocol,
			 char *payload, size_t payload_len);
int send_fake_virtual_tcp(int sock,
			  struct sockaddr *udp_addr, socklen_t udp_addr_len,
			  uint32_t virtual_ip_src, uint32_t virtual_ip_dst,
			  char *payload, size_t payload_len);

// Compute the IP checksum (TCP uses the same checksum)
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
			 uint32_t virtual_ip_src, uint32_t virtual_ip_dst,
			 uint8_t protocol,
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
  header.protocol = protocol;

  // Virtual IP source and destination are hard-coded here
  // You will want to do something different!
  header.saddr = virtual_ip_src;
  header.daddr = virtual_ip_dst;

  uint16_t packet_len = header_size + payload_len;
  header.tot_len = htons(packet_len);

  uint16_t checksum = ip_sum((char *)&header, header_size);
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

// The TCP checksum is computed based on a "pesudo-header" that
// combines the (virtual) IP source and destination address, protocol value,
// as well as the TCP header and payload
//
// This is one example one way to combine all of this information
// and compute the checksum.  This is not a particularly fast way,
// as it involves copying all of the data into one buffer.
// Yet, it works, and you can use it.
//
// For more details, see the "Checksum" component of RFC793 Section 3.1,
// https://www.ietf.org/rfc/rfc793.txt (pages 14-15)
uint16_t compute_tcp_checksum(uint32_t virtual_ip_src, uint32_t virtual_ip_dst,
			      struct tcphdr *tcp_header, char *payload, size_t payload_len)
{

    struct pseudo_header {
	uint32_t ip_src;
	uint32_t ip_dst;
	uint8_t zero;
	uint8_t protocol;
	uint16_t tcp_length;
    };

    struct pseudo_header ph;

    size_t ph_len = sizeof(struct pseudo_header);
    size_t hdr_len = sizeof(struct tcphdr);
    assert(ph_len == 12);
    assert(hdr_len == 20);

    // Now fill in the pesudo header
    memset(&ph, 0, sizeof(struct pseudo_header));
    ph.ip_src = virtual_ip_src;
    ph.ip_dst = virtual_ip_dst;
    ph.protocol = 6;  // TCP's assigned IP protocol number is 6

    // From RFC: "The TCP Length is the TCP header length plus the
    // data length in octets (this is not an explicitly transmitted
    // quantity, but is computed), and it does not count the 12 octets
    // of the pseudo header."
    ph.tcp_length = htons(hdr_len + payload_len);

    size_t total_len = ph_len + hdr_len + payload_len;
    char buffer[total_len];
    memset(buffer, 0, total_len);
    memcpy(buffer, &ph, ph_len);
    memcpy(buffer + ph_len, tcp_header, hdr_len);
    memcpy(buffer + ph_len + hdr_len, payload, payload_len);

    uint16_t checksum = ip_sum(buffer, total_len);

    return checksum;
}

int send_fake_virtual_tcp(int sock,
			  struct sockaddr *udp_addr, socklen_t udp_addr_len,
			  uint32_t virtual_ip_src, uint32_t virtual_ip_dst,
			  char *payload, size_t payload_len)
{
    // The total size of the packet can't be larger than the MTU
    // plus the size of the IP and TCP headers
    assert((sizeof(struct iphdr) + sizeof(struct tcphdr) + payload_len)
	   < MAX_PACKET_BYTES);

    struct tcphdr header;
    size_t header_size = sizeof(struct tcphdr);

    memset(&header, 0, header_size);

    // Similar to IP header, header length is always 20 bytes because no TCP options
    header.th_off = 5;

    // Set some random source ports
    header.th_sport = htons(12345);
    header.th_dport = htons(80);

    // Hard-code the sequence number, ACK number, flags, and window size
    header.th_seq = htonl(1);
    header.th_ack = htonl(1);
    header.ack = 1;
    header.th_win = 65535;
    // Header checksum value should be zero when computing checksum
    header.th_sum = 0;

    uint16_t this_sum = compute_tcp_checksum(virtual_ip_src, virtual_ip_dst,
					     &header, payload, payload_len);
    header.th_sum = this_sum;

    // Now build and send the segment
    char buffer[MAX_PACKET_BYTES];

    // Copy the header and payload into one buffer
    memcpy(buffer, &header, header_size);
    memcpy(buffer + header_size, payload, payload_len);

    int bytes_sent = send_fake_virtual_ip(sock, udp_addr, udp_addr_len,
					  virtual_ip_src, virtual_ip_dst,
					  6, buffer, header_size + payload_len);

    return bytes_sent;
}

int main(int argc, char **argv)
{
    int sock;

    if (argc != 2) {
	printf("%s:  sender <message>\n", argv[0]);
	exit(1);
    }

    // Message we will send (from command line)
    char *message = argv[1];
    size_t message_length = strlen(message);

    // Hard-code UDP source and dest
    char *udp_src_addr = "127.0.0.1";
    char *udp_src_port = "5000";
    char *udp_dst_addr = "127.0.0.1";
    char *udp_dst_port = "5001";

    // Create a socket bound on this UDP port
    sock = create_bind_sock(udp_src_addr, udp_src_port);

    // Resolve the UDP destination address from command line
    struct sockaddr dest_addr;
    socklen_t dest_len = sizeof(struct sockaddr);
    get_addr_by_name(udp_dst_addr, udp_dst_port,
		     &dest_addr, &dest_len);

    // Pick some arbitrary virtual IP addresses for this test packet
    uint32_t virtual_ip_src = get_ip_addr("192.168.0.1");
    uint32_t virtual_ip_dst = get_ip_addr("192.168.0.2");


    int bytes_sent = send_fake_virtual_tcp(sock, &dest_addr, dest_len,
					   virtual_ip_src, virtual_ip_dst,
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
