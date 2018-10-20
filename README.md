# 1337Speak Detector!
A rather simple 1337speak detector that uses a mapping dictionary and wordlist
to detect and higlight instances of 1337speak written to stdin.

## Example
Run `./1337Detect -h` for usage.

Example below is on binary 
```
$ strings test.jpg | ./1337Detect
Dict size:52
Wordlist size:5460
9Nice one for checking metadata. n07 50 53cr37 m374 d474
```

`53cr37` should be highlighted in your terminal.
