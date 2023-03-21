# SourceCodeX

A program to convert source code repos to epubs.

## Example

```bash
$ fd -e go | git ls-files | sourcecodex
```

## To-do

-   Generate epub
-   Syntax highlighting
-   Add links between files (like a dependency graph)
-   Add an index for symbols, with links to its definition, and references
-   Link symbols in the code to their entry on the index
