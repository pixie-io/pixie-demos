# Tables are Hard: Basic Demo

Demo for the second step in [this tutorial](https://blog.px.dev/tables-are-hard-2).

[Back to Top](../README.md)

In this step, we install the [react-table](https://react-table.tanstack.com) library and use it to render our basic data.

How this was created:
* Start with the result from [step 1](../1-basics/README.md)
* `npm install react-table`
* Follow `react-table`'s [quick start](https://react-table.tanstack.com/docs/quick-start) to convert our basic HTML table into a `react-table` consumer.
* Changed the `makeData` function in `src/makeData.js` to provide what `react-table` expects.
* Renamed `makeData` to `useData` in order to use `React.useMemo`.
* Re-run `npm start` (need to restart this when changing dependencies and adding/renaming files)