import * as React from 'react';

import { useStaticData } from './utils/useData.js';
import Table from './Table.js';

import './App.css';

function App() {
  const data = useStaticData(100);
  return (
    <main className="App">
      <Table data={data} />
    </main>
  );
}

export default App;
