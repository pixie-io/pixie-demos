# Tables are Hard: Basic Demo

Demo for the third step in [this tutorial](https://blog.px.dev/tables-are-hard-2).

[Back to Top](../README.md)

In this step, we add filtering and sorting functionality.

How this was created:
* Start with the result from [step 2](../2-react-table/README.md)
* Add filtering:
  * Create a basic `Filter` component for the user to type anything that might exist in the table's data
  * Import `useGlobalFilter` from `react-table`
  * Add it to the list of plugins in the `useTable` call
  * Call the `setGlobalFilter` function that it provides when the input in our new `Filter` component changes.
* Add sorting:
  * Import `useSortBy` from `react-table`
  * Add it to the list of plugins like we did with `useGlobalFilter`
  * Add props to the table header to attach props and event handlers that `useSortBy` creates for us.
  * Add a visual representation of sort state, using props that `useSortBy` provides for us.