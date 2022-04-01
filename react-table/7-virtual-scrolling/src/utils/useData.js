import * as React from 'react';

import './useData.css';

const TimestampCell = ({ value }) => (
  <span className='Cell-Timestamp'>
    <span>{new Date(value).toLocaleString()}</span>
  </span>
);
const LatencyCell = ({ value }) => {
  let color = 'bad';
  if (value <= 50) color = 'good';
  else if (value <= 100) color = 'weak';

  return <span className={`Cell-Latency ${color}`}>{ value }</span>;
};
const EndpointCell = ({ value }) => {
  return <span className='Cell-Endpoint'>{ value }</span>;
};
const StatusCell = ({ value }) => {
  const digit = +value.split(' ')[0];
  let range = 500;
  if (digit < 300) range = 200;
  else if (digit < 400) range = 300;
  else if (digit < 500) range = 400;
  return <span className={`Cell-StatusCode range-${range}`}>{ value }</span>;
};

function randomFrom(array) {
  return array[Math.floor(Math.random() * array.length)];
}

function generateRow(minTimestamp, maxTimestamp) {
  const status = randomFrom([
    '200 OK', '301 Moved Permanently', '404 Not Found', "418 I'm a teapot", '501 Not Implemented'
  ]);
  return {
    timestamp: minTimestamp + Math.floor(Math.random() * (maxTimestamp - minTimestamp)),
    latencyMs: 5 + Math.floor(Math.random() * 150),
    // Fun bit of lore: https://wiki.ubuntu.com/DevelopmentCodeNames
    endpoint: '/user/' + randomFrom([
      'bendy-badger', 'happy-hippo', 'giant-ape', 'grumpy-groundhog', 'phlegmatic-pheasant'
    ]),
    status,
  };
}

const columns = [
  { Header: 'Timestamp', Cell: TimestampCell, accessor: 'timestamp' },
  { Header: 'Latency',   Cell: LatencyCell,   accessor: 'latencyMs' },
  { Header: 'Endpoint',  Cell: EndpointCell,  accessor: 'endpoint' },
  { Header: 'Status',    Cell: StatusCell,    accessor: 'status' },
];

// react-table expects memoized columns and data, so we export a React hook to permit doing that.
export function useStaticData(numRows = 20) {
  const data = React.useMemo(() => {
    const maxTimestamp = Date.now();
    const minTimestamp = maxTimestamp - (1000 * 60 * 60 * 24 * 7);
    return Array(numRows).fill(0).map(() => generateRow(minTimestamp, maxTimestamp));
  }, [numRows]);
  return { columns, data };
}
