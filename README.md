# troute
Go implementation of the traceroute command

### install
```bash
go install github.com/Wa4h1h/troute@latest
```
OR
### generate build
```bash
make
```

### Features
* UDP tracing(does not require root)
* ICMP tracing(requires root)
* Concurrent hops/probes execution

### Usage
```bash
Usage: troute [options] host
Use troute -h or --help for more information.
Options:
  -I	use icmp echo probes
  -U	use udp packet probes (default true)
  -ch int
    	number of concurrent ttls (only for UDP and ICMP) (default 1)
  -cp int
    	number of concurrent probes pro ttl (default 3)
  -max-ttl int
    	specifies the maximum number of hops (max ttl value) (default 30)
  -n int
    	number of probes pro ttl (default 3)
  -start-ttl int
    	specifies with what TTL to start (default 1)
  -t int
    	probe timeout in seconds (default 3)
```

### Example output
```bash
troute google.com                                                                
tracing google.com (172.217.18.14) with max hops 30
1	66.215.140.254 (66.215.140.254) 5.300333ms  5.317625ms  5.252667ms
2	i68973bd1.versanet.de. (104.151.59.209) 8.743375ms  8.756417ms  8.964833ms
3	24.40.148.208 (24.40.148.208) 10.90375ms  10.966458ms  10.92075ms
4	89.246.109.249 (89.246.109.249) 13.214125ms
	72.14.204.149 (72.14.204.149) 22.343709ms  22.352167ms
5	72.14.204.148 (72.14.204.148) 25.689083ms * *
6	* * *
7	172.253.50.150 (172.253.50.150) 26.865167ms  26.84725ms  26.921958ms
8	172.253.66.139 (172.253.66.139) 23.978917ms
	192.178.109.126 (192.178.109.126) 26.806792ms  26.817792ms
9	fra02s19-in-f14.1e100.net. (172.217.18.14) 12.294584ms  12.327ms  12.480791ms
```