package main

import (
	"flag"
	"os/exec"

	"github.com/golang/glog"
)

const a1 string = "-c"

func main() {
	flag.Parse()
	defer glog.Flush()

	const a2 = a1
	glog.Info("run command")
	ret := exec.Command("bash", a2, "id")
	glog.Infof("%s", ret)
}
