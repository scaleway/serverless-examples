import { Metadata } from "next";
import { notFound } from "next/navigation";
import { getAllPosts, getPostBySlug } from "../../../lib/api";
import { CMS_NAME } from "../../../lib/constants";
import markdownToHtml from "../../../lib/markdownToHtml";
import Alert from "../../_components/alert";
import Container from "../../_components/container";
import Header from "../../_components/header";
import { PostBody } from "../../_components/post-body";
import { PostHeader } from "../../_components/post-header";


export default async function Post({ params }: Params) {
  const post = await getPostBySlug(params.slug)

  if (!post) {
    return notFound();
  }

  return (
    <main>
      <Alert preview={post.preview} />
      <Container>
        <Header />
        <article className="mb-32">
          <PostHeader
            title={post.title}
            coverImage={post.coverImage}
            date={post.date}
            author={post.author}
          />
          <PostBody content={post.content} />
        </article>
      </Container>
    </main>
  );
}

type Params = {
  params: {
    slug: string;
  };
};

export async function generateMetadata({ params }: Params): Promise<Metadata> {
  const post = await getPostBySlug(params.slug);

  if (!post) {
    return notFound();
  }

  const title = `${post.title} | Next.js Blog Example with ${CMS_NAME}`;

  return {
    openGraph: {
      title,
      images: [post.ogImage.url],
    },
  };
}




