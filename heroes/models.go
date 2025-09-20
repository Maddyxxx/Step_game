package heroes

type Hero interface {
	Move()
}

type ClassHero struct {
	Class     string
	Level     int
	Health    int
	Mana      int
	Skills    []interface{}
	Equipment []interface{}
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
