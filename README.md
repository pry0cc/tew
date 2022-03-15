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
tew -x file.xml | tee output.txt
tew -x file.xml | httpx -json -o http.json
```

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
- [ ] Clean up DNSX Parsing module and sort unique

# Credit
- @hakluke - Thank you man for helping me fix that dumb bug :) 
- @vay3t - Go Help
- @BruceEdiger - Go Help
- @mortensonsam - Go help!!
