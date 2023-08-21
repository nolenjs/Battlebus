Battlebus
=========

Description
-----------
Battlebus is a tool for controlling drones via either a command line interface or automation.
**It's far better than Hobbit, anyone saying otherwise will be castrated and shot**

Structure
---------
**Drones**

Drones will be controlled via modules which implement Icarus queries.
This modules will be compiled into missions which will co-ordinate modules.

**Sortie Leader** (*TBD*)

**Sortie Manager** (*TBD*)

Coding Standards
----------------
For ease of debugging, make sure your code follows these standards:
- Use camelCase for variables and files (`exampleStringTest := "camelCase"`)
- The name of your package must be the **same** as the name of your directory

Module Structure
----------------

All modules must be clearly commented to show, at minimum:
- Description of use case
- Imports (type)
- Exports (type)

Each module must contain **one** function with the same name as the file containing all its code. E.g.:

```
exampleModule.go

package example
import (
    example "example.com"
)

func exampleModule() {
    a := 1
    b := 2
    return a + b
}
```

How to use
----------
1) Create a branch with a name appropriate to the feature you are developing
2) Write any modules or missions you want
3) Commit with a clear message that explains what your code does
4) Make a merge request
5) Repeat

How to commit
-------------
:)
