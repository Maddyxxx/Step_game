package heroes

type Hero interface {
	Move()
	Attack()
	ShowEquipment() []string
	ChangeEquipment(newEquipment Equipment)
}

type ClassHero struct {
	Class     string
	Level     int
	Health    int
	Mana      int
	Skills    []Skills
	Equipment []Equipment
}

func (h ClassHero) Move() {

}

func (h ClassHero) Attack() {

}

func (h ClassHero) ShowEquipment() []string {
	equipment := make([]string, len(h.Equipment))
	for _, e := range h.Equipment {
		equipment = append(equipment, e.Name)
	}
	return equipment
}

func (h ClassHero) ChangeEquipment(newEquipment Equipment) {

}

// Heroes
var (
	Warrior = ClassHero{
		Class:  "warrior",
		Level:  1,
		Health: 100,
		Mana:   0,
	}
)

var Heroes = []string{
	Warrior.Class,
}
