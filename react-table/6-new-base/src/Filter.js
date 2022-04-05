import * as React from 'react';

import styles from './Filter.module.css';

export default function Filter({ onChange }) {
  const [value, setValue] = React.useState('');

  const onChangeWrapper = React.useCallback((event) => {
    const v = event.target.value.trim();
    setValue(v);
    onChange(v);
  }, [onChange]);

  return (
    <div className={styles.Filter}>
      <input type="text" value={value} placeholder='Search rows...' onChange={onChangeWrapper} />
    </div>
  );
}