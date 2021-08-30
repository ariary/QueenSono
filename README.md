# QueenSono
A Golang Package for Data Exfiltration with ICMP protocol. Try to imitate https://github.com/ytisf/PyExfil (and others) with the idea that the target machine does not necessary have python installed (so provide a binary could be usefupl)

***Notes:***
- only work on Linux 
- need cap_net_raw capabilities
- Only rely on the fact that ICMP protocol isn't monitored (could be detected with a rigorous inspection of ICMP packet content or frequency)

## ICMP
https://github.com/cyb3rw01f/icmpExfiltrater (avec command on serverside)
https://github.com/martinoj2009/ICMPExfil

## options
- size of packet
- delay
- filename or simple message 
- integrity check
- server: crossbar advancement
