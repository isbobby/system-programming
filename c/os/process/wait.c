#include <unistd.h>
#include <stdio.h>
#include <sys/wait.h>

#define BACKTRACE_BUFFER 100

int main() {
    printf("hello from parent (pid:%d)\n", (int) getpid());
    int rc = fork();
    char * placholder;
    if (rc < 0) {
        fprintf(stderr, "fork failed\n");
    } else if (rc == 0) {
        printf("child execute: (pid:%d)\n", (int) getpid());
        printf("enter any key to continue child process, run ps now to see active process\n");
        scanf(placholder, "%d");
    } else {
        int rc_wait = wait(NULL);
        printf("child process waited for\n");
        printf("enter any key to continue, run ps now to see active process\n");
        scanf(placholder, "%d");
        printf("parent continues in (main), parent of %d (pid: %d)\n", rc, (int) getpid());
        return 0;
    }
}
