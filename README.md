## echoer
自定义流程驱动系统自定义业务事件中间件 (WIP)

### echoer 设计与开发总结

1. FSM数学模型,资源定义为flow/flow_run/step/action/action_run
2. FSL(DSL)支持action定义,flow定义,flow_run定义

### api-server
1. 前端业务api

### flow controller
1. 实现可基于FSL语言/api自定义flow(流程)的控制
2. 自动状态转换驱动

### action controller
1. 延时队列实现协调通知驱动


### client
1. 流程与FSL
![graph ](./pkg/fsm/fss/example/fsm.jpg)

```
action ci
	addr = "compass.ym/tekton";
	method = http;
	args = (str project,str version,int retry_count);
	return = (SUCCESS | FAIL);
action_end

action approval
	addr = "nz.compass.ym/approval";
	method = http;
	args = (str work_order,int version);
	return = (AGREE | REJECT | NEXT | FAIL);
action_end

action deploy_1
	addr = "compass.ym/deploy";
	method = http;
	args = (str project, int version);
	return = (SUCCESS | FAIL);
action_end

action approval_2
	addr = "nz.compass.ym/approval2";
	method = http;
	args = (str project, int version);
	return = (AGREE | REJECT | FAIL);
action_end

action notify
	addr = "nz.compass.ym/approval2";
	method = http;
	args = (str project, int version);
	return = (AGREE | REJECT | FAIL);
action_end

flow_run my_flow_run
	step A => (SUCCESS->D | FAIL->A) {
		action = "ci";
		args = (project="https://github.com/yametech/compass.git",version="v0.1.0",retry_count=10);
	};
	deci D => ( AGREE -> B | REJECT -> C | NEXT -> E | FAIL -> D ) {
		action="approval";
		args=(work_order="nz00001",version=12343);
	};
	step B => (FAIL->B | SUCCESS->C) {
		action="deploy";
		args=(env="release",version=12343);
	};
	step E => (REJECT->C | AGREE->B | FAIL->E) {
		action="deploy";
		args=(env="test",version=12343);
	};
    step C => (FAIL->C){
        action="notify";
        args=(work_order="nz00001",version=12343);
    };
flow_run_end
```

2. 终端访问,可以基于FSL定义业务及流程
```
☁  cli [master] ⚡  go run main.go

echoer
        /\_/\                                                 ##         .
      =( °w° )=                                         ## ## ##        ==
        )   (     // 📒 🤔🤔🤔  ♻︎                       ## ## ## ## ##    ===
       (__ __)           === == ==                /""""""""""""""""\___/ ===
 /"""""""""""""" //\___/ === == ==                           ~~/~~~~ ~~~ ~~~~ ~~ ~ /  ===- ~~~
{                       /  == =-                  \______ o          _,/
 \______ O           _ _/                          \      \       _,'
  \    \         _ _/                               '--.._\..--''
    \____\_______/__/__/

>
>
> flow my_flow
	step A => (SUCCESS->D | FAIL->A) {
        action = "ci";
	};
	deci D => ( AGREE -> B | REJECT -> C | NEXT -> E | FAIL -> D ) {
		action="approval";
	};
	step B => (FAIL->B | SUCCESS->C) {
		action="deploy_1";
	};
	STEP E => (REJECT->C | AGREE->B | FAIL->E) {
		action="approval_2";
	};
    step C => (FAIL->C) {
        action="notify";
    };
flow_end
/
E| factory translation flow (flow my_flow	step A => (SUCCESS->D | FAIL->A) {        action = "ci";	};	deci D => ( AGREE -> B | REJECT -> C | NEXT -> E | FAIL -> D ) {		action="approval";	};	step B => (FAIL->B | SUCCESS->C) {		action="deploy_1";	};	STEP E => (REJECT->C | AGREE->B | FAIL->E) {		action="approval_2";	};    step C => (FAIL->C) {        action="notify";    };flow_end) error: (not without getting action (ci) definition)
>
>
> action ci
	addr = "compass.ym/tekton";
	method = http;
	args = (str project,str version,int retry_count);
	return = (SUCCESS | FAIL);
action_end
/
action approval
	addr = "nz.compass.ym/approval";
	method = http;
	args = (str work_order,int version);
	return = (AGREE | REJECT | NEXT | FAIL);
action_end
/
action deploy_1
	addr = "compass.ym/deploy";
	method = http;
	args = (str project, int version);
	return = (SUCCESS | FAIL);
action_end
/
action approval_2
	addr = "nz.compass.ym/approval2";
	method = http;
	args = (str project, int version);
	return = (AGREE | REJECT | FAIL);
action_end
/
action notify
	addr = "nz.compass.ym/approval2";
	method = http;
	args = (str project, int version);
	return = (AGREE | REJECT | FAIL);
action_end
/
OK
> OK
> OK
> OK
> OK
> flow my_flow
	step A => (SUCCESS->D | FAIL->A) {
        action = "ci";
	};
	deci D => ( AGREE -> B | REJECT -> C | NEXT -> E | FAIL -> D ) {
		action="approval";
	};
	step B => (FAIL->B | SUCCESS->C) {
		action="deploy_1";
	};
	STEP E => (REJECT->C | AGREE->B | FAIL->E) {
		action="approval_2";
	};
    step C => (FAIL->C) {
        action="notify";
    };
flow_end
> /
OK
> get flow my_flow
/
R| {"_id":"5f6b248766a5e1c8c78be432","metadata":{"kind":"flow","labels":null,"name":"my_flow","version":1600857223},"spec":{"steps":[{"action_name":"ci","returns":{"FAIL":"A","SUCCESS":"D"}},{"action_name":"approval","returns":{"AGREE":"B","FAIL":"D","NEXT":"E","REJECT":"C"}},{"action_name":"deploy_1","returns":{"FAIL":"B","SUCCESS":"C"}},{"action_name":"approval_2","returns":{"AGREE":"B","FAIL":"E","REJECT":"C"}},{"action_name":"notify","returns":{"FAIL":"C"}}]}}
>
```