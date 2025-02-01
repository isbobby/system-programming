#include <execinfo.h>
#include <stdlib.h>
#include <stdio.h>

#define BACKTRACE_BUFFER 100

void odd(int *n);
void even(int * n);

void even(int * n) {
    *n -= 1;
    odd(n);
}

void odd(int *n) {
    *n -= 1;
    if (*n < 0) {
        int nptrs;
        void *buffer[BACKTRACE_BUFFER];
        nptrs = backtrace(buffer, BACKTRACE_BUFFER);
        
        char **strings;
        strings = backtrace_symbols(buffer, nptrs);

        for (size_t j = 0; j < nptrs; j++) {
            printf("%s\n", strings[j]);
        }

        return;
    }

    if (*n % 2 == 1) {
        odd(n);
    } else {
        even(n);
    }
};

int main() {
    int n = 5;
    odd(&n);
}
