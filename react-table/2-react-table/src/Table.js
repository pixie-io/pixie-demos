import styles from './Table.module.css'
import { useTable } from 'react-table';

export default function Table({ data: { columns, data } }) {
  const reactTable = useTable({ columns, data });

  const {
    getTableProps,
    getTableBodyProps,
    headerGroups,
    rows,
    prepareRow
  } = reactTable;

  return (
    <table {...getTableProps()} className={styles.Table}>
      <thead>
        {headerGroups.map(group => (
          <tr {...group.getHeaderGroupProps()}>
            {group.headers.map(column => (
              <th {...column.getHeaderProps()}>
                {column.render('Header')}
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
  );
}