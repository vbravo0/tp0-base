from common.utils import Bet 

def bet_from_string(s: str) -> Bet:
  fields = s.split(',')
  return Bet(fields[0], fields[1], fields[2], fields[3], fields[4], fields[5])

def bets_from_chunk(chunk: str) -> list[Bet]:
  lines = chunk.split("\n")
  bets = []
  for line in lines:
    bet = bet_from_string(line)
    bets.append(bet)
  return bets

def bet_documents_to_chunk(bets: list[Bet]) -> str:
  documents = [bet.document for bet in bets]
  return ','.join(documents)
