#include <stdlib.h>
#include <stdio.h>
#include <unistd.h>
#include <sys/types.h>

int main() {
    pid_t pid;
    pid = fork();
    switch (pid)
    {
    case -1:
        perror("fork fail");
        exit(EXIT_FAILURE);
        break;
    case 0:
        execlp("./child", NULL);
        perror("execlp"); // exec returns only on error
        exit(EXIT_FAILURE);
    default:
        puts("parent executed and exited");
        break;
        _exit(EXIT_SUCCESS);
    }
}