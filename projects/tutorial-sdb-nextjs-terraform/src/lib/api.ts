'use server'
import { Post } from "@/interfaces/post";
import fs from "fs";
import matter from "gray-matter";
import { join } from "path";
import { Client } from 'pg'
import { getPostFromRow } from "./utils"

const client = new Client({
    ssl: {
      rejectUnauthorized: false,
    }
  })

const connect = async () => {
  await client.connect()
}

connect()

export async function getPostBySlug(slug: string): Promise<Post> {
  const data = await client.query('SELECT * FROM posts WHERE slug=$1;',[slug])
  const post = getPostFromRow(data.rows[0])
  return post;
}

export async function getAllPosts(): Promise<Post[]> {
  const data = await client.query('SELECT * FROM posts;')
  const posts = data.rows.map((row) => (getPostFromRow(row)))
  return posts;
}

