package common

type Bet struct {
	agency    string
	name      string
	surname   string
	document  string
	birthdate string
	number    string
}

func NewBet(
	agency string, name string, surname string, document string, birthdate string, number string) *Bet {
	bet := &Bet{
		agency:    agency,
		name:      name,
		surname:   surname,
		document:  document,
		birthdate: birthdate,
		number:    number,
	}
	return bet
}
