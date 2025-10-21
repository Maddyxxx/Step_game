package heroes

type Hero interface {
	TellHistory() string
	Move()
	Attack() int
	UseAbility(ability string) (string, int)
	LevelUp()
	ShowEquipment() []string
	ChangeEquipment(newEquipment Equipment)
}

type ClassHero struct {
	Class     string
	Level     int
	Health    int
	Mana      int
	Skills    map[string]Skill
	Equipment map[string]Equipment
}

func (h ClassHero) TellHistory() string {
	return Histories[h.Class]
}

func (h ClassHero) Move() {

}

func (h ClassHero) Attack() int {
	return h.Level + h.Equipment["weapon"].Strength
}

func (h ClassHero) UseAbility(ability string) (string, int) {
	for _, skill := range h.Skills {
		if skill.Name == ability {
			if h.Mana > skill.ManaCost {
				h.Mana -= skill.ManaCost
				return "", skill.Damage + h.Level*10
			} else {
				return "недостаточно маны", 0
			}
		} else {
			return "ошибка, скилл не найден", 0
		}
	}
	return "ошибка, скилл не применен", 0
}

func (h ClassHero) LevelUp() {
	h.Level++
	h.Health += 10
	// каждый новый уровень дать возможность прокачать любой скилл либо выбрать новый
	//h.SkillUp()
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

var Heroes = map[string]ClassHero{
	"warrior": {
		Class:     "warrior",
		Level:     1,
		Health:    100,
		Mana:      0,
		Skills:    WarriorSkills,
		Equipment: WarriorStartEquipment,
	},
}

var HeroList = []string{
	"warrior",
}

var Histories = map[string]string{
	"warrior": "Эта история про могучего воителя, который по преданиям был потомком самого Зевса",
}
