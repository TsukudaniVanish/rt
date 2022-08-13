package main 

import (
	"fmt"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/ui"
)

type LogEntry struct {
    Text string
    Color gruid.Color
    Dups int // conseutive duplicates of same message
}

func (e *LogEntry) String() (s string){
    if e.Dups == 0 {
        s = e.Text
        return 
    }
    s = fmt.Sprintf("%s (%dx)", e.Text, e.Dups)
    return
}
