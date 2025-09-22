package heroes

type Equipment struct {
	Class      string
	Type       string
	Name       string
	Strength   int
	Durability int
}

var wSword1 = Equipment{
	"warrior",
	"меч",
	"Тренировочный меч",
	1,
	50,
}

var wChest1 = Equipment{
	"warrior",
	"грудак",
	"Тренировочный грудак",
	1,
	50,
}

var WarriorStartEquipment = map[string]Equipment{
	"меч":    wSword1,
	"грудак": wChest1,
}
