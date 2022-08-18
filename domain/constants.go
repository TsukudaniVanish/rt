package domain 

import (
	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/rl"
)

const AttrReverse = 1 << iota 

const (
    PlayerName = "You"
    UIWidth = 80
    UIHight = 24
    LogLines = 2 
    StatusLines = 3
)

const (
	ColorFOV gruid.Color = iota + 1
	ColorPlayer 
	ColorEnemy
    ColorConsumable
    ColorLogPlayerAttack
    ColorLogEnemyAttack
    ColorLogSpecial
    ColorStatusHealthy
    ColorStatusWounded
)

const (
	Wall rl.Cell = iota
	Floor 
	MinCaveSize = 400
    MapWidth = UIWidth
    MapHight = UIHight - LogLines - StatusLines
)

const (
	MaxLOS = 10
)

const (
    ErrNoShow = "ErrNoShow"
    ErrNoTargeting = "error no targeting"
    HealRate = 50
)