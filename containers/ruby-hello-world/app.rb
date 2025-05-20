require 'webrick'

# Read the PORT environment variable, default to 8000 if not set
port = ENV['PORT'] || 8000

server = WEBrick::HTTPServer.new(Port: port.to_i)

server.mount_proc '/' do |req, res|
  res.body = "Hello from Scaleway!"
end

trap('INT') { server.shutdown }
server.start
