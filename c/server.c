#include <sys/socket.h>
#include <unistd.h>
#include <stdio.h>
#include <netinet/in.h>
#include <arpa/inet.h>


#include "log.h"

#define PORT 8080
#define BUFFER_SIZE 1024

// main driver function
// 1. create a socket
// 2. we can configure socket, but do nothing here first
// 3. close socket 
int main() {
    int sockfd;
    sockfd = socket(AF_INET, SOCK_STREAM, 6);
    if (sockfd == -1) {
        perror("Error opening socket");
        return 1;
    }
    log_with_time("openned socket with fd:%d", sockfd);

    struct sockaddr_in addr;
    addr.sin_family = AF_INET;
    addr.sin_addr.s_addr = inet_addr("127.0.0.1");
    addr.sin_port = htons(PORT); 
    if (bind(sockfd, (struct sockaddr *) &addr, sizeof(addr)) < 0) {
        perror("Failed to bind socket");
        return 1;
    };

    int listenBacklog = 1;
    if (listen(sockfd, listenBacklog) < 0) {
        perror("Failed to listen on port");
        return 1;
    }

    socklen_t client_addr_size;
    struct sockaddr_in client_addr;
    client_addr_size = sizeof(client_addr);
    log_with_time("waiting for connection at port:%d", PORT);
    for (;;) {
        int clientfd = accept(sockfd, (struct sockaddr *) &client_addr, &client_addr_size);
        if (clientfd == -1) {
            perror("Failed to accept client connection");
            return 1;
        }

        char client_addr_str[INET_ADDRSTRLEN];
        inet_ntop(AF_INET, &client_addr.sin_addr, client_addr_str, INET_ADDRSTRLEN);
        log_with_time("received connection, new socket:%d, client address:%s", clientfd, client_addr_str);
        
        char buffer[BUFFER_SIZE];
        ssize_t data_len = read(clientfd, buffer, BUFFER_SIZE);
        log_with_time("received data (%d):[%sb], sending response to client", data_len, buffer);

        char write_back[1024];
        sprintf(write_back, "received data (%zdb)", data_len);
        if (write(clientfd, write_back, strlen(write_back)) < 0) {
            log_with_time("failed to respond to client");
        }

        close(clientfd);
    }
    
    log_with_time("closed socket with fd:%d", sockfd);
    close(sockfd);

    return 0;
}
