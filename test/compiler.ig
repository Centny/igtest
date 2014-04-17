$a=1
$b=$aa
//for 1
FOR i R 1~9 P $AA ROF
//for 2
FOR i R 1~9 
 P $AA ROF
//for 3
FOR i R 1~9 
 P $AA 
ROF
//for 4
FOR i R 1~9 
 P $a
 P $b 
ROF
//for 5
FOR i R 1~9 
 P $a
 FOR k IN A~B~C
  P $b 
 ROF
ROF
//for 6
FOR i R 1~9 
 P $AA 
 ROF
//for 7
FOR i R 1~9 
 P $AA ROF

HG http://github.com
HP http://github.com
HR POST http://github.com
BC $a+$a
P $a
SET $ab 123
GET $ab
EX /bin/echo abc
SUB abc