# Monitoring glaciers with Scaleway Serverless Functions

Associated Scaleway blog post [link](https://blog.scaleway.com/how-to-collect-and-store-data-periodically-with-serverless-functions-2/)

## Instructions

1. Install serverless framework : [Quick Start](https://github.com/scaleway/serverless-scaleway-functions#quick-start)
1. Ensure `caas/.env` and `faas/.env` files are filled with your informations (also you must add `.env` to your `.gitignore`)
1. Then at the root of the directory, run `serverless deploy` (a prompt will ask you to install serverless compose if it's not already on your machine)

## Settings

In [faas](/faas/handler.go#L47) folder you must provide a valid URL pointing to the webcam image provider you want.
For example you can use your favorite search engine to find results about "Webcam Mont Blanc", providers often
offers easy way to retreive images from the webcam.

Example of `handler.go` modification :

```go
const unformattedURL = "https://data.skaping.com/chamonix/aiguille-du-midi/%d/%02d/%d/09-31.jpg"
```

## Documentation

This project is a part of a blog post to introduce some functionalities about serverless capabilities.

You can find documentations here :

- [Serverless Functions](https://www.scaleway.com/en/docs/compute/functions/reference-content/)
- [Serverless Containers](https://www.scaleway.com/en/docs/compute/containers/reference-content/)
- [Databases](https://www.scaleway.com/en/docs/managed-databases/postgresql-and-mysql/)
- [Object Storage](https://www.scaleway.com/en/docs/storage/object/)
- [Scaleway serverless plugin](https://github.com/scaleway/serverless-scaleway-functions) for deployment

## Contributing

Nature of this project is to demonstrate how easy is to develop a project using serverless approach.

Feel free to open issues or pull requests.
