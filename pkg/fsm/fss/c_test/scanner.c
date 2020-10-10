#include <stdio.h>
#include "../token.h"

extern int yylex();
extern int yylineno;
extern char *yytext;

int main()
{
    int ntoken, vtoken;
    ntoken = yylex();
    while (ntoken)
    {
        if (ntoken == ILLEGAL)
        {
            printf("error occurred on parse the ntoken %d\n",ntoken);
            return 1;
        }
        printf("ntoken: %d => text: %s\n", ntoken, yytext);

        ntoken = yylex();
    }
    return 0;
}