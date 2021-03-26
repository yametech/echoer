/* Copyright 2020 laik.lj@me.com. All rights reserved. */
/* Use of this source code is governed by a Apache */
/* license that can be found in the LICENSE file. */

/* simplest version of calculator */

%{
package fss

import "strings"

var (
  print      = __yyfmt__.Print
  printf     = __yyfmt__.Printf
)
%}

// fields inside this union end up as the fields in a structure known
// as fssSymType, of which a reference is passed to the lexer.
%union {
  // flow
  Flow string
  Steps []Step
  _step Step
  _return Return
  _returns Returns
  _action Action
  _string string
  _identifier string
  _list []interface{}
  _args []Param
  _param Param
  _dict map[string]interface{}
  _variable string
  _number int64
  _secret map[string]string

  // action
  ActionStatement
  _addr []string
  _capem string
  _type ActionMethodType
  _grpc ActionMethodType
  _http ActionMethodType
  _https ActionMethodType
}

// any non-terminal which returns a value needs a type, which is
// really a field name in the above union struct

//%type <_action> action_stmt
//%type <steps> step_expr

%token ILLEGAL EOL
%token IDENTIFIER NUMBER_VALUE ID STRING_VALUE
%token LIST DICT
%token FLOW FLOW_END STEP ACTION ARGS DECI ACTION_END ADDR METHOD FLOW_RUN FLOW_RUN_END RETURN HTTPS SECRET
%token CAPEM
%token LPAREN RPAREN LSQUARE RSQUARE LCURLY RCURLY SEMICOLON COMMA COLON
%token HTTP GRPC
%token INT STR
%token ASSIGN OR AND
%token TO DEST

%start program

%%
program: flow_run_stmt
	{
		flowRunSymPoolPut($1.Flow,$1);
	}
	| action_stmt
	{
		actionSymPoolPut($1.ActionStatement.Name,$1);
	}
	| flow_stmt
	{
		flowSymPoolPut($1.Flow,$1);
	}
	;

action_stmt:
	ACTION STRING_VALUE action_content_addr_stmt action_content_method_stmt action_content_args_stmt action_return_stmt ACTION_END
	{
        $$.ActionStatement = ActionStatement{
            Name: $2._string,
            Addr: $3._addr,
            Type: $4._type,
            Args: $5._args,
            Returns: $6._returns,
        }
	}
	|
	ACTION IDENTIFIER action_content_addr_stmt action_content_method_stmt action_content_args_stmt action_return_stmt ACTION_END
    {
        $$.ActionStatement = ActionStatement{
            Name: $2._identifier,
            Addr: $3._addr,
            Type: $4._type,
            Args: $5._args,
            Returns: $6._returns,
        }
    }
    	|
	ACTION STRING_VALUE action_content_addr_stmt action_content_method_stmt action_content_secret_stmt action_content_args_stmt action_return_stmt ACTION_END
	{
	$$.ActionStatement = ActionStatement{
	    Name: $2._string,
	    Addr: $3._addr,
	    Type: $4._type,
	    Secret: $5._secret,
	    Args: $6._args,
	    Returns: $7._returns,
	}
	}
	|
	ACTION IDENTIFIER action_content_addr_stmt action_content_method_stmt action_content_secret_stmt action_content_args_stmt action_return_stmt ACTION_END
    {
	$$.ActionStatement = ActionStatement{
	    Name: $2._identifier,
	    Addr: $3._addr,
	    Type: $4._type,
	    Secret: $5._secret,
	    Args: $6._args,
	    Returns: $7._returns,
	}
    }
	;

action_content_addr_stmt:
	|ADDR ASSIGN STRING_VALUE SEMICOLON { $$._addr = append($$._addr,strings.Split(strings.Trim($3._string,"\""),",")...); }
	;

action_content_method_stmt:
	|METHOD ASSIGN HTTP SEMICOLON { $$._type = $3._http }
	|METHOD ASSIGN GRPC SEMICOLON { $$._type = $3._grpc }
	|METHOD ASSIGN HTTPS SEMICOLON { $$._type = $3._https }
	;

action_content_secret_stmt:
	SECRET ASSIGN action_content_secret_args_stmt SEMICOLON
	{
		$$._secret=make(map[string]string);
		$$._secret["capem"] = $3._capem;
	}

action_content_secret_args_stmt:
	|LPAREN CAPEM ASSIGN STRING_VALUE RPAREN
	{
		$$._capem = $4._string;
	}
	|LPAREN CAPEM ASSIGN IDENTIFIER RPAREN
	{
		$$._capem = $4._identifier;
	}
	;


action_content_args_stmt:
	|ARGS ASSIGN action_args_stmt SEMICOLON { $$._args = $3._args; }
	;


action_args_stmt:
	LPAREN action_args_content_stmts RPAREN { $$ = $2 }
	;

