IGR
======
ig script runner.

###Install
```
go install github.com/Centny/igtest/igr

```

###Usage
``` 
igr [-l -m mode -R type -r file] file ...
 -h show this
 -l show log
 -R the report data type
 -r the report store path
 -m mode in YOE/EYOE/NONE
  YOE assert on execute all command
  EYOE assert on expression
  NONE no assert in command(ony Y/N)
```