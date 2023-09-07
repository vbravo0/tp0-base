package common

import (
	"fmt"
	"strings"
)

func bet_to_string(bet *Bet) string {
	return fmt.Sprintf(
		"%v,%v,%v,%v,%v,%v",
		bet.agency,
		bet.name,
		bet.surname,
		bet.document,
		bet.birthdate,
		bet.number,
	)
}

func bet_from_string(s string) *Bet {
	fields := strings.Split(s, ",")
	return NewBet(fields[0], fields[1], fields[2], fields[3], fields[4], fields[5])
}

func bets_from_array(a []string) []*Bet {
	bets := make([]*Bet, 0)
	for _, s := range a {
		bets = append(bets, bet_from_string(s))
	}
	return bets
}

func bets_to_chunk(bets []*Bet) string {
	chunk := ""
	for i, bet := range bets {
		if i == len(bets)-1 {
			chunk += bet_to_string(bet)
		} else {
			chunk += bet_to_string(bet) + "\n"
		}
	}
	return chunk
}

func bets_array_add_agency(agency string, a []string) []string {
	res := make([]string, 0)
	for _, s := range a {
		res = append(res, agency+","+s)
	}
	return res
}

func bet_documents_from_chunk(chunk string) []string {
	documents := make([]string, 0)
	if len(chunk) == 0 {
		return documents
	}

	lines := strings.Split(chunk, ",")
	for _, line := range lines {
		documents = append(documents, line)
	}
	return documents
}
