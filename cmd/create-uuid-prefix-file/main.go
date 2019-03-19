package main

import (
	"flag"

	"github.com/m-lab/go/flagx"
	"github.com/m-lab/go/rtx"
	"github.com/m-lab/uuid/prefix"
)

var (
	filename = flag.String("filename", "/var/local/uuid/prefix",
		"The file to which the prefix should be written. If the file "+
			"does not exist, it will be created.  If the file's directory "+
			"does not exist, creation will fail.")
)

func main() {
	flag.Parse()
	flagx.ArgsFromEnv(flag.CommandLine)

	rtx.Must(prefix.Generate(*filename), "Could not create prefix file.")
}
