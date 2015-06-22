# Locating main.go

Godo compiles and runs a file at the relative path `Gododir/main.go`. If the path does not exist at the current directory, parent directories are searched.

For example, given this directory structure:

```
mgutz/
    Gododir/
        main.go
    project1/
        Gododir/
            main.go
    project2

```

* Running `godo` inside of project1 uses `project1/Gododir/main.go`.
* Running `godo` inside of project2 uses `mgutz/Gododir/main.go`


