#if defined(_WIN32) && defined(_WINDOWS)
#define inline 
#define snprintf _snprintf
int strncasecmp(char s1, char s2, register int n);
#endif
