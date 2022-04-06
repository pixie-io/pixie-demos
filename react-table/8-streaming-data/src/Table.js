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

const TableContext = React.createContext(null);

/**
 * By memoizing this, we ensure that react-window can recycle rendered rows that haven't changed when new data comes in.
 */
const BodyRow = React.memo(({ index, style }) => {
  const { rows, instance: { prepareRow } } = React.useContext(TableContext);
  const row = rows[index];
  prepareRow(row);
  return (
    <div className={styles.Row} {...row.getRowProps({ style })}>
      {row.cells.map(cell => (
        <div className={styles.BodyCell} {...cell.getCellProps()}>
          {cell.render('Cell')}
        </div>
      ))}
    </div>
  );
});

/**
 * Setting outerElementType on FixedSizeList lets us override properties on the scroll container. However, FixedSizeList
 * redraws this on every data change, so the wrapper component needs to be memoized for scroll position to be retained.
 *
 * Note: If the list is sorted such that new items are added to the top, the items in view will still change
 * because the ones that _were_ at that scroll position were pushed down.
 * This can be accounted in a more complete implementation, but it's out of scope of this demonstration.
 */
const ForcedScrollWrapper = React.memo((props, ref) => (
  // Instead of handling complexity with when the scrollbar is/isn't visible for this basic tutorial,
  // instead force the scrollbar to appear even when it isn't needed. Not great, but out of scope.
  <div {...props} style={{ ...props.style, overflowY: 'scroll' }} forwardedRef={ref}></div>
));

const Table = React.memo(({ data: { columns, data } }) => {
  const reactTable = useTable({
      columns,
      data,
      autoResetSortBy: false,
      autoResetResize: false,
      disableSortRemove: true,
      initialState: {
        sortBy: [{ id: 'timestamp', desc: false }],
      },
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
    setGlobalFilter
  } = reactTable;

  const [latestFilter, setLatestFilter] = React.useState('');
  React.useEffect(() => {
    setGlobalFilter(latestFilter);
  }, [setGlobalFilter, latestFilter, data]);

  const { width: scrollbarWidth } = useScrollbarSize();

  const [fillContainer, setFillContainer] = React.useState(null)
  const fillContainerRef = React.useCallback((el) => setFillContainer(el), []);
  const { height: fillHeight } = useContainerSize(fillContainer);

  const [visibleStart, setVisibleStart] = React.useState(1);
  const [visibleStop, setVisibleStop] = React.useState(1);
  const viewportDetails = React.useMemo(() => {
    const count = visibleStop - visibleStart + 1;
    let text = `Showing ${visibleStart + 1} - ${visibleStop + 1} / ${rows.length} records`;
    if (rows.length === 500) text += ' (most recent only)';

    if (count <= 0) {
      text = 'No records to show';
    } else if (count >= rows.length) {
      text = '\xa0'; // non-breaking space
    }
    return text;
  }, [rows.length, visibleStart, visibleStop]);

  const onRowsRendered = React.useCallback(({ visibleStartIndex, visibleStopIndex }) => {
    setVisibleStart(visibleStartIndex);
    setVisibleStop(visibleStopIndex);
  }, []);

  const context = React.useMemo(() => ({
    instance: reactTable,
    rows,
    // By also watching reactTable.state specifically, we make sure that resizing columns is reflected immediately.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }), [reactTable, rows, reactTable.state]);

  return (
    <TableContext.Provider value={context}>
      <div className={styles.root}>
        <header>
          <ColumnSelector columns={allColumns} />
          <Filter onChange={setLatestFilter} />
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
              <FixedSizeList
                outerElementType={ForcedScrollWrapper}
                itemCount={rows.length}
                height={fillHeight - 56}
                itemSize={34}
                onItemsRendered={onRowsRendered}
              >
                {BodyRow}
              </FixedSizeList>
            </div>
          </div>
        </div>
        <div className={styles.ViewportDetails} style={{ marginRight: scrollbarWidth }}>
          {viewportDetails}
        </div>
      </div>
    </TableContext.Provider>
  );
});

export default Table;
