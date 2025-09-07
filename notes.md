
scope.yml
ingnore-issues.yml
data/
    scope-full.yml <- maybe call this full-inventory.yaml
    scope-current.yml <- not necessary
    issues-current.yml


---

scope.yml
ingnore-issues.yml
data/
    discovered-surface.yml
    discovered-issues.yml

---

overview.ymd <- a human-readable recap of scope-current and issues-current.yml

on every run:
scope ---upodate--> scope-full

        _/
      _/
     /

scope-full --scan--> scope-current.new   --diff with previous --> alert of new stuff, replace old file, update scope-full
                     issues-current.new  --diff with previous --> alert of new stuff, replace old file,


what is scope full exactly?
- a complete list of all assets ever seen by the scanner. including old, or deleted ones.

we don't need to keep a history of events for our functionality.
a history can always be reconstructed by the commit history, or 
by git-bisect

what is ignore-issues?
on every scan, the issues-current file is updated.
evefry single item in that list will be considered an issue, and put into a summary.
if some of the elements in that list are false positives, the identifier of that 
event must be put into ignore-issues.
on a subsequent scan, the issue will not be searched again, or will not be included in the list.


scope.yml
inspired by owasp-amass

```
scope:
  domains:
    - owasp.org
  ips:
    - 192.0.2.1
    - 192.168.0.10-192.168.0.20
  endpoints:
    - owasp.org:8727/teamname/

  cidrs:
    - 192.0.2.0/24
  ports:
    - 80
    - 443

blacklist:
  domains:
    - owasp.org

```

scope-full will be similar, maybe it will include endpoints:
- owasp.org/grafana
- owasp.org:8080/
- owasp.org:8080/asd
but probably endpoints are not necessary.
endpoints can be easily generated from permutations of what's in the scope, and the generation will 
clearly happen every time.
endpoints are what will end up in the scope-current.

it could be useful to enhance stuff in scope-full with the source where it came from.
useful to who? reviewers that want to remove noise, and make sure we are scanning our own stuff.
source ideas:
"scope"
"subfinder.domain.'owasp.com'"
"httpx"


---

alerts:
informative: when the current-scope changes, show what's new.
warnings: when the current-issues change, show what's new. show also what is not new, and 
          in both cases provide the commands to disable the warning

a report should be in the format:

    new surface: 
    - example
    issues detected:
    - example. if false positive, add this line to ignore issues: ``
    old issues:
    - example: if false positive, add this line to ignore issues: ``



---

git scanning:

trufflehog
https://github.com/trufflesecurity/trufflehog

git scanning is unrealistic, mostly because if you know the repo you want to watch,
it's easier to add trufflehog pipelines inside of them.


---

## design - interface

The user interaction will happen trough configuration files, gitlab issues, and CI pipelines execution.
The configuration should be as simple and as intuitive as possible. To achieve this we will 
follow these paradigms:
- sensible defaults. The system should work reliably with minimal configurations.
- configuration, documentation and technical terminologies should be detached from the internal 
  logic and data structures. e.g: no git-style docs.
- everything exposed to the user will be saved as human-redable yaml, with intuitive namings.

## design - database

Can we store everything as yaml, even the internal state?
This is an interesting opportunity we have, since the state
will mostly be an extension of the scope data structure.

by saving everything in a file, we can rely on version control and existing diff tools to keep 
the surface easy to understand.

data size estimates:
assuming a large organization with automated prefixes,
we might be in the thousands or domains / thousands of IPs.
assuming an average of 30bytes per entry
    e.g:
    a long list of domains with an uuidv4 would have 40 char per entry on average.
    a long list of ips would have 15 chars per entry on average
that's just a few thousands lines, at around 30Kib. 

What about urls? 
urls are larger, we can assume 100bytes per entry.
we might get a combinatory explosions of urls from each domain/ip.
- assuming only 10 urls per entry, 1000domains, and 1000 ips:
we'll gett 10*1000*2 = 20.000 entries --> 2Mib total.
- assuming 100 urls per entry, we'll get 200.000 entries --> 20Mib.

In conclusion
In the worst case scenario, for a huge Org, our surface will be stored in a 20Mib file.
- we'll have to store it in git, push it and pull it from the CI.
- we'll have to process it in memory, performing several data copies.

considerations:

gitalb stops showing the diff preview after 200KB, 50.000 lines
For an average organization the surface file will be below the 200Kib limit.
github max file size is 100Mb. after that, a file push will fail.




## data flow | data transformation programming frameworks

https://github.com/trustmaster/goflow

ELT pipelines ?

## we can easily make our own:
https://retejs.org/examples

reactflow
    https://www.repree.net/

https://flume.dev/

litegraph

the graph should compile into something like:

exclusions = init()
scope = init()
scopeFull = init()

deep_copy(scope, to:PIPE)
extract_domain(PIPE.urls, to:PIPE.domains)
remove(exclusions, to:PIPE)




### pipeline interface ideas:

the pipeline could be composed of several operators, with data 
flowing trough them in a functional manner.

each operator could conform to some specific interface, and we could have a list 
of all the operators defined somewhere. 

if the interface includes the []configfiles, and validateConfig(),
we could then print at startup the config files expected or the config files read.
and we could have a centralized way to validate config files at startup.


for an initial prototype pipeline operators can be regular functions, composed manually.
If we eventually want to implement a GUI pipeline compiler, we can 
define all operators as a struct with a execute function, all implementing the operator interface.
the execution function wil ltake as imput the pipeline memory map, and maybe some 
indexes into that map for input output. The details are not important now, 
just the general idea of having a dynamic runtime with typed data.


some of the operators could be designed to print the results, and the count of results.
The web interface could be used in debug mode to start actual scans, and in that scenario 
the output of each module would be useful


### going out of scope

the risk of going out of scope is mitigated by the fact that we don't run invasive scans.
still, it would be nice if we could add updated lists of known companies to the exclusions.
e.g: github
