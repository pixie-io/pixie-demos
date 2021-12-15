function randomFrom(array) {
  return array[Math.floor(Math.random() * array.length)];
}

/** Make a silly word that rhymes with Goomba (the Mario mushroom enemies) */
function sillyWord() {
  const leadConsonants = ['B', 'D', 'F', 'G', 'L', 'T', 'V', 'Z'];
  const middles = ['oom', 'oon', 'um', 'un'];
  const midConsonants = ['b', 'd']
  const ends = ['a', 'ah', 'u', 'uh', 'o'];

  const word = randomFrom(leadConsonants) + randomFrom(middles) + randomFrom(midConsonants) + randomFrom(ends);

  // Try again if we accidentally picked something offensive
  if (['Goombah'].includes(word)) return sillyWord();
  else return word;
}

function sillyName() {
  return `${sillyWord()} ${sillyWord().substring(0, 3)}`;
}

export default function useData(numRows = 20, numCols = 5) {
  const headers = ['ID', 'Name', 'Friend', 'Score', 'Temperament'];
  const genFuncs = [
    () => Math.ceil(Math.random() * 100),
    sillyName,
    sillyWord,
    () => Math.floor(Math.random() * 10_000) / 100,
    () => randomFrom(['Goofy', 'Wacky', 'Silly', 'Funny', 'Serious']),
  ];
  return {
    headers: Array(numCols).fill(0).map((_, h) => headers[h % headers.length]),
    rows: Array(numRows).fill(0).map(
      () => Array(numCols).fill(0).map(
        (_, c) => genFuncs[c % genFuncs.length]()
      )
    ),
  };
}