# litetable-cli
The Litetable CLI is a command-line utility for easily working and prototyping with any Litetable server.
**Litetable CLI is in active development**

<img src="./images/litetable-logo-min.png" alt="LiteTable Logo" width="350px">



1. Install the latest CLI:
   ```bash
   curl -fsSL https://raw.githubusercontent.com/litetable/litetable-cli/main/install.sh | bash
   ``` 
   or

   ```bash
   go get github.com/litetable/litetable-cli
   ```
2. Initialize a new LiteTable database:
   ```bash
    litetable service init
    ```

3. Start the LiteTable server:
   ```bash
   litetable service start
   ```

4. Stop the LiteTable server:
   ```bash
   litetable service stop
   ```

With an initialized server, you can start writing data to it. The first write is to always
create a supported column family, which is accomplished by a `create` command.
```bash
litetable create --family <my_family>
```

A valid column family is required for every read and write command.

### Create some data to your column family:
1. With a running server, create a new column family:
   ```bash
   litetable create --family wrestlers
   ```

2. Create a new record for that column family
   ```bash
      litetable write -k champ:1 -f wrestlers -q firstName -v John -q lastName -v Cena -q  championships -v 15
      ```
3. Append more data to the row key
   ```bash
      litetable write -k champ:1 -f champions -q championships -v 16 &&
      litetable write -k champ:1 -f champions -q championships -v 17
      ```
4. Read the data back
   ```bash
   litetable read -k champ:1 -f wrestlers
   ```

5. Delete a column qualifier
   ```bash
   litetable delete -k champ:1 -f wrestlers -q championships
   ```

6. Delete with custom TTL (number of seconds before garbage collection)
   ```bash
   litetable delete -k champ:1 -f wrestlers -q championships --ttl 300
   ```
