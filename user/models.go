package user

type UserModel struct {
	ID           int
	TelegramID   int
	Username     string
	CurrentState string
	Hero         string
	Inventory    []string
	Progress     []string
}
