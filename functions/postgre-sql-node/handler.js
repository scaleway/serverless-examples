const pg = require('pg')

const PG_HOST=process.env.PG_HOST;
const PG_USER=process.env.PG_USER;
const PG_DATABASE=process.env.PG_DATABASE;
const PG_PASSWORD=process.env.PG_PASSWORD;
const PG_PORT=parseInt(process.env.PG_PORT,10);

const pool = new pg.Pool({
  user: PG_USER,
  host: PG_HOST,
  database: PG_DATABASE,
  password: PG_PASSWORD,
  port: PG_PORT
})

exports.handle = async (event, context, callback) => {
  try {
      const { rows } = await query("SELECT *")
      console.log(JSON.stringify(rows[0]))
      const response = {
          "statusCode": 200,
          "headers": {
              "Content-Type" : "application/json"
          },
          "body": JSON.stringify(rows),
          "isBase64Encoded": false
      };
      callback(null, response);
  } catch (err) {
      console.log('Database ' + err)
      callback(null, 'Database ' + err);
  }
  };

async function query (q) {
  const client = await pool.connect()
  let res
  try {
    await client.query('BEGIN')
    try {
      res = await client.query(q)
      await client.query('COMMIT')
    } catch (err) {
      await client.query('ROLLBACK')
      throw err
    }
  } finally {
    client.release()
  }

  return res
}