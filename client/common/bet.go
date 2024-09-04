package common

import (
	"strconv"
)

type Bet struct {
	agency    int
	firstname string
	surname   string
	document  string
	bornDate  string
	betValue  string
}

func NewSingleBet(agency int, firstname string, surname string, document string, bornDate string, betValue string) *Bet {
	singleBet := &Bet{
		agency:    agency,
		firstname: firstname,
		surname:   surname,
		document:  document,
		bornDate:  bornDate,
		betValue:  betValue,
	}

	return singleBet
}

func (bet *Bet) ToArray() []string {
	array := []string{
		strconv.Itoa(bet.agency),
		bet.firstname,
		bet.surname,
		bet.document,
		bet.bornDate,
		bet.betValue,
	}

	return array
}

func (bet *Bet) GetDocument() string {
	return bet.document
}

func (bet *Bet) GetBetNumber() string {
	return bet.betValue
}
