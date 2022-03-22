import * as React from 'react';

import styles from './Table.module.css'
import Filter from './Filter.js';
import ColumnSelector from './ColumnSelector.js';
import {
  useTable,
  useFlexLayout,
  useGlobalFilter,
  useSortBy,
  useResizeColumns,
} from 'react-table';

import './Table.css';

export default function Table({ data: { columns, data } }) {
  const reactTable = useTable({
      columns,
      data
    },
    useFlexLayout,
    useGlobalFilter,
    useSortBy,
    useResizeColumns
  );

  const {
    getTableProps,
    getTableBodyProps,
    headerGroups,
    rows,
    allColumns,
    prepareRow,
    setGlobalFilter
  } = reactTable;

  return (
    <div className='Table-root'>
      <header>
        <ColumnSelector columns={allColumns} />
        <Filter onChange={setGlobalFilter} />
      </header>
      <table {...getTableProps()} className={styles.Table}>
        <thead>
          {headerGroups.map(group => (
            <tr {...group.getHeaderGroupProps()}>
              {group.headers.map(column => (
                <th {...column.getHeaderProps()}>
                  <div {...column.getSortByToggleProps()}>
                    {column.render('Header')}
                    <span>
                      {column.isSorted ? (
                        column.isSortedDesc ? ' 🔽' : ' 🔼'
                      ): ''}
                    </span>
                  </div>
                  <div {...column.getResizerProps()} className={[styles.ResizeHandle, column.isResizing && styles.ResizeHandleActive].filter(x=>x).join(' ')}>
                    &#x22EE;
                  </div>
                </th>
              ))}
            </tr>
          ))}
        </thead>
        <tbody {...getTableBodyProps()}>
          {rows.map((row) => {
            prepareRow(row);
            return (
              <tr {...row.getRowProps()}>
                {row.cells.map(cell => (
                  <td {...cell.getCellProps()}>
                    {cell.render('Cell')}
                  </td>
                ))}
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}