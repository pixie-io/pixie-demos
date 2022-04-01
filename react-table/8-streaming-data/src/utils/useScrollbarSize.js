import * as React from 'react';

/**
 * Determines the dimensions of the browser scrollbars when they would be shown.
 * Browsers may overlay scrollbars rather than putting them in-layout. In those cases they're treated as 0-width/height.
 */
 export function useScrollbarSize() {
  return React.useMemo(() => {
    const scroller = document.createElement('div');
    scroller.setAttribute('style', 'width: 100vw; height: 100vh; overflow: scroll; position: absolute; top: -100vh;');
    document.body.appendChild(scroller);
    const width = scroller.offsetWidth - scroller.clientWidth;
    const height = scroller.offsetHeight - scroller.clientHeight;
    document.body.removeChild(scroller);
    return { width, height };
  }, []);
}
