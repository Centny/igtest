M begin
$a=1
$b=2
/* block comment */
/* 
block comment 
*/
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

// $(/js/a)

SUB sub.ig kk=111111 aa=444444 -cookie -ig
SUB sub.ig kk=111111 aa=444444
$aa=555555
$kk=666666
SUB sub.ig -CTX

BC $a*$b
SET aa 123
GET aa
HR GET http:\/\/localhost #data=dkk
// P $dkk
HR POST http:\/\/localhost
HP http:\/\/localhost
HG http:\/\/localhost
EX /bin/echo abc
EX ls #data=ls
P $ls
$x1=123456
$x2=123654
EX /bin/echo $x1 $x2 #data=ex
P $x1 $x2 $ex
// $p='/tmp/kk kk.log'
// EX cat $p #data=ex
// P $ex
// EX pwd #data=pwd
P $CWD
$fp=@($CWD+/+sub.ig)
P file path $fp
SUB $fp kk=$fp aa=testing
$res=@[EX /bin/echo 1]
Y $res

W /tmp/igtest testing
R /tmp/igtest #data=igt
Y $igt
W /tmp/igtest2 {"abb":111}
R /tmp/igtest2 #data=igt2
Y $(/igt2/abb)
D /tmp/igtest
D /tmp/igtest2

$ary=["A","B"]
$len=$(/ary/@len)
P $(/ary/0)
M mmmm

//
$val=100
$cval=1>2
$cval2=$val<1
$cval3=$val*100
$cval4=$val/100
$cval5=@{$cval4}
P $cval $cval2 $cval3 $cval4 $cval5

$A_B=1111
P $A_B

$sss=a==a
Y $sss
Y a==a
$va=a
$vb=a
$vc=$va==$vb
Y $va==$vb
Y $vc
N a==1