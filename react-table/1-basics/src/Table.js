import styles from './Table.module.css'

export default function Table({ data: { headers, rows } }) {
  return (
    <table className={styles.Table}>
      <thead>
        {headers.map(h => <th>{h}</th>)}
      </thead>
      <tbody>
        {rows.map((row) => (
          <tr>
            {row.map(cell => <td>{cell}</td>)}
          </tr>
        ))}
      </tbody>
    </table>
  );
}