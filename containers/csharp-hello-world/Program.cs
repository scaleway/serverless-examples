var builder = WebApplication.CreateBuilder(args);

// Get PORT from environment, default to 8080 if not set
var port = Environment.GetEnvironmentVariable("PORT") ?? "8080";
builder.WebHost.UseUrls($"http://*:{port}");

var app = builder.Build();

app.MapGet("/", () => "Hello from Scaleway!");

app.Run();
