
$a=1
$b=2

//
FOR I R 0~3 //jkkk
 P $I
ROF

//
$total=0
FOR I R 0~100
 $total=@{$total+$I}
ROF
P $total

FOR I IN A~B~$total
 P $I
ROF

$pp=aaa
FOR I IN A~B~$total
 $pa=@($pp+$I)
 P $pa
ROF

$js={"a":11,"b":0}
P $(/js/a)

$(/js/a)=111111

P $(/js/a)

Y $(/js/a)
N $(/js/b)

$(/js/a)

SUB sub.ig kk=111111 aa=444444 -cookie -ig
SUB sub.ig kk=111111 aa=444444

BC $a*$b
SET aa 123
GET aa
HR GET http:\/\/192.168.1.1 #data=dkk
// P $dkk
HR POST http:\/\/192.168.1.1
HP http:\/\/192.168.1.1
HG http:\/\/192.168.1.1
EX /bin/echo abc
EX ls #data=ls
P $ls
$x1=123456
$x2=123654
EX /bin/echo $x1 $x2 #data=ex
P $x1 $x2 $ex
$p='/tmp/kk kk.log'
EX cat $p #data=ex
P $ex
EX pwd #data=pwd
// P $pwd
$fp=@($pwd+/+test/sub.ig)
// P $fp
SUB $fp kk=$fp aa=testing
$res=@[EX /bin/echo 1]
Y $res