# Create a serverless scraping architecture

This is the source code for the tutorial: [Create a serverless scraping architecture, with Scaleway Messaging and Queuing SQS, Serverless Functions and Managed Database](https://www.scaleway.com/en/docs/tutorials/create-serverless-scraping).

In this tutorial we show how to set up a simple application which reads [Hacker News](https://news.ycombinator.com/news) and processes the articles it finds there asynchronously, using Scaleway serverless products. 

## Requirements

This example assumes you are familiar with how serverless functions work. If needed, you can
check [Scaleway official documentation](https://www.scaleway.com/en/docs/serverless/functions/quickstart/)

This example is written using Python and Terraform, and assumes you have [set up authentication for the Terraform provider](https://registry.terraform.io/providers/scaleway/scaleway/latest/docs#authentication). 


## Context

**The architecture deployed in this tutorial consists of two functions, two triggers, a SQS queue, and a RDB instance.**
*The producer function, activated by a recurrent cron trigger, scrapes HackerNews for articles published in the last 15 minutes and pushes the title and URL of the articles to an SQS queue created with Scaleway Messaging and Queuing.*
*The consumer function, triggered by each new message on the SQS queue, consumes messages published to the queue, scrapes some data from the linked article, and then writes the data into a Scaleway Managed Database.*


## Setup
Once you have cloned this repository, you just need to package your functions and deploy them using Terraform. 
```bash
cd scraper
pip install -r requirements.txt --target ./package
zip -r functions.zip handlers/ package/
cd ../consumer
pip install -r requirements.txt --target ./package
zip -r functions.zip handlers/ package/
cd ../terraform 
terraform init
terraform apply
```


## Running

Everything is already up and running! 
You can check correct execution by using the Scaleway cockpit, and by connecting to your RDB instance to see results.

```bash
psql -h <DB_INSTANCE_IP> --port <DB_INSTANCE_PORT> -d hn-database -U worker
```
