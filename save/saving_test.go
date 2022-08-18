package save

import (
	"bytes"
	"encoding/gob"

	"game"

	"testing"
)

func TestSaveLoad(t *testing.T){
	RegisterEntity()



	g := game.NewGame()

	println("Log")

	log := g.Logs
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err := enc.Encode(log)
	if err != nil {
		t.Fatal(err)
	}
	data := buf.Bytes()
	dec := gob.NewDecoder(bytes.NewReader(data))
	log2 := &[]game.LogEntry{}
	err = dec.Decode(log2)
	if err != nil {
		t.Fatal(err)
	}

	println("Map")
	gmap := g.Map
	buf = &bytes.Buffer{}
	enc = gob.NewEncoder(buf)
	err = enc.Encode(gmap)
	if err != nil {
		t.Fatal(err)
	}
	data = buf.Bytes()
	dec = gob.NewDecoder(bytes.NewReader(data))
	gmap2 := &game.GameMap{}
	err = dec.Decode(gmap2)
	if err != nil {
		t.Fatal(err)
	}

	println("ECS")
	ecs := g.ECS
	buf = &bytes.Buffer{}
	enc = gob.NewEncoder(buf)
	err = enc.Encode(ecs)
	if err != nil {
		t.Fatal(err)
	}
	data = buf.Bytes()
	dec = gob.NewDecoder(bytes.NewReader(data))
	ecs2 := &game.ECS{}
	err = dec.Decode(ecs2)
	if err != nil {
		t.Fatal(err)
	}

	println("Game")
	buf = &bytes.Buffer{}
	enc = gob.NewEncoder(buf)
	err = enc.Encode(g)
	if err != nil {
		t.Fatal(err)
	}
	data = buf.Bytes()
	dec = gob.NewDecoder(bytes.NewReader(data))
	g2 := &game.Game{}
	err = dec.Decode(g2)
	if err != nil {
		t.Fatal(err)
	}

	println("Encode Decode")

	data, err = EncodeNoGzip(g)
	if err != nil {
		t.Fatal(err)
	}
	_, err = DecodeNoGzip(data)
	if err != nil {
		t.Fatal(err)
	}
}