/* Copyright 2020 laik.lj@me.com. All rights reserved. */
/* Use of this source code is governed by a Apache */
/* license that can be found in the LICENSE file. */

enum
{
    ILLEGAL = 10000,
    EOL = 10001,

    IDENTIFIER = 260, // X
    NUMBER_VALUE = 261, // digit
    STRING_VALUE = 264, // "string"

    LIST = 262, // []
    DICT = 263, // {...}


    FLOW = 265, // FLOW $name
    FLOW_END = 266,	 // END
    DECI = 267, // DECI $node
    STEP = 268, // STEP $node
    ACTION = 269, // ACTION
    ARGS = 270, // ARGS=$?

    LPAREN = 271,	 // (
    RPAREN = 272,	 // )
    LSQUARE = 273,	 // [
    RSQUARE = 274,	 // ]
    LCURLY = 275,	 // {
    RCURLY = 276,	 // }
    ASSIGN = 277,	 // =
    SEMICOLON = 278, // ;
    OR = 279,	 // |
    AND = 280,	 // &
    TO = 281,	 // =>
    COMMA = 282,	 // ,
    COLON = 283,	 // :
    DEST = 284,	// ->

    ADDR = 285, // action ADDR keyword
    METHOD = 286, // action METHOD keyword
    ACTION_END = 288, // action ACTION_END keyword
    INT = 291,
    STR = 292,

    HTTP = 289,
    GRPC = 290,

    FLOW_RUN = 293,
    FLOW_RUN_END = 294,
    RETURN = 295,
    SECRET = 296,
    HTTPS = 297,
    CAPEM = 299,
};
