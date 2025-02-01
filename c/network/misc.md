# Misc Development Log
## C Header Files
C allows us to separate function signatures and variable from their implementations.

We can refer to these headers and reuse the same function / variables without duplicating them.

```c
// in log.h
void log_with_time(char *format, ...) {
    ..
}
```
