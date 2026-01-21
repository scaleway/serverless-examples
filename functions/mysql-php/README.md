# MySQL RDB and PHP functions

A small example of a PHP function that creates a table in a MySQL database and inserts a row into it.

## Deploying

1. Create a MySQL database from the Scaleway console.
   1. Make sure to create a user as well.
2. Create a new function namespace and a PHP function running using PHP 8.2
3. Copy the contents of `handler.php` into the function code editor.
4. Add the following environment variables to the function:
   - `MYSQL_HOST`: The hostname of the MySQL database
   - `MYSQL_USER`: The username of the MySQL user
   - `MYSQL_PASSWORD`: The password of the MySQL user
   - `MYSQL_DB`: The name of the MySQL database
5. Deploy the function.
6. Grab the function URL and use it to call the function.

## Usage

Call the function with a body to insert a row into the table.

```console
curl -X POST -d '{"email": "hello@acme.org"}' <function-url>
Inserted hello@acme.org into the database
```
