# SQL Injection Demo
This is a demo of how Pixie can be used to capture SQL injections on an application. In
This demo, we will spin up a
[DVWA web application](https://hub.docker.com/r/vulnerables/web-dvwa) that is vulnerable
to SQL injection monitored by Pixie, run
[sqlmap](https://github.com/sqlmapproject/sqlmap) (a sql injection tool) against that
application, and detect the SQL injections at the database level  using a PxL script.


## Prerequisites
* [Kubernetes](https://kubernetes.io/docs/tasks/tools/) to deploy the vulnerable web
application monitored by Pixie.
* A Pixie account. Follow these instructions here.

## Deploy the Vulnerable Application
* Deploy demo application `kubectl apply -f ./dvwa`
* Login with username: `admin` , password: `password`
* Follow instructions on webpage and click `Create / Reset Database` 
* Relogin with username: `admin` , password: `password`

## SQL Injection
DVWA was designed with an SQL injection that originates from taking raw user input in a URL query parameter. The path of the vulernability is `http://____domain____/vulnerabilities/sqli/?id=<SQL-Injection-Point>&Submit=Submit#`. An attacker could supply a crafted value for the ID query param which ultimately would lead to an SQL injection. 

At the database level, the raw query will look like:
`SELECT First_Name,Last_Name FROM users WHERE ID=<SQL-Injection-Point>;`


As an example, you can view try `1' union select 1,@@version#` as the `id` value. This will append the database version to the results by including a `union select` injection. 
Example: `http://____domain____/vulnerabilities/sqli/?id=1%27+union+select+1%2C%40%40version%23%26Submit%3DSubmit%23&Submit=Submit#`

## Automating finding SQL injections with SQLMap
SQLMap is a CLI tool that automates finding SQL injections via bruteforce and huerstic methods. https://github.com/sqlmapproject/sqlmap . To use SQLMap you will need the cookie from DVWA after logging in, please copy the PHPSESSID value. 

```
export DVWA_COOKIE='PHPSESSID=<YOUR-PHP-SESSID>; security=low'
sqlmap -u 'http://____domain____/vulnerabilities/sqli/?id=1&Submit=Submit#' -cookie $DVWA_COOKIE
```
SQLMap question answers: `Y`, `n`, `N`

*sqlmap question 1 is `Y`  - found a vuln, do you want to skip trying other DB types*

*sqlmap question 2 is `n`  - found a vuln, do you want to try all MySQL tests*

*sqlmap question 3 is `N`  - found a vuln, asking do you want to keep testing*

## Capture the SQL Injections using PxL
1. Browse to your vulnerable cluster on `https://work.withpixie.ai/`.
1. Under the script drop down select `Scratch Pad`.
1. Replace the PxL Script contents with the contents of `scripts/sql_injections.pxl`.
1. Replace the Vis Spec contents with the contents of `scripts/vis.json`.
1. Click Run.
1. You should now see the SQL Injection queries run by sqlmap in the data table.
