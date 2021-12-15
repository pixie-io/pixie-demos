import * as React from 'react';

import styles from './Table.module.css'
import Filter from './Filter.js';
import { useTable, useGlobalFilter, useSortBy } from 'react-table';

export default function Table({ data: { columns, data } }) {
  const reactTable = useTable({
      columns,
      data
    },
    useGlobalFilter,
    useSortBy
  );

  const {
    getTableProps,
    getTableBodyProps,
    headerGroups,
    rows,
    prepareRow,
    setGlobalFilter
  } = reactTable;

  return (
    <>
      <Filter onChange={setGlobalFilter} />
      <table {...getTableProps()} className={styles.Table}>
        <thead>
          {headerGroups.map(group => (
            <tr {...group.getHeaderGroupProps()}>
              {group.headers.map(column => (
                <th {...column.getHeaderProps(column.getSortByToggleProps())}>
                  {column.render('Header')}
                  <span>
                    {column.isSorted ? (
                      column.isSortedDesc ? ' 🔽' : ' 🔼'
                    ): ''}
                  </span>
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
    </>
  );
}