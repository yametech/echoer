package fss

import (
	"fmt"
	"testing"
)

const ci = `
action ci
	addr = "compass.ym/tekton";
	method = http;
	args = (str project,str version,int retry_count);
	return = (SUCCESS | FAIL);
action_end
`

const approval = `
action approval
	addr = "nz.compass.ym/approval";
	method = http;
	args = (str work_order,int version);
	return = (AGREE | REJECT | NEXT | FAIL);
action_end
`

const notify = `
action notify
	addr = "nz.compass.ym/approval2";
	method = http;
	args = (str project, int version);
	return = (AGREE | REJECT | FAIL);
action_end
`

const deploy_1 = `
action deploy_1
	addr = "compass.ym/deploy";
	method = http;
	args = (str project, int version);
	return = (SUCCESS | FAIL);
action_end
`

const approval_2 = `
action approval_2
	addr = "nz.compass.ym/approval2";
	method = http;
	args = (str project, int version);
	return = (AGREE | REJECT | FAIL);
action_end
`

var actions = map[string]string{
	"ci":         ci,
	"approval":   approval,
	"notify":     notify,
	"deploy_1":   deploy_1,
	"approval_2": approval_2,
}

const my_flow = `
flow my_flow
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
`

const my_flow_run = `
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
`

const ci_https = `
action ci_https
	addr = "compass.ym/tekton";
	method = https;
	secret = (capem="xxadsa");
	args = (str project,str version,int retry_count);
	return = (SUCCESS | FAIL);
action_end
`

func Test_example_https(t *testing.T) {
	val := parse(NewFssLexer([]byte(ci_https)))
	if val != 0 {
		fmt.Println("syntax error")
	}
	ac, err := actionSymPoolGet("ci_https")
	if err != nil {
		t.Fatal(err)
	}
	_ = ac
}

func Test_example(t *testing.T) {
	// flex + goyacc
	// flow
	fmt.Println("--------flow")
	val := parse(NewFssLexer([]byte(my_flow)))
	if val != 0 {
		t.Fatal("syntax error")
	}

	fs, err := flowSymPoolGet("my_flow")
	if err != nil {
		t.Fatal(err)
	}

	for _, step := range fs.Steps {
		_ = step
	}

	fmt.Println("--------flow run")
	// flow run
	val = parse(NewFssLexer([]byte(my_flow_run)))
	if val != 0 {
		t.Fatal("syntax error")
	}

	fr, err := flowRunSymPoolGet("my_flow_run")
	if err != nil {
		t.Fatal(err)
	}
	for _, step := range fr.Steps {
		_ = step
	}

	// action
	for key, value := range actions {
		val = parse(NewFssLexer([]byte(value)))
		if val != 0 {
			fmt.Println("syntax error")
		}
		ac, err := actionSymPoolGet(key)
		if err != nil {
			t.Fatal(err)
		}
		_ = ac
	}
}

func TestNewFlowRunFSLParser(t *testing.T) {
	fss, err := NewFlowRunFSLParser().Parse(my_flow_run)
	if err != nil {
		t.Fatal(err)
	}
	_ = fss
}

func Benchmark_NewFlowRunFSLParser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fss, err := NewFlowRunFSLParser().Parse(my_flow_run)
		if err != nil {
			b.Fatal(err)
		}
		_ = fss
	}
}

func Benchmark_action_example(b *testing.B) {
	// flex + goyacc
	for i := 0; i < b.N; i++ {
		if parse(NewFssLexer([]byte(ci))) != 0 {
			b.Fatal("unknown failed")
		}
	}
}

func Benchmark_flow_example(b *testing.B) {
	// flex + goyacc
	for i := 0; i < b.N; i++ {
		if parse(NewFssLexer([]byte(my_flow))) != 0 {
			b.Fatal("unknown failed")
		}
	}
}

func Benchmark_flow_run_example(b *testing.B) {
	// flex + goyacc
	for i := 0; i < b.N; i++ {
		if parse(NewFssLexer([]byte(my_flow_run))) != 0 {
			b.Fatal("unknown failed")
		}
	}
}
