import * as React from 'react';

/**
 * Oversimplified way to get the scroll container dimensions of an element.
 * If those dimensions vary with the size of those children, and those children try to use this,
 * there will be a loop. Use a more complex technique like `react-virtualized-autosizer` for those scenarios.
 */
export function useContainerSize(el) {
  return React.useMemo(() => {
    if (el) {
      return {
        width: el.clientWidth,
        height: el.clientHeight,
      };
    }

    return { width: 0, height: 0 };
  }, [el]);
}
