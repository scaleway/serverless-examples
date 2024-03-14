---
meta:
  title: Deploy your Full Serverless Next.js application using Serverless Containers and Serverless SQL Database
  description: This page explains how to deploy your full serverless Next.js application using Serverless Containers and Serverless SQL Database
content:
  h1: Deploy your Full Serverless Next.js application using Serverless Containers and Serverless SQL Database
  paragraph: This page explains how to deploy your full serverless Next.js application using Serverless Containers and Serverless SQL Database
dates:
  validation: 2024-02-08
  posted: 2024-02-08
---

This tutorial will let you deploy your full serverless Next.js application using Serverless Containers and Serverless SQL Database in a few minutes.

<Message type="requirement">
  - You have a **Scaleway account** and can log into the [Scaleway console](https://console.scaleway.com/) 
  - You have installed the [Scaleway CLI](https://www.scaleway.com/en/docs/developer-tools/scaleway-cli/quickstart/) and initialized it with your [Scaleway credentials](https://www.scaleway.com/en/docs/developer-tools/scaleway-cli/quickstart/#how-to-configure-the-scaleway-cli)
  - You have **Docker Engine** installed, either through [Docker Desktop](https://docs.docker.com/engine/install/) (recommended for MacOS and Windows) or directly [from binaries](https://docs.docker.com/engine/install/binaries/)
  - You have the right [IAM permissions](https://www.scaleway.com/en/docs/identity-and-access-management/iam/reference-content/policy/) to push containers image to Container Registry and to create Serverless Containers. You specifically require **ContainerRegistryFullAccess**, **ContainersFullAccess** and **ServerlessSQLDatabaseFullAccess** permission sets.
  - You have [Node >=18.0.0](https://nodejs.org/en/download) installed.
</Message>

<Message type="tip">
  If you encounter any unexpected error, please check beforehand that **you meet all requirements** mentioned above. You can find **safety checks and commands** that you can run to make sure you are perfectly setup in the [Troubleshooting](troubleshooting) section.
</Message>

To begin with either start by:
- [Deploying your application step by step using Scaleway CLI](#deploy-next.js-application-frontend-using-serverless-containers) to understand each detailed actions performed and subresources required
- [Deploy your application using a Terraform template](#deploy-next.js-application-using-terraform-templates) to deploy your application faster and have a ready to use Infrastructure as Code. This part requires you to have [Terraform](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli) installed and configured with [Scaleway Terraform Provider](https://registry.terraform.io/providers/scaleway/scaleway/latest/docs).

## Deploy Next.js application frontend using Serverless Containers

1. First, make sure your local environment variables are properly setup:
   ```bash
   scw info 
   ```
   This command should display your access_key and secret_key (first characters only) as the last two lines. Make sure the ORIGIN column displays `env (SCW_ACCESS_KEY)` and `env (SCW_SECRET_KEY)`, and not `default profile`.
   ```bash
   KEY                      VALUE                                       ORIGIN
   (...)
   access_key               SCWF9DETBF829TAFB1TB                        env (SCW_ACCESS_KEY)
   secret_key               9a1bcf92-xxxx-xxxx-xxxx-xxxxxxxxxxxx        env (SCW_SECRET_KEY)
   ```
   If this is not the case, add these environment variables to your local environment:
   ```bash
   export SCW_ACCESS_KEY=$(scw config get access-key)
   export SCW_SECRET_KEY=$(scw config get secret-key)
   ```
2. Go to the folder where you want to store your project repository and initialize a blog template application with the following command:
    ```bash
    npx create-next-app --example blog-starter my-nextjs-blog
    ```
    where `my-nextjs-blog` will be the newly created folder name and can be changed to your own project name.  
    Accept to install dependencies:
    ```bash
    Need to install the following packages:
    create-next-app@14.0.3
    Ok to proceed? (y) y
    ```
    This will create a folder `my-nextjs-blog`. Go into it with:
    ```bash
    cd my-nextjs-blog
    ```
3. Confirm you can run the application locally with the following commands:
    ```bash
    npm run dev
    ```
    Go to [http://localhost:3000](http://localhost:3000) in your browser. The blog template should display.

4. Create a `Dockerfile`:
    ```bash
    touch Dockerfile
    ```
    Add the following content to it:
    ```bash
    # syntax=docker/dockerfile:1

    FROM node:20-alpine
     
    #Built time arguments used for pre-rendering
    ARG PGHOST
    ARG PGPORT
    ARG PGDATABASE
    ARG PGUSER
    ARG PGPASSWORD

    WORKDIR /usr/app
    COPY ./ ./
    RUN npm install
    RUN npm run build
     
    #Web application configuration
    ENV PORT=8080
     
    #Database configuration used for dynamically rendered data. These default values should be overwritten by container runtime environment variables
    ENV PGHOST=localhost
    ENV PGPORT=5432
    ENV PGDATABASE=database
    ENV PGUSER=user
    ENV PGPASSWORD=password
     
    CMD npm run start
    ```
5. Create a `.dockerignore` file by copying `.gitignore` file:
    ```bash
    cp .gitignore .dockerignore
    ```
    This will let Docker ignore node dependencies and build files which will be generated during container build.
6. Build your application container with:
    ```bash
    docker build -t my-nextjs-blog .   
    ```

    You can check that your container is running correctly locally with:
    ```bash
    docker run -it -p 8080:8080 my-nextjs-blog
    ```
    Go to http://localhost:8080 in your browser. The blog template should display.
    <Message type="tip">
    When connecting to the webpage, your terminal running docker might display the following error `Error: ENOENT: no such file or directory, open '/usr/app/_posts/%5Bslug%5D.md'`. This will be fixed in further steps.
    </Message> 

7. Create a [Container Registry Namespace](https://www.scaleway.com/en/docs/containers/container-registry/concepts/#namespace):
    ```bash
    export REGISTRY_ENDPOINT=$(scw registry namespace create -o json | jq -r '.endpoint')
    ```
8. Push your containerized application to [Container Registry](https://www.scaleway.com/en/docs/containers/container-registry/quickstart/).  
    First login to **Container Registry** from your local terminal:
    ```bash
    docker login $REGISTRY_ENDPOINT -u nologin --password-stdin <<< "$SCW_SECRET_KEY"
    ```

    Then, tag and push your container image to **Container Registry**:
    ```bash
    docker tag my-nextjs-blog:latest $REGISTRY_ENDPOINT/my-nextjs-blog:latest
    docker push $REGISTRY_ENDPOINT/my-nextjs-blog:latest
    ```
9. Deploy your container to [Serverless Containers](https://www.scaleway.com/en/docs/serverless/containers/quickstart/).  
   First, create a [Serverless Containers Namespace](https://www.scaleway.com/en/docs/serverless/containers/concepts/#namespace):
   ```bash
   export CONTAINER_NAMESPACE_ID=$(scw container namespace create name="my-nextjs-blog-ns" -o json | jq -r '.id')
   ```
   Deploy your application to **Serverless Containers**:
   ```bash
   scw container container create name="my-nextjs-blog" namespace-id=$CONTAINER_NAMESPACE_ID registry-image=$REGISTRY_ENDPOINT/my-nextjs-blog:latest
   ```
   where `my-nextjs-blog` can be changed to any name you want to give to your running container.  

   Go to the provided url in your browser displayed next to `DomainName` property. The blog template should display.
   <Message type="tip">
   The first deployment can take up to one or two minutes. You can check the deployment status with the following command `scw container container list name=my-nextjs-blog `. As soon as the status is set as `ready`, you should be able to access the website in your browser.
   </Message>

10. Congratulations, your Next.js application frontend is **now online**!
   Your application cannot store any data persistently yet, so you will add serverless storage in the next section.


## Add serverless storage to your application using Serverless SQL Database

1. Create a [Serverless SQL Database](https://www.scaleway.com/en/docs/serverless/sql-databases/reference-content/serverless-sql-databases-overview/) using Scaleway CLI:
   ```bash
   export PGHOST=$(scw sdb-sql database create name=tutorial-nextjs-blog-db cpu-min=0 cpu-max=4 -o json | jq -r '.endpoint' | cut -d "/" -f3 | cut -d ":" -f1 )
   export PGPORT='5432'
   export PGDATABASE='tutorial-nextjs-blog-db'
   export PGUSER=$(scw iam api-key get $SCW_ACCESS_KEY -o json | jq -r '.user_id')
   export PGPASSWORD=$SCW_SECRET_KEY
   ```
   This commands create a Serverless SQL Database and adds its configuration information into your default PostgreSQL local environment variables. These environment variables are automatically recognized by most PostgreSQL compatible tool.
2. Connect to your database and add sample data using `psql` tool:
   ```bash
   psql 
   ```
   You don't have to provide any arguments as psql uses automically `PGHOST`, `PGPORT`, `PGDATABASE`, `PGUSER`, `PGPASSWORD` environment variables as default arguments.  
   You should now see an input field with the name of your database:
   ```bash
   psql (15.3, server 16.1 (Debian 16.1-1.pgdg120+1))
   SSL connection (protocol: TLSv1.3, cipher: TLS_AES_128_GCM_SHA256, compression: off)
   Type "help" for help. 

   tutorial-nextjs=> 
   ```
   Create a table structure with the following query: 
   ```sql
   create table posts (title char(100), excerpt text, coverImage text, date date, author_name char(50), author_picture text, ogImage_url text, slug char(50), content text);
   ```
   Then, add data using the following command:
   ```sql
   insert into posts values 
     ('Learn How to Pre-render Pages Using Static Generation with Next.js', 
     'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Praesent elementum facilisis leo vel fringilla est ullamcorper eget. At imperdiet dui accumsan sit amet nulla facilities morbi tempus.',
     '/assets/blog/hello-world/cover.jpg',
     '2020-03-16T05:35:07.322Z',
     'Tim Neutkens',
     '/assets/blog/authors/tim.jpeg',
     '/assets/blog/hello-world/cover.jpg',
     'hello-world',
     'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Praesent elementum facilisis leo vel fringilla est ullamcorper eget. At imperdiet dui accumsan sit amet nulla facilities morbi tempus. Praesent elementum facilisis leo vel fringilla. Congue mauris rhoncus aenean vel. Egestas sed tempus urna et pharetra pharetra massa massa ultricies.'),
     ('Dynamic Routing and Static Generation', 
     'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Praesent elementum facilisis leo vel fringilla est ullamcorper eget. At imperdiet dui accumsan sit amet nulla facilities morbi tempus.',
     '/assets/blog/dynamic-routing/cover.jpg',
     '2020-03-16T05:35:07.322Z',
     'JJ Kasper',
     '/assets/blog/authors/jj.jpeg',
     '/assets/blog/dynamic-routing/cover.jpg',
     'dynamic-routing',
     'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Praesent elementum facilisis leo vel fringilla est ullamcorper eget. At imperdiet dui accumsan sit amet nulla facilities morbi tempus. Praesent elementum facilisis leo vel fringilla. Congue mauris rhoncus aenean vel. Egestas sed tempus urna et pharetra pharetra massa massa ultricies.'),
     ('Preview Mode for Static Generation', 
     'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Praesent elementum facilisis leo vel fringilla est ullamcorper eget. At imperdiet dui accumsan sit amet nulla facilities morbi tempus.',
     '/assets/blog/preview/cover.jpg',
     '2020-03-16T05:35:07.322Z',
     'Joe Haddad',
     '/assets/blog/authors/joe.jpeg',
     '/assets/blog/preview/cover.jpg',
     'preview',
     'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Praesent elementum facilisis leo vel fringilla est ullamcorper eget. At imperdiet dui accumsan sit amet nulla facilities morbi tempus. Praesent elementum facilisis leo vel fringilla. Congue mauris rhoncus aenean vel. Egestas sed tempus urna et pharetra pharetra massa massa ultricies.');
   ```
   Exit database query interface by typing:
   ```bash
   exit
   ```
3. Add `node-postgres` dependency to your application:
   ```bash
   npm install pg && npm install --save-dev @types/pg
   ```
   This library will be needed to connect from your frontend to your database.
4. Your application folder should be called `my-nextjs-blog` and structured such as:
   ```bash
   my-nextjs-blog/
       .next/
       app/
       node_modules/
       public/
       pages/
       app/
       package.json
       tsconfig.json
   ```
   You need to edit 3 files to have a proper integration with your database.
   <Message type="tip">
    If you want to avoid doing these modifications manually, you can also directly clone the final code from this tutorial from a remote repository with `git clone git@gitlab.infra.online.net:fpagny/tutorial-sdb-nextjs-terraform.git`. Note that this repository differs slightly in configuration as it will dynamically render pages and not pre-render them in the Nextjs app build phase.
   </Message> 

   - First, edit the file located at `src/lib/api.ts` by replacing its content with the following code. This will switch frontend content from local data to your **Serverless SQL Database**.
   ```typescript
   import { Post } from "@/interfaces/post";
   import fs from "fs";
   import matter from "gray-matter";
   import { join } from "path";
   import { Client } from 'pg'
 
 
   export type Row = {
     slug: string;
     title: string;
     date: Date;
     coverimage: string;
     author_name: string;
     author_picture: string;
     excerpt: string;
     ogimage_url: string;
     content: string;
     preview?: boolean;
   };
 
   const client = new Client({
       ssl: {
         rejectUnauthorized: false,
       }
     })
 
   const connect = async () => {
     await client.connect()
   }
 
   connect()
 
   //This default function handles Nextjs pre-rendering with undefined values
   export function getPostFromRow(row: Row){
     if (row === undefined){ 
       var postRow: Row = {
         slug: "Title",
         title: "Excerpt",
         date: new Date("2024-01-16"),
         coverimage: "",
         author_name: "",
         author_picture: "",
         excerpt: "",
         ogimage_url: "",
         content: "",
         preview: false
       } 
     } else{
       var postRow = row
     }
 
     const post = {
       title: postRow.title,
       excerpt: postRow.excerpt,
       coverImage: postRow.coverimage,
       date: postRow.date.toISOString(),
       author: {
         name: postRow.author_name,
         picture: postRow.author_picture
       },
       ogImage: {
         url: postRow.ogimage_url
       },
       slug: postRow.slug,
       content: postRow.content,
       preview: false
     }
     return post;
   }
 
 
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
   ```
   - Secondly, edit the file located at `src/app/page.tsx` by replacing its content with the following code. This will update the main application home page that displays a posts list.
   ```typescript
   import Container from "@/app/_components/container";
   import { HeroPost } from "@/app/_components/hero-post";
   import { Intro } from "@/app/_components/intro";
   import { MoreStories } from "@/app/_components/more-stories";
   import { getAllPosts } from "../lib/api";
   
   export default async function Index() {
     const allPosts = await getAllPosts();
   
     const heroPost = allPosts[0];
   
     const morePosts = allPosts.slice(1);
   
     return (
       <main>
         <Container>
           <Intro />
           <HeroPost
             title={heroPost.title}
             coverImage={heroPost.coverImage}
             date={heroPost.date}
             author={heroPost.author}
             slug={heroPost.slug}
             excerpt={heroPost.excerpt}
           />
           {morePosts.length > 0 && <MoreStories posts={morePosts} />}
         </Container>
       </main>
     );
   }
   ```
   - Finally, edit the file located at `src/app/posts/[slug]/page.tsx` by replacing its content with the following code. This will update pages that displays a single post.
   ```typescript
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
   
   export async function generateStaticParams() {
     const posts = await getAllPosts();
   
     return posts.map((post) => ({
       slug: post.slug,
     }));
   }
   ```
5. Confirm you can run the application locally with the following command:
   ```bash
   npm run dev
   ```
   Go to [http://localhost:3000](http://localhost:3000) in your browser. The blog template should display with the updated data from your database (images and titles of the first blog post have changed!)
   
   **Congrats!** You could already deploy your application by:
   - Building a new container version with `docker build -t my-nextjs-blog .` as long as your also provide the right environment variables as build arguments
   - Push the new container to **Container Registry** with `docker tag my-nextjs-blog:latest $REGISTRY_ENDPOINT/my-nextjs-blog:latest` and `docker push $REGISTRY_ENDPOINT/my-nextjs-blog:latest`
   - Deploy your container to **Serverless Containers** with `scw container container deploy` as long as you provide the right environment variables to your container (as in step 2.). 
   However, your application would then connect with your [user credentials](https://www.scaleway.com/en/docs/identity-and-access-management/iam/concepts/#user), which is a bad security practice. To secure your deployment, you will now add a dedicated [IAM application](https://www.scaleway.com/en/docs/identity-and-access-management/iam/concepts/#application), give it the minimum required permissions, and provide its credentials to your application.

## Secure and redeploy your application

1. Create an [IAM application](https://www.scaleway.com/en/docs/identity-and-access-management/iam/concepts/#application): 
   ```bash
   export SCW_APPLICATION_ID=$(scw iam application create name=tutorial-nextjs-app -o json | jq -r '.id')
   ```
   The `SCW_APPLICATION_ID` environment variable will store the IAM appplication id so you can use it in later commands.

2. Create an [IAM policy](https://www.scaleway.com/en/docs/identity-and-access-management/iam/concepts/#policy) giving your application rights to access the database:
   ```bash
   scw iam policy create application-id=$SCW_APPLICATION_ID name=tutorial-nextjs-policy rules.0.organization-id=$(scw config get default-organization-id) rules.0.permission-set-names.0=ServerlessSQLDatabaseFullAccess
   ```
   The permission `ServerlessSQLDatabaseFullAccess` lets your application read and write data or create new databases. You can restrict this later to fit your database use cases, using `ServerlessSQLDatabaseRead` or `ServerlessSQLDatabaseReadWrite` permissions sets.

3. Create an [IAM Secret Key](https://www.scaleway.com/en/docs/identity-and-access-management/iam/concepts/#api-key):
   ```bash
   export SCW_APPLICATION_SECRET=$(scw iam api-key create application-id=$SCW_APPLICATION_ID -o json | jq -r '.secret_key')
   ```
4. Build your application container with the updated code:
    ```bash
    docker build -t my-nextjs-blog --build-arg PGHOST --build-arg PGDATABASE --build-arg PGPORT --build-arg PGUSER=$SCW_APPLICATION_ID --build-arg PGPASSWORD=$SCW_APPLICATION_SECRET . 
    ```
    This time build command requires environment variables because Nextjs pre-rendering will access your database at built time to load pages content.

    You can check that your container is running correctly locally with:
    ```bash
    docker run -it -p 8080:8080 -e PGHOST -e PGDATABASE -e PGPORT -e PGUSER=$SCW_APPLICATION_ID -e PGPASSWORD=$SCW_APPLICATION_SECRET my-nextjs-blog
    ```
    Go to http://localhost:8080 in your browser. The blog template should display.
5. Push your containerized application to [Container Registry](https://www.scaleway.com/en/docs/containers/container-registry/quickstart/):
    ```bash
    docker tag my-nextjs-blog:latest $REGISTRY_ENDPOINT/my-nextjs-blog:latest
    docker push $REGISTRY_ENDPOINT/my-nextjs-blog:latest
    ```
6. Redeploy your container to [Serverless Containers](https://www.scaleway.com/en/docs/serverless/containers/quickstart/). 
   First, get your container id with:
   ```bash
   export CONTAINER_ID=$(scw container container list name="my-nextjs-blog" -o json | jq -r '.[0].id')
   ```
   ```bash
   scw container container deploy $CONTAINER_ID
   ```
   Refresh your browser page displaying the blog. An updated version should display. 
7. **Congrats!** You have deployed your full serverless Nextjs application!
   
   To go further you can:
   - **Edit the source code** locally, apply `scw container deploy` again, and see the **new version being live** in seconds.
   - Inspect your **newly created resources** in the Scaleway Console
     - You can display your newly created **registry namespace** and **container image** in the [Container Registry Console Page](https://console.scaleway.com/registry/namespaces) 
     - You can display your newly created **serverless container namespace** and **container deployment** in the [Serverless Containers Console Page](https://console.scaleway.com/containers/namespaces)
     - You can display your newly created **serverless sql database** in the [Serverless SQL Database Console Page](https://console.scaleway.com/serverless-db/databases)
   - **Fine tune deployment options** such as autoscaling, targeted regions and much more. You can find more information by typing `scw container deploy --help` in your terminal.


## Deploy Next.js application using Terraform templates

1. Install and download the template source code from this repository:
    ```bash
    git clone git@gitlab.infra.online.net:fpagny/tutorial-sdb-nextjs-terraform.git my-nextjs-blog
    ```
    This repository extends the blog demo application from `npx create-next-app --example blog-starter my-nextjs-blog` and add required code to fetch data from a database instead of static markdown files.

    This will create a folder `my-nextjs-blog`. Go into it with:
    ```bash
    cd my-nextjs-blog
    ```
2. Install dependencies:
    ```bash
    npm install
    ```
3. Confirm you can run the application locally with the following commands:
    ```bash
    npm run dev
    ```
    Go to [http://localhost:3000](http://localhost:3000) in your browser. The blog template should display.
4. First, make sure your local environment variables are properly setup:
    ```bash
    scw info 
    ```
    This command should display your access_key and secret_key (first characters only) as the last two lines. Make sure the ORIGIN column displays `env (SCW_ACCESS_KEY)` and `env (SCW_SECRET_KEY)`, and not `default profile`.
    ```bash
    KEY                      VALUE                                       ORIGIN
    (...)
    access_key               SCWF9DETBF829TAFB1TB                        env (SCW_ACCESS_KEY)
    secret_key               9a1bcf92-xxxx-xxxx-xxxx-xxxxxxxxxxxx        env (SCW_SECRET_KEY)
    ```
    If this is not the case, add these environment variables to your local environment:
    ```bash
    export SCW_ACCESS_KEY=$(scw config get access-key)
    export SCW_SECRET_KEY=$(scw config get secret-key)
    ```
5. Build your application container with:
    ```bash
    docker build -t my-nextjs-blog .   
    ```
    <Message type="tip">
      An error message `Error: connect ECONNREFUSED` might display at build time. This is linked to Nextjs prerendering stage without any database to connect to yet, but will raise any issue at runtime because all pages will be dynamically rendered. 
    </Message>

    You can check that your container is running correctly locally with:
    ```bash
    docker run -it -p 8080:8080 my-nextjs-blog
    ```
    Go to http://localhost:8080 in your browser. You should see the webpage but with an error message `Application error: a server-side exception has occurred (see the server logs for more information)`. This is expected as your application is running but not plugged to a database yet.
6. Create a [Container Registry Namespace](https://www.scaleway.com/en/docs/containers/container-registry/concepts/#namespace):
    ```bash
    export REGISTRY_ENDPOINT=$(scw registry namespace create -o json | jq -r '.endpoint')
    ```
7. Push your containerized application to [Container Registry](https://www.scaleway.com/en/docs/containers/container-registry/quickstart/).  
    First login to **Container Registry** from your local terminal:
    ```bash
    docker login $REGISTRY_ENDPOINT -u nologin --password-stdin <<< "$SCW_SECRET_KEY"
    ```
    Then, tag and push your container image to **Container Registry**:
    ```bash
    docker tag my-nextjs-blog:latest $REGISTRY_ENDPOINT/my-nextjs-blog:latest
    docker push $REGISTRY_ENDPOINT/my-nextjs-blog:latest
    ```
8. Create a folder to store your terraform files in a dedicated repository:
    ```bash
    cd ..
    mkdir terraform-nextjs-blog
    ``` 

    Go into the terraform directory and create a `main.tf` terraform file:
    ```bash
    cd terraform-nextjs-blog
    touch main.tf
    ``` 

    Your application folder should now be structured as:
    ```bash
    my-nextjs-blog/
    terraform-nextjs-blog/
        main.tf
    ```
9. Add the following content to your `main.tf` file:
    ```bash
      terraform {
      required_providers {
        scaleway = {
          source = "scaleway/scaleway"
        }
      }
      required_version = ">= 0.13"
      }
       
      variable "REGISTRY_ENDPOINT" {
      type = string
      description = "Container Registry endpoint where your application container is stored"
      }
       
      variable "DEFAULT_PROJECT_ID" {
      type = string
      description = "Project id where your resources will be created"
      }
       
      resource scaleway_container_namespace main {
        name = "tutorial-nextjs-blog-tf"
        description = "Namespace created for full serverless Nextjs app deployment"
      }
       
      resource scaleway_container main {
        name = "tutorial-nextjs-blog-tf"
        description = "Container for Nextjs blog"
        namespace_id = scaleway_container_namespace.main.id
        registry_image = "${var.REGISTRY_ENDPOINT}/my-nextjs-blog:latest"
        port = 8080
        cpu_limit = 560
        memory_limit = 1024
        min_scale = 0
        max_scale = 5
        timeout = 600
        max_concurrency = 80
        privacy = "public"
        protocol = "http1"
        deploy = true
         
        environment_variables = {
            "PGUSER" = scaleway_iam_application.app.id,
            "PGHOST" = trimsuffix(trimprefix(regex(":\\/\\/.*:",scaleway_sdb_sql_database.database.endpoint), "://"),":")
            "PGDATABASE" = scaleway_sdb_sql_database.database.name,
            "PGPORT" = trimprefix(regex(":[0-9]{1,5}",scaleway_sdb_sql_database.database.endpoint), ":")
        }
        secret_environment_variables = {
          "PGPASSWORD" = scaleway_iam_api_key.api_key.secret_key,
        }
      }
       
      resource scaleway_iam_application "app" {
      name = "tutorial-nextjs-app-tf"
      }
       
      resource scaleway_iam_policy "db_access" {
      name = "tutorial-nextjs-policy-tf"
      description = "Gives tutorial Nextjs app access to Serverless SQL Database"
      application_id = scaleway_iam_application.app.id
      rule {
        project_ids = ["${var.DEFAULT_PROJECT_ID}"]
        permission_set_names = ["ServerlessSQLDatabaseReadWrite"]
      }
      }
       
      resource scaleway_iam_api_key "api_key" {
      application_id = scaleway_iam_application.app.id
      }
       
      resource scaleway_sdb_sql_database "database" {
      name = "tutorial-nextjs-tf"
      min_cpu = 0
      max_cpu = 8
      }
       
      output "database_connection_string" {
      // Output as an example, you can give this string to your application
      value = format("postgres://%s:%s@%s",
        scaleway_iam_application.app.id,
        scaleway_iam_api_key.api_key.secret_key,
        trimprefix(scaleway_sdb_sql_database.database.endpoint, "postgres://"),
      )
      sensitive = true
      }
       
      output "container_url" {
      // Output as an example, you can give this string to your application
      value = scaleway_container.main.domain_name
      sensitive = true
      }
    ```
    The terraform file creates several resources:
    - A [Serverless Container](https://www.scaleway.com/en/docs/serverless/containers/quickstart/) which host your NextJS application. This container is a group structure called a [Serverless Container Namespace](https://www.scaleway.com/en/docs/serverless/containers/concepts/#namespace)
    - A [Serverless SQL Database](https://www.scaleway.com/en/docs/serverless/sql-databases/reference-content/serverless-sql-databases-overview/) which stores posts data
    - An [IAM application](https://www.scaleway.com/en/docs/identity-and-access-management/iam/concepts/#application) used as a principal for your application
    - An [IAM secret key](https://www.scaleway.com/en/docs/identity-and-access-management/iam/concepts/#api-key) used as credentials to authenticate your application to the database
    - An [IAM policy](https://www.scaleway.com/en/docs/identity-and-access-management/iam/concepts/#policy) used giving the right permissions to your application principal  
    Then, initialize terraform module:
    ```bash
    terraform init
    ``` 
10. Provide `REGISTRY_ENDPOINT` and `DEFAULT_PROJECT_ID` environment variables to Terraform with:
    ```bash
    export TF_VAR_REGISTRY_ENDPOINT=$REGISTRY_ENDPOINT
    export TF_VAR_DEFAULT_PROJECT_ID=$(scw config get default-project-id) 
    ```
    Create and review Terraform plan:
    ```bash
    terraform plan
    ```
    Deploy your application using Terraform:
    ```bash
    terraform apply
    ```
    The output will provide the urls to which you can access your application and the connection string for your database. Since they contain sensitive value, they are hidden by default, but you can display them with `terraform output -json` command:
    ```bash
    terraform output -json
    {
     "container_url": {
       "sensitive": true,
       "type": "string",
       "value": "tutorialnextjsblogtfaxtypxrf-tutorial-nextjs-blog-tf.functions.fnc.fr-par.scw.cloud"
     },
     "database_connection_string": {
       "sensitive": true,
       "type": "string",
       "value": "postgres://a9773bf6-f6b7-40cc-9ae8-7f24e64c6531:604db597-c770-46ea-b785-94bf39536e6a@650c9680-1100-48e4-b5a6-ff2ff5bcf142.pg.sdb.fr-par.scw.cloud:5432/tutorial-nextjs-tf?sslmode=require"
     }
    }
    ```
    You can access your application by accessing the container url from your browser. However no post will display because your database in still empty. Let's add data in the next step.
11. Connect to your database using the link provided by Terraform's output:
    ```bash
    psql $(terraform output -json | jq -r '.database_connection_string.value')
    ```
    You should now see an input field with the name of your database:
    ```bash
    psql (15.3, server 16.1 (Debian 16.1-1.pgdg120+1))
    SSL connection (protocol: TLSv1.3, cipher: TLS_AES_128_GCM_SHA256, compression: off)
    Type "help" for help. 
 
    tutorial-nextjs=> 
    ```
    Create a table structure and add data using the following command: 
    ```sql
    create table posts (title char(100), excerpt text, coverImage text, date date, author_name char(50), author_picture text, ogImage_url text, slug char(50), content text);
 
    insert into posts values 
      ('Learn How to Pre-render Pages Using Static Generation with Next.js', 
     'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Praesent elementum facilisis leo vel fringilla est ullamcorper eget. At imperdiet dui  accumsan sit amet nulla facilities morbi tempus.',
      '/assets/blog/hello-world/cover.jpg',
      '2020-03-16T05:35:07.322Z',
      'Tim Neutkens',
      '/assets/blog/authors/tim.jpeg',
      '/assets/blog/hello-world/cover.jpg',
      'hello-world',
     'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Praesent elementum facilisis leo vel fringilla est ullamcorper eget. At imperdiet dui  accumsan sit amet nulla facilities morbi tempus. Praesent elementum facilisis leo vel fringilla. Congue mauris rhoncus aenean vel. Egestas sed tempus urna et pharetra pharetra massa massa ultricies.'),
      ('Dynamic Routing and Static Generation', 
     'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Praesent elementum facilisis leo vel fringilla est ullamcorper eget. At imperdiet dui  accumsan sit amet nulla facilities morbi tempus.',
      '/assets/blog/dynamic-routing/cover.jpg',
      '2020-03-16T05:35:07.322Z',
      'JJ Kasper',
      '/assets/blog/authors/jj.jpeg',
      '/assets/blog/dynamic-routing/cover.jpg',
      'dynamic-routing',
     'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Praesent elementum facilisis leo vel fringilla est ullamcorper eget. At imperdiet dui  accumsan sit amet nulla facilities morbi tempus. Praesent elementum facilisis leo vel fringilla. Congue mauris rhoncus aenean vel. Egestas sed tempus urna et pharetra pharetra massa massa ultricies.'),
      ('Preview Mode for Static Generation', 
     'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Praesent elementum facilisis leo vel fringilla est ullamcorper eget. At imperdiet dui  accumsan sit amet nulla facilities morbi tempus.',
      '/assets/blog/preview/cover.jpg',
      '2020-03-16T05:35:07.322Z',
      'Joe Haddad',
      '/assets/blog/authors/joe.jpeg',
      '/assets/blog/preview/cover.jpg',
      'preview',
     'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Praesent elementum facilisis leo vel fringilla est ullamcorper eget. At imperdiet dui  accumsan sit amet nulla facilities morbi tempus. Praesent elementum facilisis leo vel fringilla. Congue mauris rhoncus aenean vel. Egestas sed tempus urna et pharetra pharetra massa massa ultricies.');
    ```
    Exit database query interface by typing:
    ```bash
    exit
    ```
12. To check that everything works fine, go to the container url in your browser. The blog application should display.
    You can show the last terraform output (including the container url) with the following command:
    ```bash
    terraform output -json
    ```
13. **Congrats!** Your serverless application is now online!

    To go further you can:
    - **Edit the source code** locally, then:
      - Build a new container version with `docker build -t my-nextjs-blog .`
      - Push the new container to **Container Registry** with `docker tag my-nextjs-blog:latest $REGISTRY_ENDPOINT/my-nextjs-blog:latest` and `docker push $REGISTRY_ENDPOINT/my-nextjs-blog:latest`
      - Update and apply a new container deployment using `export CONTAINER_ID=$(scw container container list name="my-nextjs-blog" -o json | jq -r '.[0].id')` and `scw container container deploy $CONTAINER_ID`, and see the **new version being live** in seconds.
    - Inspect your **newly created resources** in the Scaleway Console
      - You can display your newly created **registry namespace** and **container image** in the [Container Registry Console Page](https://console.scaleway.com/registry/namespaces) 
      - You can display your newly created **serverless container namespace** and **container deployment** in the [Serverless Containers Console Page](https://console.scaleway.com/containers/namespaces)
      - You can display your newly created **serverless sql database** in the [Serverless SQL Database Console Page](https://console.scaleway.com/serverless-db/databases)
    - **Fine tune deployment options** such as autoscaling, targeted regions and much more. You can find more information by typing `scw deploy --help` in your terminal.
    - **Change security configuration** for your container. It is currently public and anyone with the link can access it, but you can make it [private to require authentication](https://www.scaleway.com/en/developers/api/serverless-containers/#authentication).
 

## Troubleshooting

If you encounter any issue, please first check you met all the prerequisite requirements.

- You have installed and configured [Scaleway CLI](https://www.scaleway.com/en/docs/developer-tools/scaleway-cli/quickstart/)
  - Running the following command in your terminal: 
    ```bash
    scw account project get
    ```
    should display your current project and organization id:
    ```
    ID              aedb23ea-7afa-4a5d-9f6c-27db072f1527
    Name            default
    OrganizationID  aedb23ea-7afa-4a5d-9f6c-27db072f1527
    CreatedAt       1 year ago
    UpdatedAt       1 year ago
    Description     -
    ```
  You can also find and compare your project and organization id in the [Scaleway Console Settings](https://console.scaleway.com/project/settings) 

- You have **Docker Engine** installed
  - Running the following command in your terminal:
    ```bash
    docker -v
    ```
    should display your currently installed docker version:
    ```
    Docker version 24.0.6, build ed223bc820
    ```
- You have the right [IAM permissions]((https://www.scaleway.com/en/docs/identity-and-access-management/iam/reference-content/policy/)), specifically **ContainersRegistryFullAccess**, **ContainersFullAccess** and **ServerlessSQLDatabaseFullAccess**
  - Running the following command in your terminal:
    ```bash
    scw registry namespace create
    ```
    should let you create a **registry namespace** and display its id and name:
    ```
    ID              5b955c35-c7a9-4340-9034-05a8171eca6a
    Name            cli-ns-peaceful-visvesvaraya
    Description     -
    OrganizationID  aedb23ea-7afa-4a5d-9f6c-27db072f1527
    ProjectID       aedb23ea-7afa-4a5d-9f6c-27db072f1527
    Status          ready
    StatusMessage   -
    Endpoint        rg.fr-par.scw.cloud/cli-ns-peaceful-visvesvaraya
    IsPublic        false
    Size            0 B
    CreatedAt       now
    UpdatedAt       now
    ImageCount      0
    Region          fr-par
    ```
    You can then delete your newly created namespace with the command:
    ```
    scw registry namespace delete {namespace-id}
    ``` 
    where `{namespace-id}` would be `5b955c35-c7a9-4340-9034-05a8171eca6a` in this example.

  - Running the following command in your terminal:
    ```bash
    scw container namespace create
    ```
    should let you create a **container namespace** and display its id and name:
    ```
    ID                   af106877-e327-4b0b-aa6e-f03f72609f9d
    Name                 cli-cns-sad-mendel
    OrganizationID       aedb23ea-7afa-4a5d-9f6c-27db072f1527
    ProjectID            aedb23ea-7afa-4a5d-9f6c-27db072f1527
    Status               pending
    RegistryNamespaceID  -
    RegistryEndpoint     -
    Description          -
    Region               fr-par
    ```
    You can then delete your newly created namespace with the command:
    ```bash
    scw container namespace delete {namespace-id}
    ``` 
    where `{namespace-id}` would be `af106877-e327-4b0b-aa6e-f03f72609f9d` in this example.

  - Running the following command in your terminal:
    ```bash
    scw sdb database create
    ```
    should let you create a **serverless sql database** and display its id and name:
    ```
    ID                   b8e406877-e327-4b0b-aa6e-f03f72609f9d
    Name                 cli-cns-sad-mendel
    OrganizationID       aedb23ea-7afa-4a5d-9f6c-27db072f1527
    ProjectID            aedb23ea-7afa-4a5d-9f6c-27db072f1527
    Status               pending
    Description          -
    Region               fr-par
    ```
    You can then delete your newly created database with the command:
    ```bash
    scw sdb database delete {database-id}
    ``` 
    where `{namespace-id}` would be `b8e406877-e327-4b0b-aa6e-f03f72609f9d` in this example.


