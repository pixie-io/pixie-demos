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


## Attack the Vulnerable Application using  sqlmap


## Capture the SQL Injections using PxL
1. Browse to your vulnerable cluster on `https://work.withpixie.ai/`.
1. Under the script drop down select `Scratch Pad`.
1. Replace the PxL Script contents with the contents of `scripts/sql_injections.pxl`.
1. Replace the Vis Spec contents with the contents of `scripts/vis.json`.
1. Click Run.
1. You should now see the SQL Injection queries run by sqlmap in the data table.
