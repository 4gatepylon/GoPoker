package testing

import "testing"
import "fmt"

import "github.com/4gatepylon/GoPoker/poker/cardset"

// Golang is smart and every file that is Test<string that doesn't start with a lowercase letter>
// will be run when you do "go test" and we can use fmt.Errof to throw errors

// I might upgrade to Bazel later (though I can write my own rule in that case....)

func TestRoyalFlushExist(t *testing.T) {

}

func TestRoyalFlushDoesNotExist(t *testing.T) {

}

func TestQuadsExist(t *testing.T) {

}

func TestQuadsDoNotExist(t *testing.T) {

}