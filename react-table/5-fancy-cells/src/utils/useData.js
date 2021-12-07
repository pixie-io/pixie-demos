import * as React from 'react';

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

const genericHeaders = ['ID', 'Name', 'Friend', 'Score', 'Temperament'];
const genericDataFuncs = [
  () => Math.ceil(Math.random() * 100),
  sillyName,
  sillyWord,
  () => Math.floor(Math.random() * 10_000) / 100,
  () => randomFrom(['Goofy', 'Wacky', 'Silly', 'Funny', 'Serious']),
];
const fancyCellRenderers = [
  function IdCell({ value }) {
    return <>{value}</>;
  },
  function NameCell({ value }) {
    return <strong>{value}</strong>;
  },
  function WordCell({ value }) {
    return <span style={{ fontStyle: 'italic' }}>{value}</span>
  },
  function ScoreCell({ value }) {
    const color = value < 50 ? 'pink' : 'aquamarine';
    return <span style={{ color }}>{value}</span>;
  },
  function TemperamentCell({value}) {
    return <>{value}</>;
  },
];

// react-table expects memoized columns and data, so we export a React hook to permit doing that.
export default function useData(numRows = 20, numCols = 5) {
  const columns = React.useMemo(() => (
    Array(numCols).fill(0).map((_, h) => {
      const name = genericHeaders[h % genericHeaders.length];
      const id = `col${h}`;
      return {
        Header: name,
        Cell: fancyCellRenderers[h % fancyCellRenderers.length],
        accessor: id
      };
    })
  ), [numCols]);

  const data = React.useMemo(() => (
    Array(numRows).fill(0).map(() => {
      const row = {};
      for (let c = 0; c < numCols; c++) {
        row[`col${c}`] = genericDataFuncs[c % genericDataFuncs.length]();
      }
      return row;
    })
  ), [numRows, numCols]);

  return { columns, data };
}