#include <unistd.h>
#include <stdio.h>
#include <execinfo.h>

#define BACKTRACE_BUFFER 100

void trace() {
    int nptrs;
    void *buffer[BACKTRACE_BUFFER];
    nptrs = backtrace(buffer, BACKTRACE_BUFFER);
    
    char **strings;
    strings = backtrace_symbols(buffer, nptrs);

    for (size_t j = 0; j < nptrs; j++) {
        printf("%s\n", strings[j]);
    }
} 

int main() {
    printf("hello from parent (pid:%d)\n", (int) getpid());
    int rc = fork();
    if (rc < 0) {
        fprintf(stderr, "fork failed\n");
    } else if (rc == 0) {
        printf("\n======== child log trace ========\n");
        printf("child execute: (pid:%d)\n", (int) getpid());
        trace();
    } else {
        printf("\n======== parent log trace ========\n");
        printf("parent continues in (main), parent of %d (pid: %d)\n", rc, (int) getpid());
        trace();
    }
    return 0;
}

