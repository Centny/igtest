Integration test framework by golang
=========

##Store Value Path
* /abc/v
* /abc/1/a
* /abc/@length

##Variable
* $a=1
* $b=$a
* $c=$a
* $d=@{$a+$b}
* $d=$(/abc/v)
* $d=$(/abc/@length) //array length
* $d=$(/abc/1) 		 //array value

##Expression
* @{$a*$b}
* @{$a+$b}
* @{$a/$b}
* @{$a/$b-$c}
* @{($a+$b)*$c}
* @{$a<$c}
* @($a+abc+$b) string join
* @[command expression]

##Commands
####Y $a 
assert value is valid

####N $a 
assert value is invalid

####BC $a*$b
execute the bc command

####SET path value
store value to context.

* path:value store path
* value:store value

```
SET /a/b 123
SET /a/b $a1
SET $a2 $a1
```
####GET path
get the value from context,equal $(path)

* path:value store path

```
GET /a/b
```

####HR method url k1=v1,...,^h1=v3,...,%f1=path #path
send http require to url by query argument and http header,setting response data to /path

* method:GET or POST
* url:target http url
* k1=v1:adding one request argument by name(k1)
* ^h1=v3:adding one request header by name(h1)
* %f1=path:post file to server by name(f1)
* \#path:set the response data to path

```
R POST http://www.google.com ka=a,kb=b,^Content-Type=image/png,%file=/tmp/aa
```

####HG url k1=v1,..,^h1=v3,
see ```R```

####HP url k1=v1,..,^h1=v3,
see ```R```

####EX cmd arg1 arg2 ...
execute external command

* cmd:external command
* arg1,arg2:external command arguments.

```
EX /bin/echo abc 123
```

####P val1 val2,...
print message

* val1,val2:the message will be printed.

```
P abc 123
```

####FOR k R|IN 1~5|A~b~c \<cmds\> ROF 
loop commands,end by ROF

* k:the loop item
* R:loop the value range
* IN:loop the option values
* 1~5:loop range
* A~b~c:loop option
* cmds:all cmd to loop
* ROF:the for loop end

```
FOR k R 1~5
 P $k
OFR

FOR k IN A~B~C
 P $k
OFR
```

####INC path
include the script file to self

* path:script path

```
INC /tmp/ab.s
```

####SUB path k1=v1 k2=v2
execute subscript by argument k1 and k2

* path:file path or http url
* k1=v1:transfter key value to subprocess.
* `-cooki` transfter cookie to subprocess.
* `-CTX` transfter context to subprocess(cookie and value).
* `-ig` ignore error for subprocess.

```
SUB /tmp/ab.s va=123 vb=abc
```

####R filepath
read file context

####W filepath content
write file

####D filepath
delete file