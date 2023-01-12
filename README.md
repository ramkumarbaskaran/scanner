# Install:
1. Clone the application in your $workspace/github.com
2. Start Mysql Server ( I would choose Nosql for this type of application. To make it simple I have choosen Mysql here.)
3. Create Database and Table. Database migration script is in app/database folder.
4. $workspace/github.com/scanner/app/.env file have all the below environment variables. Please update the values which is specific to your environment.
```
    SERVER_PORT=8080
    DB_DRIVER=mysql
    DB_USER=
    DB_PASSWORD=
    DB_HOST=localhost
    DB_PORT=3306
    DB_DATABASE=
    SEARCH_PATTERN=public_key,prefix_key
    SCANClONEFOLDER=/tmp
    SCANClONEFOLDERPREFIX=repo
    NOOFWORKERS=3
```
# Test:
```
cd $workspace/github.com/scanner
make test 
```
# Build:
```
cd $workspace/github.com/scanner
make build
```
# Run:
```
cd $workspace/github.com/scanner
make run
```

# Endpoint test:

1.  Create new repo:
  ```sh
  curl -d "name=test&url=https://github.com/test/test"  -X  POST "http://localhost:8080/api/repo"
  ```

2. List all repos:
  ```sh
  curl -X GET "http://localhost:8080/api/repos"
   ```

3. Get repo by ID:
  ```sh
   curl -X GET "http://localhost:8080/api/repo/{repoID}"
   ```

4. Update repo:
  ```sh
   curl -X PUT "http://localhost:8080/api/repo/{repoID}" -d "name=test&url=https://github.com/test/test"
   ```

5. Delete repo:
  ```sh
   curl -X DELETE "http://localhost:8080/api/repo/{repoID}"
   ```

6. Scan repo:
  ```sh
   curl -X  POST "http://localhost:8080/api/repo/{repoID}/scan"
   ```

7. List all results:
  ```sh
   curl -X  GET "http://localhost:8080/api/scan/results"
   ```

8. Get result by ID:
  ```sh
   curl -X  GET "http://localhost:8080/api/scan/result/{resultID}"
   ```

 Note:  Replace host and port number with your host and port.

# Architecture:
![architecture](https://user-images.githubusercontent.com/3071990/211971856-1b787448-8326-4dcd-b40a-2c6ab47ce141.jpeg)
<code>
We can add more instances when required. We can use monitoring tools to check memory consumption, CPU utilisations  and disk space. If it reaches the threshold limit we can create an automation script to spin up more instances.
</code>


