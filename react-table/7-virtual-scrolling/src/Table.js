import * as React from 'react';

import styles from './Table.module.css'
import Filter from './Filter.js';
import ColumnSelector from './ColumnSelector.js';
import { useContainerSize } from './utils/useContainerSize';
import { useScrollbarSize } from './utils/useScrollbarSize';

import {
  useTable,
  useFlexLayout,
  useGlobalFilter,
  useSortBy,
  useResizeColumns,
} from 'react-table';
import { FixedSizeList } from 'react-window';

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

  const { width: scrollbarWidth } = useScrollbarSize();

  const [fillContainer, setFillContainer] = React.useState(null)
  const fillContainerRef = React.useCallback((el) => setFillContainer(el), []);
  const { height: fillHeight } = useContainerSize(fillContainer);

  return (
    <div className={styles.root}>
      <header>
        <ColumnSelector columns={allColumns} />
        <Filter onChange={setGlobalFilter} />
      </header>
      <div className={styles.fill} ref={fillContainerRef}>
        <div {...getTableProps()} className={styles.Table}>
          <div className={styles.TableHead}>
            {headerGroups.map((group, gi) => (
              <div className={styles.Row} {...group.getHeaderGroupProps()}>
                {group.headers.map((column, ci) => (
                  <div className={styles.HeaderCell} {...column.getHeaderProps({
                    style: (gi === headerGroups.length - 1 && ci === group.headers.length - 1)
                    ? { marginRight: scrollbarWidth }
                    :  {}
                  })}>
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
                  </div>
                ))}
              </div>
            ))}
          </div>
          <div className={styles.TableBody} {...getTableBodyProps()}>
            {/*
              * This has a few problems once we use 10,000 rows and add virtual scrolling:
              * - At 10,000 rows, sorting and filtering slow down quite a lot. We'll address that later.
              *   Virtual scrolling speeds up scrolling and resizing columns, but sorting and filtering still chug.
              */}
            <FixedSizeList
              outerElementType={(props, ref) => (
                // Instead of handling complexity with when the scrollbar is/isn't visible for this basic tutorial,
                // we'll instead force the scrollbar to appear even when it isn't needed. Suboptimal, but out of scope.
                <div {...props} style={{ ...props.style, overflowY: 'scroll' }} forwardedRef={ref}></div>
              )}
              itemCount={rows.length}
              height={fillHeight - 56}
              itemSize={34}
            >
              {({ index, style }) => {
                const row = rows[index];
                prepareRow(row);
                return <div className={styles.Row} {...row.getRowProps({ style })}>
                  {row.cells.map(cell => (
                    <div className={styles.BodyCell} {...cell.getCellProps()}>
                      {cell.render('Cell')}
                    </div>
                  ))}
                </div>
              }}
            </FixedSizeList>
          </div>
        </div>
      </div>
    </div>
  );
}