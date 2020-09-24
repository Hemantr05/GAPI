package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// docker run <container> cmd args
// go run main.go run cmd args
func main(){
	switch os.Args[1]{
		case "run":
			run()
		case "child":
			child()
		default:
			panic("What??")
	}
}

func run(){
        cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

        // go run main.go args[1] args[2] ....
       // cmd := exec.Command(os.Args[2], os.Args[3:]...)
        cmd.Stdin = os.Stdin
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr

	// syscall.CLONE_NEWUTS - New hostname
	// syscall.CLONE_NEWPID - New Process PID
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID, 
	}

        must(cmd.Run())
}




func child(){
        fmt.Printf("running %v PID %d\n", os.Args[2:], os.Getpid())

        // go run main.go args[1] args[2] ....
        cmd := exec.Command(os.Args[2], os.Args[3:]...)
        cmd.Stdin = os.Stdin
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}


        must(cmd.Run())
}


func must(err error){
	if err != nil{
		panic(err)
	}
}
