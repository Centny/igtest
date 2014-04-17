Example
======

```
//comment
/*
comment block
*/
$a=1  //define variable
$b=2
//
$total=0
FOR I R 0~100 //loop
 $total=@{$total+$I}
ROF
P $total //print

FOR I IN A~B~$total //for in
 P $I
ROF

$js={"a":11,"b":0} 	//store object value(json format)
P $(/js/a) 			//get value by path

$(/js/a)=111111		//set value by path
P $(/js/a)

Y $(/js/a)			//assert value true(0,nil,empty string is false)
N $(/js/b)			//assert value false

//execute subprocess by cooke and ignore error
//transfter value by key kk and aa
SUB sub.ig kk=111111 aa=444444 -cookie -ig 
SUB sub.ig kk=111111 aa=444444

BC $a*$b  	//execute math,equal @{$a+$b}

SET aa 123 	//store value to path aa,equal $(/aa)=123
GET aa		//get value by path aa,equal $(/aa)

//execute one http request to url,
//and set the response data to path /hr/ab
HR GET http:\/\/192.168.1.1 #data=/hr/ab
HR POST http:\/\/192.168.1.1
HP http:\/\/192.168.1.1
HG http:\/\/192.168.1.1

//execute external command
//and set the response data to path /ex/eh
EX /bin/echo abc #data=/ex/eh

P $CWD //get the current script file location

D filepath //delete file
R filepath #abc //read file to abc
W filepath $abc //write abc value to file

```