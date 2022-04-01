import * as React from 'react';

import { useStreamingData } from './utils/useData.js';
import Table from './Table.js';

import './App.css';

function App() {
  const data = useStreamingData();
  return (
    <main className="App">
      <Table data={data} />
    </main>
  );
}

export default App;
