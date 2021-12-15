import useData from './utils/useData.js';
import Table from './Table.js';

import './App.css';

function App() {
  const data = useData();
  return (
    <div className="App">
      <main>
        <Table data={data} />
      </main>
    </div>
  );
}

export default App;
