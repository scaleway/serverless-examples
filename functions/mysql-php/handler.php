<?php
function insertRecord($mysqli, $tableName, $data) {
    // Check if mysqli object is null
    if ($mysqli === null) {
        die("Database connection is not established.");
    }

    // Create table if it doesn't exist
    $createTableSql = "CREATE TABLE IF NOT EXISTS `$tableName` (
        id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
        email VARCHAR(50),
        reg_date TIMESTAMP
    )";

    if ($mysqli->query($createTableSql) === FALSE) {
        die("Error creating table: " . $mysqli->error);
    }

    // Prepare SQL statement
    $stmt = $mysqli->prepare("INSERT INTO $tableName (email, reg_date) VALUES (?, NOW())");
    if ($stmt === false) {
        die("Error preparing statement: " . $mysqli->error);
    }

    // Bind parameters
    $stmt->bind_param("s", $data['email']);

    // Execute statement
    if ($stmt->execute() === false) {
        die("Error executing statement: " . $stmt->error);
    }

    echo "Record inserted successfully";
    $stmt->close();
}

function handle($event, $context) {
    $dbHost = getenv('MYSQL_HOST');
    $dbUsername = getenv('MYSQL_USER');
    $dbPassword = getenv('MYSQL_PASSWORD');
    $dbName = getenv('MYSQL_DB');

    $mysqli = new mysqli($dbHost, $dbUsername, $dbPassword, $dbName);
    $tableName = "users";

    if ($mysqli === null) {
        die("Database connection is not established.");
    }

    if ($mysqli->connect_error) {
        die("Connection failed: " . $mysqli->connect_error);
    }

    // Extract user data from the request
    $data = json_decode($event['body'], true);
    if ($data === null || $data['email'] === null) {
        return [
            "body" => "Email is required",
            "statusCode" => 400,
        ];
    }

    try {
        insertRecord($mysqli, $tableName, $data);
    } catch (Exception $e) {
        return [
            "body" => "Error: " . $e->getMessage(),
            "statusCode" => 500,
        ];
    }
  
	return [
		"body" => "Inserted " . $data['email'] . " into the database",
		"statusCode" => 200,
	];
}
