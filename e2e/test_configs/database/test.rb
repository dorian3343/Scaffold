require 'net/http'
require 'json'

# Start the server and redirect logging
system("start /b go run main.go > NUL 2>&1")
sleep(5) # Wait for server to start

# Test 1 -> Check for null response
uri = URI('http://localhost:8080/get_user')
response = Net::HTTP.get(uri)

if response.strip != "null"
  puts "Wrong json response received in Test 1: #{response} should be null"
  system("taskkill /f /im main.exe > NUL 2>&1")
  File.delete("main.db")
  exit 1
end

# Test 2 -> add user correctly 
json_data = { "Name" => "John Doe", "Age" => 30 }
json_string = json_data.to_json

# Define the URL
url = URI.parse('http://localhost:8080/post_user')

# Create the HTTP request
http = Net::HTTP.new(url.host, url.port)
request = Net::HTTP::Post.new(url.request_uri)
request.body = json_string
request['Content-Type'] = 'application/json'
response = http.request(request)

# Check the response
if response.body.strip != "null"
  puts "Wrong json response received in Test 2: #{response.body} should be null"
  system("taskkill /f /im main.exe > NUL 2>&1")
  File.delete("main.db")
  exit 1
end

# Test 3 -> Check for user response
json = '[{"age":30,"id":1,"name":"John Doe"}]'
data = JSON.parse(json)

uri = URI('http://localhost:8080/get_user')
response = Net::HTTP.get(uri)
response_data = JSON.parse(response)

# Check if the response matches the expected data
if response_data != data
  puts "Response does not match expected data."
  puts "Expected: #{data}"
  puts "Received: #{response_data}"
  system("taskkill /f /im main.exe > NUL 2>&1")
  File.delete("main.db")
  exit 1
end

# Clean up
system("taskkill /f /im main.exe > NUL 2>&1")
File.delete("main.db")
exit 0
