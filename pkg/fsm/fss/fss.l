/* Copyright 2020 laik.lj@me.com. All rights reserved. */
/* Use of this source code is governed by a Apache */
/* license that can be found in the LICENSE file. */

%option noyywrap
%option caseless

%{
#include "token.h"
%}


%%
"FLOW"                               { return FLOW; }
"FLOW_END"                           { return FLOW_END; }
"FLOW_RUN"                           { return FLOW_RUN; }
"FLOW_RUN_END"                       { return FLOW_RUN_END; }
"DECI"                               { return DECI; }
"STEP"                               { return STEP; }
"ACTION"                             { return ACTION; }
"RETURN"                             { return RETURN; }
"ACTION_END"                         { return ACTION_END; }
"ARGS"                               { return ARGS; }
"ADDR"                               { return ADDR; }
"METHOD"                             { return METHOD; }
"HTTP"                               { return HTTP; }
"HTTPS"                              { return HTTPS; }
"GRPC"                               { return GRPC; }
"INT"                                { return INT; }
"STR"                                { return STR; }
"SECRET"                             { return SECRET; }
"CAPEM"                              { return CAPEM; }

"->"   { return DEST; }
"("    { return LPAREN; }
")"    { return RPAREN; }
"["    { return LSQUARE; }
"]"    { return RSQUARE; }
"{"    { return LCURLY; }
"}"    { return RCURLY; }
"=>"   { return TO; }
"="    { return ASSIGN; }
";"    { return SEMICOLON; }
"|"    { return OR; }
"&"    { return AND; }
":"    { return COLON; }
","    { return COMMA; }

[0-9]+                                 { return NUMBER_VALUE; }
\[.*\]                                 { return LIST; }
list\((.*,.*)\)                        { return LIST; }
dict\((.*=.*)\)                        { return DICT; }
\"([^\"]*)\"                           { return STRING_VALUE; }
\`([^\`].*)\`                          { return STRING_VALUE; }
[_a-z_A-Z_][_a-z_A-Z_0-9_]*            { return IDENTIFIER; }

\n     { /* ignore whitespace */ }
[ \t]  { /* ignore whitespace */ }
.      {  return ILLEGAL; }

%%