# SourceCodex

A program to convert source code repos to epubs.

## Examples

The program receives as input the filenames to include:

```bash
$ fd -e go | sourcecodex
```

Here's another example with more filters and files sorted by last commit date:

```bash
$ git ls-tree -r --name-only HEAD \
    | grep -E '\.ts$' \
    | grep -v -E '(^test|__test__|\.test\.\w+$|\.d\.ts$)' \
    | xargs -I{} git log -1 --format='%at {}' -- {} \
    | sort -n -r \
    | cut -d ' ' -f2- \
    | sourcecodex -title "Cool project" -author "Cool author" -output "source.epub"
```

## To-do

-   Set metadata
-   Generate a cover
-   Syntax highlighting
-   Add links between files (like a dependency graph)
-   Add an index for symbols, with links to its definition, and references
-   Link symbols in the code to their entry on the index
-   Sort files with a recommended reading order
-   Autoformat code for a better ebook reading experience
-   Add line numbers
