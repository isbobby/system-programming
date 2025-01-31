#include <time.h>
#include <stdarg.h>
#include <string.h>
#include <stdio.h>

void log_with_time(char *format, ...) {
    int time_stamp_size = 23;
    int new_line_size = 2;

    char time_buff[time_stamp_size];
    time_t rawtime = time(0);
    strftime(time_buff, time_stamp_size, "[%Y-%m-%d %H:%M:%S] ", localtime(&rawtime));

    char log_buff[time_stamp_size + strlen(format) + new_line_size];
    strcpy(log_buff, time_buff);
    strcat(log_buff, format);
    strcat(log_buff, "\n");

    va_list args;
    va_start(args, format);
    vprintf(log_buff, args);
    va_end(args);
}