import mongoose from "mongoose";
import { randomIntegerBetween } from "@std/random";

const MONGODB_CERT_PATH = "/tmp/mongodb_cert.pem";

const MONGODB_URI = Deno.env.get("MONGODB_URI")
const MONGODB_CERT = Deno.env.get("MONGODB_CERT")

let mongooseConnected = false;

async function getMongooseConnection() {
  if (!mongooseConnected) {
    await mongoose.connect(MONGODB_URI || "", {
      tls: true,
      tlsCAFile: MONGODB_CERT_PATH,
    } as mongoose.ConnectOptions);

    mongooseConnected = true;
  }
}

const CHECK_CONNECTION = new URLPattern({ pathname: "/check_connection" });
const CREATE_PERSON = new URLPattern({ pathname: "/person/:name" });
const LIST_PEOPLE = new URLPattern({ pathname: "/people" });

async function handler(req: Request): Promise<Response> {
  const url = new URL(req.url);

  if (CHECK_CONNECTION.test(url)) {
    return await check_connection(req);
  } else if (CREATE_PERSON.test(url)) {
    return await create_person(req);
  } else if (LIST_PEOPLE.test(url)) {
    return await list_people(req);
  } else {
    return new Response("Not Found", { status: 404 });
  }
}

const peopleSchema = new mongoose.Schema({ name: String, age: Number, });
const ExamplePeople = mongoose.model('Example', peopleSchema);

async function create_person(req: Request): Promise<Response> {
  await getMongooseConnection();

  const url = new URL(req.url)
  const name = url.pathname.split("/person/")[1];

  const examplePerson = new ExamplePeople({
    name: name,
    age: randomIntegerBetween(1, 100),
  });
  await examplePerson.save();

  return new Response(`Created person: ${examplePerson.name}`, { status: 201 });
}

async function list_people(_req: Request): Promise<Response> {
  await getMongooseConnection();

  const people = await ExamplePeople.find().exec();
  const peopleNames = people.map(person => person.name).join(", ");

  return new Response(`People: ${peopleNames}`, { status: 200 });
}

async function check_connection(_req: Request): Promise<Response> {
  try {
    await getMongooseConnection();

    const connectionState = mongoose.connection.readyState;
    const body = `Mongoose connection state: ${connectionState}`;

    return new Response(body, { status: 200 });
  } catch (error) {
    return new Response(`MongoDB connection failed: ${error}`, { status: 500 });
  }
}

if (import.meta.main) {
  if (MONGODB_CERT) {
    await Deno.writeTextFile(MONGODB_CERT_PATH, MONGODB_CERT);
  }

  Deno.serve({ port: 8080, hostname: "0.0.0.0" }, handler);
}
