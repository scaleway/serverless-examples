
export function getPostFromRow(row: any){

    const post = {
      title: row.title,
      excerpt: row.excerpt,
      coverImage: row.coverimage,
      date: row.date.toISOString(),
      author: {
        name: row.author_name,
        picture: row.author_picture 
      },
      ogImage: {
        url: row.ogimage_url
      },
      slug: row.slug,
      content: row.content,
      preview: false
    }
  return post;
}
