require 'net/http'
require 'json'

# Start the server and redirect logging
system("start /b go run main.go > NUL 2>&1")
sleep(5)

# Test 1
uri = URI('http://localhost:8080/Greeting')
response = Net::HTTP.get(uri)
# Convert response to JSON -> Check json response with CORS
json_response = JSON.parse(response)

if json_response["message"] != "Hello"
    puts "Wrong json response received in Test 1: #{json_response}"
    system("taskkill /f /im main.exe > NUL 2>&1")
    exit 1
end


# Test 2 -> Check JSON response without CORS
uri = URI('http://localhost:8080/status')
response = Net::HTTP.get(uri)
# Convert response to JSON
json_response = JSON.parse(response)

if json_response["status"] != "OK"
    puts "Wrong json response received in Test 2: #{json_response}"
    system("taskkill /f /im main.exe > NUL 2>&1")
    exit 1
end

# Test 3 -> Check int response
uri = URI('http://localhost:8080/int')
response = Net::HTTP.get(uri)

# Convert response to JSON
json_response = JSON.parse(response)

if json_response != 69
    puts "Wrong json response received in Test 3: #{json_response}"
    system("taskkill /f /im main.exe > NUL 2>&1")
    exit 1
end

# Kill the process
system("taskkill /f /im main.exe > NUL 2>&1")
exit 0