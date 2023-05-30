go mod tidy

[websocket_client](http://www.easyswoole.com/wstool.html)



ws://127.0.0.1:8080?user=123

ws://127.0.0.1:8080?user=456

ws://127.0.0.1:8080?user=789



INFO[0004] update protocol                               id=abc listen=":8080" module=server

INFO[0004] get userId                                    id=abc listen=":8080" module=server

INFO[0004] add user: 123                                 id=abc listen=":8080" module=server

INFO[0007] update protocol                               id=abc listen=":8080" module=server

INFO[0007] get userId                                    id=abc listen=":8080" module=server

INFO[0007] add user: 456                                 id=abc listen=":8080" module=server

INFO[0014] update protocol                               id=abc listen=":8080" module=server

INFO[0014] get userId                                    id=abc listen=":8080" module=server

INFO[0014] add user: 789                                 id=abc listen=":8080" module=server

INFO[0026] recv msg: hello, every one from user: 123    

INFO[0026] 123 send to 789: hello, every one ----from 123 

INFO[0026] 123 send to 456: hello, every one ----from 123 

INFO[0038] recv msg: hello, world from user: 456        

INFO[0038] 456 send to 123: hello, world ----from 456   

INFO[0038] 456 send to 789: hello, world ----from 456   

INFO[0051] recv msg: hello, websocket! from user: 789   

INFO[0051] 789 send to 123: hello, websocket! ----from 789 

INFO[0051] 789 send to 456: hello, websocket! ----from 789 