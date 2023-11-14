package main

import (
	"fmt"
	"encoding/json"
	"os"
)

type Cgroups struct {
	Cgroups []Cgroup
}

type Cgroup struct {
	Name string `json:"name"`
	Path string `json:"path:`
	Param string `json:"param"`
	Value []byte `json:"value,string"`
}



func main() {
	limits := Cgroups {
		[]Cgroup{
			{
				"pids",
				"/sys/fs/cgroup/pids/container",
				"pids.max",
				[]byte("20"),
			},
			{
				"memory",
				"/sys/fs/cgroup/memory/container",
				"memory.limit_in_bytes",
				[]byte("1000000"),
			},
		},
	}

	b, err := json.MarshalIndent(limits,"","\t")
	if err != nil {
		fmt.Println("error", err)
	}

	os.Stdout.Write(b)

}
