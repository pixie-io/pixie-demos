# Tables are Hard: Basic Demo

Demo for the second setup in [this tutorial](https://blog.px.dev/tables-are-hard-3).

[Back to Top](../README.md)

In this step, we add virtual scrolling to the table to handle thousands of rows without hanging the browser.

How this was created:
* Start with the result from [step 1](../6-new-base/README.md)
* `npm install react-window`
* Use `FixedSizeList` from `react-window` within `Table.js`, which presents a few noticeable issues
* Fix issues with a bit of DOM bounding box math and some CSS tweaks (`useScrollbarSize.js`, `useContainerSize.js`, `App.css`, `Table.module.css`)
