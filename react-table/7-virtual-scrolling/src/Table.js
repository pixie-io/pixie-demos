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
import { FixedSizeList } from 'react-window';

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
          {/*
            * This has a few problems once we use 10,000 rows and add virtual scrolling:
            * - At 10,000 rows, sorting and filtering slow down quite a lot. We'll address that later.
            *   Virtual scrolling makes scrolling and resizing columns fast again, but sorting and filtering still chug.
            * - The table's height is no longer dynamic. This can be fixed by detecting the parent's dimensions.
            * - If the user's browser shows layout-impacting scrollbars (Firefox does so by default for example),
            *   the header is not part of the scrolling area and thus has a different width than the scroll body.
            *   This can be fixed by detecting how wide the scrollbar is and whether it's present, then using that
            *   to adjust the <thead/> width accordingly.
            */}
          <FixedSizeList
            itemCount={rows.length}
            height={300}
            itemSize={34}
          >
            {({ index, style }) => {
              const row = rows[index];
              prepareRow(row);
              return <tr {...row.getRowProps({ style })}>
                {row.cells.map(cell => (
                  <td {...cell.getCellProps()}>
                    {cell.render('Cell')}
                  </td>
                ))}
              </tr>
            }}
          </FixedSizeList>
        </tbody>
      </table>
    </div>
  );
}