# tew
`tew` is a simple, quick 'n' dirty nmap parser for converting nmap xml output files to IP:Port notation.

For example:

```
1.1.1.1:80
1.1.1.1.1:443
```

This is useful for internal penetration tests and can be piped to httpx easily. As it is go, it compiles into a neat and tidy binary! 

# Installation

## Go install
```
go install github.com/pry0cc/tew@latest
```

## Binaries
Binaries are available for most platforms and archectectures in the [releases page](https://github.com/pry0cc/tew/releases/latest).

# Usage
```
# Run Nmap and save to XML output

nmap -T4 1.1.1.1 8.8.8.8 -oX file.xml

tew -x file.xml
tew -x file.xml -o output.txt
tew -x file.xml | httpx -json -o http.json
```

## Stdin support
```
cat example-data/example-1/nmap.xml | go run main.go -x -
```

93.184.216.34:80
93.184.216.34:443
1.1.1.1:80
8.8.8.8:53
8.8.8.8:443
8.8.4.4:53
8.8.4.4:443
1.0.0.1:53
1.1.1.1:53
1.1.1.1:443
1.0.0.1:80
1.0.0.1:443

## DNSx Parsing
If you want to correlate DNSx JSON output, simply generate a JSON file and import it using the following syntax.
```
subfinder -d domain.com -o subs.txt
dnsx -l subs.txt -json -o dns.json
cat dns.json | jq -r '.a[]' | tee ips.txt
nmap -T4 -iL ips.txt -oX nmap.xml

tew -x nmap.xml -dnsx dns.json --vhost | httpx -json -o http.json
```

# Todo
- [x] Create auto build using github ci & autobuild
- [x] Add Arm64 for Darwin to Build
- [x] Use proper flags library
- [x] Add ability to import and use dnsx JSON & text output files - working on it!
- [x] Clean up DNSX Parsing module and sort unique
- [x] Add output text file as option
- [ ] Test on Windows, Linux & Mac for cross-compatibility

# Credit
- @hakluke - Thank you man for helping me fix that dumb bug :) 
- @vay3t - Go Help
- @BruceEdiger - Go Help
- @mortensonsam - Go help!!
- @xstp - stdin support
- https://www.golangprograms.com - A lot of the code here is copy-pasted from the internet, at the time of writing, my go skills are copy-paste :P And that's ok if it works, right?
