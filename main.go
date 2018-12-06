package main

import (
	"fmt"
	"github.com/ceph/go-ceph/rados"
	"os"
)

func main() {
	if len(os.Args) <= 3 {
		fmt.Println("Usage: cephb <cephname> <poolname> <[ls|del <oid>]>")
		return
	}

	cname := os.Args[1]
	poolname := os.Args[2]
	action := os.Args[3]

	cephConn, err := rados.NewConnWithUser(cname)
	if err != nil {
		panic(err)
	}
	err = cephConn.ReadConfigFile("ceph.conf")
	if err != nil {
		panic(err)
	}
	err = cephConn.Connect()
	if err != nil {
		panic(err)
	}
	defer cephConn.Shutdown()

	cephContext, err := cephConn.OpenIOContext(poolname)
	if err != nil {
		panic(err)
	}

	switch action {
	case "ls":
		if err = cephContext.ListObjects(func(oid string) {
			fmt.Println(oid)
		}); err != nil {
			panic(err)
		}
	case "del":
		if len(os.Args) <= 4 {
			fmt.Printf("[ERROR] invalid oid.\n")
			return
		}
		oid := os.Args[4]
		if err = cephContext.Delete(oid); err != nil {
			panic(err)
		}
	default:
	}
}
