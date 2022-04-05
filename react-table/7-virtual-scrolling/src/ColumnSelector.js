import styles from './ColumnSelector.module.css';

export default function ColumnSelector({ columns }) {
  return (
    <div className={styles.ColumnSelector}>
      <div className={styles.Label}>Show Columns:</div>
      <div className={styles.Checkboxes}>
        {columns.map(column => (
          <div key={column.id} className={styles.Checkbox}>
            <label>
              <input type='checkbox' {...column.getToggleHiddenProps()} />
              {` ${column.Header}`}
            </label>
          </div>
        ))}
      </div>
    </div>
  );
}