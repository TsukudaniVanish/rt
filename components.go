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

func (st *Status) Heal(n int) (healedHP int) {
    st.HP += n 
    if st.HP > st.MaxHP {
        n -= st.HP - st.MaxHP
        st.HP = st.MaxHP
    }
    healedHP = n 
    return 
}

// style contains information relative to default graphical represantation of an entity
type Style struct {
    Rune rune 
    Color gruid.Color
}

type Inventory struct {
    Items []int 
}
