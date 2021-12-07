# Tables are Hard: Basic Demo

Demo for the fourth step in [this tutorial](https://blog.px.dev/tables-are-hard-2).

[Back to Top](../README.md)

In this step, we add features to resize, reorder, and hide columns.

How this was created:
* Start with the result from [step 3](../3-sort-filter/README.md)
* Add column resizing:
  * Add `useFlexLayout` from `react-table` to define how column widths should be calculated
  * Add `useResizeColumns` from `react-table` to add logic for manual resizing
    * When combined with `useFlexLayout`, the table remains full width, and resizing a column also resizes other columns to retain the same combined width.
    * If we had used `useBlockLayout` instead, the table's total width would change as the columns do, and the default width would not be 100%.
  * Using the props that `useResizeColumns` adds to each column, create a resize handle that invokes the plugin's methods.
  * Slightly adjust structure in HTML and CSS for the `Table` component so that sort and resize controls don't interfere with each other.
* Add column hiding:
  * This one doesn't need a plugin to be imported; it's already provided by `useTable`
  * Create a `ColumnSelector` component, which takes the `getToggleHiddenProps` method from each column in `useTable(...).allColumns` to make a checkbox toggle that column's visibility