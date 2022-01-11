package main

import (
	_ "github.com/lib/pq"
	"testing"
)

func TestPgInstance(t *testing.T){
	getPgContainerInstance()
	getPgContainerInstance()
	getPgContainerInstance()
	getPgContainerInstance()
	getPgContainerInstance()
}
