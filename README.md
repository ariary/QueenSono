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
- client :size of packet
- client: delay
- server: truncated (accept truncated data, with delay in parameter, wait for delay*NB packet and close all after
- client: filename or simple message
- server: save in file or print out

- server: integrity check
- server: crossbar advancement
- client: noreply (does not waiting for reply, in case they are filtered by fw)
