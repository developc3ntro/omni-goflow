package main

// go install github.com/developc3ntro/omni-goflow/cmd/exptester; exptester "@(lower(contact.name))"

import (
	"fmt"
	"os"

	"github.com/developc3ntro/omni-goflow/envs"
	"github.com/developc3ntro/omni-goflow/test"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: exptester <expression>")
		os.Exit(1)
	}

	output, err := expTester(os.Args[1])
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(output)
	}
}

func expTester(template string) (string, error) {
	session, _, err := test.CreateTestSession("http://localhost:49995", envs.RedactionPolicyNone)
	if err != nil {
		return "", err
	}

	run := session.Runs()[0]

	return run.EvaluateTemplate(template)
}
