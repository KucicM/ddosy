package ddosy_test

import (
	"log"
	"testing"

	ddosy "github.com/kucicm/ddosy/app"
)

func TestRepositoryBasic(t *testing.T) {
	rep := ddosy.NewRepository(":memory:")
	log.Println(rep.InsertNew(ddosy.LoadTask{}))
	log.Println(rep.InsertNew(ddosy.LoadTask{}))
	log.Println(rep.InsertNew(ddosy.LoadTask{}))
	rep.Close()

	
}
