# Tables are Hard: Basic Demo

Demo for the second setup in [this tutorial](https://blog.px.dev/tables-are-hard-3).

[Back to Top](../README.md)

In this step, we add streaming data to a table that already has virtual scrolling, filtering, and sorting.
Then we update these existing features to account for the data rapidly changing underneath them.

How this was created:
* Start with the result from [step 2](../7-virtual-scrolling/README.md)
* Modify `useData.js` and `App.js` to use data that adds a batch new rows every second, which shows several issues:
  * Scroll position jump to the top every time a batch is added
  * Virtual scrolling doesn't notice this until you scroll again
  * Sorting and filtering both reset every time a batch comes in
* ??? (next commit fixes these issues)