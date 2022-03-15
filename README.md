# tew
`tew` is a simple, quick 'n' dirty nmap parser for converting nmap xml output files to IP:Port notation.

For example:

```
1.1.1.1:80
1.1.1.1.1:443
```

This is useful for internal penetration tests and can be piped to httpx easily. As it is go, it compiles into a neat and tidy binary! 

# Installation

## Go install
```
go install github.com/pry0cc/tew@latest
```
## Binaries
Binaries are available for most platforms and archectectures. - todo

# Usage
```
# Run Nmap and save to XML output

nmap -T4 1.1.1.1 8.8.8.8 -oX file.xml

tew file.xml
tew file.xml | tee output.txt
tew file.xml | httpx -json -o http.json
```

## Future Ideas
```
# Perhaps we could process dnsx output and correlate for vhosts?
tew file.xml --dnsx dnsx.txt | httpx 
```

# Todo
- [ ] Use proper flags library
- [ ] Create auto build using github ci & autobuild
- [ ] Add ability to import and use dnsx JSON & text output files

# Credit
- @hakluke - Thank you man for helping me fix that dumb bug :) 
