# albumCQL
Simple Web API to do CRUD operations using Go and GoCQL (Cassandra DB)

PRE-REQUISITES: (PLEASE CLICK ON "RAW" VERSION TO SEE CLEARLY.)
===============
1. Implement Cassandra DB. See this link on how to install Cassandra DB on your Windows machine. (https://www.youtube.com/watch?v=EEXtVn3zAqc)

2. Make sure you are importing GoCQL library in your Go program.(import "github.com/gocql/gocql")

3. Make sure Cassandra DB is running on your windows machine as described in Step 1. Check it via "cqlsh".

4. Open "cqlsh" and create keyspace as mentioned in your golang program. In our case, we have "albumspace" as keyspace name in our main.go file.

Execute below:
--------------
cqlsh> DESCRIBE keyspaces
system_auth  system_distributed  system_traces
system_schema  system 
cqlsh> CREATE KEYSPACE albumspace WITH replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
cqlsh> DESCRIBE keyspaces
albumspace        system_auth  system_distributed  system_traces
system_schema  system 
cqlsh> 

5. In "cqlsh" itself, go inside your keyspace and create your DB Table. In our case, our table name is "albumtable" as mentioned in main.go file.

Execute below: 
--------------
cqlsh> USE albumspace;
cqlsh:albumspace> CREATE TABLE albumtable (albname TEXT PRIMARY KEY, imagelist LIST<TEXT>);
cqlsh:albumspace> SELECT * FROM albumtable;

 albname | imagelist
---------+-----------

(0 rows)
cqlsh:albumspace> 
6. Execute main.go
