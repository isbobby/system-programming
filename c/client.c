#define PORT 8080

#include <unistd.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <string.h>

#include "log.h"

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
    if (connect(sockfd, (struct sockaddr *) &addr, sizeof(addr)) < 0) {
        perror("Failed to connnect to server socket");
        return 1;
    };

    // hard coded data for now
    char * data = "(hard coded) Hi!";

    if (write(sockfd, data, strlen(data)) < 0) {
        perror("failed to write data");
        return 1;
    }
    
    log_with_time("wrote data (%db) to server", strlen(data));

    return 0;
}