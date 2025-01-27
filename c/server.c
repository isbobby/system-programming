#include <sys/socket.h>
#include <unistd.h>

// main driver function
// 1. create a socket
// 2. we can configure socket, but do nothing here first
// 3. close socket 
int main() {
    int sockfd;

    sockfd = socket(AF_INET, SOCK_STREAM, 6);

    close(sockfd);

    return 0;
}