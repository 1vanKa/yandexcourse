package main

import (
	"fmt"
	"log"
	"reflect"
)

type Spell interface {
	// название заклинания
	Name() string
	// характеристика, на которую воздействует
	Char() string
	// количественное значение
	Value() int
}

// CastReceiver — если объект удовлетворяет этом интерфейсу, то заклинание применяется через него
type CastReceiver interface {
	ReceiveSpell(s Spell)
}

func CastToAll(spell Spell, objects []interface{}) {
	for _, obj := range objects {
		CastTo(spell, obj)
	}
}

func CastTo(spell Spell, object interface{}) {
	// реализуйте эту функцию.
	if receiver_obj, ok := object.(CastReceiver); ok {
		receiver_obj.ReceiveSpell(spell)
		log.Printf("Spell %v casted to %#v", spell, object)
		return
	}
	value_obj := reflect.ValueOf(object)
	if value_obj.Kind() != reflect.Ptr || value_obj.Elem().Kind() != reflect.Struct {
		return
	}

	if field := value_obj.Elem().FieldByName(spell.Char()); field.IsValid() {
		if field.CanSet() {
			field.SetInt(int64(spell.Value()) + field.Int())
			log.Printf("Spell %v casted to %#v", spell, object)
		}
	}
}

type spell struct {
	name string
	char string
	val  int
}

func newSpell(name string, char string, val int) Spell {
	return &spell{name: name, char: char, val: val}
}

func (s spell) Name() string {
	return s.name
}

func (s spell) Char() string {
	return s.char
}

func (s spell) Value() int {
	return s.val
}

type Player struct {
	name   string
	health int
}

func (p *Player) ReceiveSpell(s Spell) {
	if s.Char() == "Health" {
		p.health += s.Value()
	}
}

type Zombie struct {
	Health int
}

type Daemon struct {
	Health int
}

type Orc struct {
	Health int
}

type Wall struct {
	Durability int
}

func main() {

	player := &Player{
		name:   "Player_1",
		health: 100,
	}

	enemies := []interface{}{
		&Zombie{Health: 1000},
		&Zombie{Health: 1000},
		&Orc{Health: 500},
		&Orc{Health: 500},
		&Orc{Health: 500},
		&Daemon{Health: 1000},
		&Daemon{Health: 1000},
		&Wall{Durability: 100},
	}

	CastToAll(newSpell("fire", "Health", -50), append(enemies, player))
	CastToAll(newSpell("heal", "Health", 190), append(enemies, player))

	fmt.Println(player)
}
