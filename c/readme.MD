# Development Journal
The development goal here is to create a TCP server in C using standard libraries.

Referring to my notes, a basic TCP server should do the following

1. create a TCP socket
2. optinally configure the socket 
3. transform the socket into a passive socket 
4. accept incoming connection and handle requests
5. close connection and clean up resources

## Creating a TCP Socket
We need to refer to [man `socket(2)`](https://man7.org/linux/man-pages/man2/socket.2.html), where we obtain the following key information
```c
#include <sys/socket.h>
int socket(int domain, int type, int protocol);
```

Following the doc, we will use `AF_INET` as the `domain` argument, and `SOCK_STREAM` as the `type` argument. Digging down further, we see this quote

>  The protocol specifies a particular protocol to be used with the socket. Normally only a single protocol exists to support a particular socket type within a given protocol family, in which case protocol can be specified as 0

Since `SOCK_STREAM` is implemented on TCP only, we can leave it as 0. `man socket(2)` also provides a link to `man protocols(5)`, which lists the protocol in `etc/protocols`.

If we do `cat /etc/protocols | grep TCP`, we get the following (on my Macbook).
```
tcp	6	TCP		# transmission control protocol
```

With the above information, we can write the following code (note, the `close()` is added here but will be discussed later). We can use `0` as protocol, but `6 (tcp)` will be functionally the same.
```c
#include <sys/socket.h>
#include <unistd.h>

int main() {
    int sockfd;
    sockfd = socket(AF_INET, SOCK_STREAM, 0);
    close(sockfd);
}
```