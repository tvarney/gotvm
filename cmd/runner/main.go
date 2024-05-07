package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kingpin"
	"github.com/tvarney/gotvm/assembler"
	"github.com/tvarney/gotvm/reference"
)

func main() {
	showBytecode := false
	showStack := false
	trace := false
	filename := ""

	argparse := kingpin.New("runner", "Run an assembly program in the VM")
	argparse.Flag("show-bytecode", "Print the raw bytecode after assembly").BoolVar(&showBytecode)
	argparse.Flag("show-stack", "Print the values on the stack at the end of the program").BoolVar(&showStack)
	argparse.Flag("trace", "Print the values of the stack after each opcode").BoolVar(&trace)
	argparse.Arg("file", "The file to assemble and run").Required().StringVar(&filename)

	if _, err := argparse.Parse(os.Args[1:]); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	lines := strings.Split(string(content), "\n")
	bytecode := assembler.Assemble(lines, assembler.ReportPrint)
	if bytecode == nil {
		fmt.Printf("Error: no bytecode assembled")
		os.Exit(1)
	}

	if showBytecode {
		fmt.Printf("==ByteCode==\n")
		for idx, op := range bytecode {
			if idx%4 == 0 && idx != 0 {
				fmt.Printf("\n")
			}
			fmt.Printf("0x%08x ", uint32(op))
		}
		fmt.Printf("\n============\n")
	}

	vm := reference.New()
	fmt.Printf("Running bytecode...\n")
	if trace {
		vm.Start(bytecode)
		for {
			if err := vm.Step(); err != nil {
				if err != reference.ErrHalt {
					fmt.Printf("Error running bytecode: %v\n", err)
				}
				break
			}
			fmt.Printf("Stack: %+v\n", vm.Stack)
		}

	} else {
		if err := vm.Execute(bytecode); err != nil {
			fmt.Printf("Error running bytecode: %v\n", err)
		}

		if showStack {
			fmt.Printf("Stack: %#v\n", vm.Stack)
		}
	}
}