action_args_content_stmts:
	| action_args_content_stmts action_args_content_stmt { $$._args = append($$._args,$2._param); }
	;

action_args_content_stmt:
	COMMA action_args_content_stmt { $$._param = $2._param }
	|INT STRING_VALUE { $$._param = Param{ Name:$2._string, ParamType: NumberType}; }
	|STR STRING_VALUE { $$._param = Param{ Name:$2._string, ParamType: StringType}; }
	|INT IDENTIFIER { $$._param = Param{ Name:$2._identifier, ParamType: NumberType}; }
	|STR IDENTIFIER { $$._param = Param{ Name:$2._identifier, ParamType: StringType}; }
	;

action_return_stmt:
	RETURN ASSIGN return_stmt SEMICOLON
	{
		$$._returns = $3._returns ;
	}
	;

flow_stmt:
	FLOW STRING_VALUE flow_step_stmts FLOW_END
    {
        $$.Flow = $2._string;
        $$.Steps = $3.Steps;
    }
    |FLOW IDENTIFIER flow_step_stmts FLOW_END
    {
        $$.Flow = $2._identifier;
        $$.Steps = $3.Steps;
    }
    ;
flow_run_stmt:
	FLOW_RUN STRING_VALUE flow_step_stmts FLOW_RUN_END
	{
		$$.Flow = $2._string;
		$$.Steps = $3.Steps;
	}
	|FLOW_RUN IDENTIFIER flow_step_stmts FLOW_RUN_END
	{
        $$.Flow = $2._identifier;
        $$.Steps = $3.Steps;
    }
	;

flow_step_stmts:
	|flow_step_stmts flow_step_stmt
	{
		$$.Steps = append($$.Steps,$2._step);
	}
	;

flow_step_stmt:
	STEP IDENTIFIER TO return_stmt flow_action_stmt SEMICOLON
	{
		$$._step = Step{ Name:$2._identifier, Action:$5._action, Returns:$4._returns, StepType: Normal }
	}
	|DECI IDENTIFIER TO return_stmt flow_action_stmt SEMICOLON
	{
		$$._step = Step{ Name:$2._identifier, Action:$5._action, Returns:$4._returns, StepType: Decision }
	}
	|STEP STRING_VALUE TO return_stmt flow_action_stmt SEMICOLON
    {
        $$._step = Step{ Name:$2._string, Action:$5._action, Returns:$4._returns, StepType: Normal }
    }
    |DECI STRING_VALUE TO return_stmt flow_action_stmt SEMICOLON
    {
        $$._step = Step{ Name:$2._string, Action:$5._action, Returns:$4._returns, StepType: Decision }
    }
	;

flow_action_stmt:
	LCURLY flow_action_content_stmt RCURLY { $$ = $2; }
	;

flow_action_content_stmt:
	ACTION ASSIGN STRING_VALUE SEMICOLON ARGS ASSIGN flow_args_stmt SEMICOLON
	{
		$$._action = Action{ Name:$3._string, Args:$7._args };
	}
	|
	ACTION ASSIGN IDENTIFIER SEMICOLON ARGS ASSIGN flow_args_stmt SEMICOLON
    {
        $$._action = Action{ Name:$3._identifier, Args:$7._args };
    }
    |
	ACTION ASSIGN STRING_VALUE SEMICOLON
    {
        $$._action = Action{ Name:$3._string };
    }
    |
    ACTION ASSIGN IDENTIFIER SEMICOLON
    {
        $$._action = Action{ Name:$3._identifier };
    }
	;

flow_args_stmt:
	LPAREN flow_args_content_stmts RPAREN { $$ = $2 }
	;

flow_args_content_stmts:
	| flow_args_content_stmts flow_args_content_stmt
	{
		$$._args = append($$._args,$2._param);
	}
	;

flow_args_content_stmt:
	COMMA flow_args_content_stmt { $$._param = $2._param }
	|IDENTIFIER ASSIGN NUMBER_VALUE { $$._param = Param{ Name:$1._identifier, ParamType: NumberType, Value: $3._number }; }
	|IDENTIFIER ASSIGN STRING_VALUE { $$._param = Param{ Name:$1._identifier, ParamType: StringType, Value: $3._string }; }
	;

return_stmt:
	LPAREN return_content_stmts RPAREN { $$ = $2; }
	|LPAREN RPAREN { $$._returns = append($$._returns,Return{State:"DONE",Next:""}); }
	;

return_content_stmts:
	return_content_stmt
	{
		$$._returns = append($$._returns, $1._return);
	}
	|return_content_stmts OR return_content_stmt
	{
		$$._returns = append($$._returns, $3._return);
	}
	;

return_content_stmt:
	IDENTIFIER DEST IDENTIFIER {
		$$._return = Return{ State:$1._identifier, Next:$3._identifier };
	}
	|IDENTIFIER
	{
		$$._return = Return{ State:$1._identifier };
    }
	;

%%