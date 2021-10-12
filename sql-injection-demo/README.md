# SQL Injection Demo
This is a demo of how Pixie can be used to capture SQL injections on a Kubernetes
application. In this demo, we will spin up a
[DVWA web application](https://hub.docker.com/r/vulnerables/web-dvwa) that is vulnerable
to SQL injection monitored by Pixie, run
[SQLMap](https://github.com/SQLMapproject/SQLMap) (a sql injection tool) against that
application, and detect the SQL injections at the database level using a PxL script.

## WARNING!
DVWA is an *intentionally* vulnerable web application. It should **NOT** be deployed to
a live web server. These instructions will cover deploying DVWA to a minikube
environment. See DVWA's [disclaimer](https://github.com/digininja/DVWA) for more
details.

## Create your Minikube cluster
1. Install [Minikube](https://minikube.sigs.k8s.io/docs/start/)
1. Run minikube. Linux users should use the kvm2driver and Mac users should use the
[hyperkit](https://minikube.sigs.k8s.io/docs/drivers/hyperkit/) driver. Other drivers,
including the docker driver, are not supported by Pixie.
    ```
    minikube start --driver=<kvm2|hyperkit> --cni=flannel --cpus=4 --memory=8000 -p=<cluster-name>
    ```
1. Verify your cluster is up and running.
    ```
    kubectl get nodes
    ```

## Deploy Pixie to your cluster
1. Follow an install guide to
[Install Pixie](https://docs.px.dev/installing-pixie/install-guides).

## Deploy the Vulnerable Application
**WARNING** This image is vulnerable to several kinds of attacks. You should only deploy
it to your `minikube` cluster.
1. Ensure that you are still running on your `minikube` environment.
    ```
    kubectl config current-context
    ```
1. `git clone` this repo and `cd` into the `sql-injection-demo` directory.
    ```
    git clone <path to repo>
    cd <repo_path>/sql-injection-demo
    ```
1. Deploy the vulnerable demo application.
    ```
    kubectl apply -f ./dvwa-k8s
    ```
1. Get the pod name for dvwa-pixie-demo.
    ```
    kubectl get pods
    ```
1. Forward the port so you can access the UI. Leave this running. You can use a
different value for 1234 if you want, just make sure you replace it in subsequent
instructions.
    ```
    kubectl port-forward <dvwa-podname> 1234:80
    ```    
1. In your browser, navigate to `localhost:1234`
1. Login with username: `admin`, password: `password`
1. Follow instructions on webpage and click `Create / Reset Database` 
1. Relogin with username: `admin`, password: `password`

## Manual SQL Injection
DVWA was designed with a SQL injection that originates from taking raw user input in a
URL query parameter. The path of the vulnerability is
`http://localhost:1234/vulnerabilities/sqli/?id=<SQL-Injection-Point>&Submit=Submit#`.
An attacker could supply a crafted value for the ID query parameter which ultimately
would lead to a SQL injection. 

At the database level, the raw query will look like:
`SELECT First_Name,Last_Name FROM users WHERE ID=<SQL-Injection-Point>;`


As an example, you can view try `1' union select 1,@@version#` as the `id` value. This
will append the database version to the results by including a `union select` injection. 

Try accessing the following URL:
```
http://localhost:1234/vulnerabilities/sqli/?id=1%27+union+select+1%2C%40%40version%23%26Submit%3DSubmit%23&Submit=Submit#`
```

## Automating finding SQL injections with SQLMap
[SQLMap](https://github.com/SQLMapproject/SQLMap) is a CLI tool that automates finding
SQL injections via bruteforce and heuristic methods.

1. In a new tab, `git clone` the [SQLMap](https://github.com/SQLMapproject/SQLMap) repo
based on their instructions in the README and `cd` into it.

1. To use SQLMap you will need the cookie from DVWA after logging in. Go to
`http://localhost:1234/phpinfo.php` and scroll down to `PHP Variables`. 

1. Copy the cookie inside the entry `$_COOKIE['PHPSESSID']`. and export it into an
environment variable. 
    ```
    export DVWA_COOKIE='PHPSESSID=<YOUR-PHP-SESSID>; security=low'
    ```
1. Run SQLMap.
    ```
    python sqlmap.py -u 'http://localhost:1234/vulnerabilities/sqli/?id=1&Submit=Submit#' -cookie $DVWA_COOKIE
    ```
1. SQLMap will prompt for answers to various questions as it runs. Answer the following
when prompted:
    * *it looks like the back-end DBMS is 'MySQL'. Do you want to skip test payloads specific for other DBMSes?* y
    * *for the remaining tests, do you want to include all tests for 'MySQL' extending provided level (1) and risk (1) values?* n
    * *found a vuln, asking do you want to keep testing* n

## Capture the SQL Injections using PxL
1. Execute a script via the Pixie CLI. This script returns the latest MySQL queries
Pixie observed on your cluster.
    ```
    px run px/mysql_data
    ```
1. Load the above view in the Pixie UI. In your browser, navigate to the URL printed at
the bottom of the CLI output at `Live UI:`.
1. Open the script editor (Ctrl + E).
1. Replace the PxL Script tab contents with the contents of `script/sql_injections.pxl`.
1. Replace the Vis Spec tab contents with the contents of `script/vis.json`.
1. Click Run.
1. You should now see the SQL Injection queries run by SQLMap in the data table.

## Clean up
1. Delete your minikube cluster.
    ```
    minikube delete
    ```
