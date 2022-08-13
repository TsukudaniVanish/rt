package main 

import  	"github.com/anaseto/gruid"

type Status struct {
    HP int 
    MaxHP int 
    Power int 
    Defence int 
}

type EnemyAI struct {
    Path []gruid.Point
}
