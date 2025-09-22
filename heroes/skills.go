package heroes

type Skill struct {
	Class       string
	Level       int
	Name        string
	Description string
	Damage      int
	ManaCost    int
}

var WarSmash = Skill{
	"warrior",
	1,
	"warSmash",
	"мощный удар",
	10,
	40,
}

var WarriorSkills = map[string]Skill{
	"мощный удар": WarSmash,
}
