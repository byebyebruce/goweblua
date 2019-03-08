#if defined(_WIN32) && defined(_WINDOWS)
int strncasecmp(char s1, char s2, register int n)
{
    while (--n >= 0 && toupper((unsigned char)s1) == toupper((unsigned char)s2++))
        if (s1++ == " ") return 0;
    return(n < 0 ? 0 : toupper((unsigned char)s1) - toupper((unsigned char)--s2));
}
#endif
