require 'net/http'
require 'json'

# Start the server and redirect logging
system("start /b go run main.go > NUL 2>&1")
sleep(5)

# Test 1
uri = URI('http://localhost:8080/Greeting')
response = Net::HTTP.get_response(uri)

cors_found = false

response.each_capitalized do |key, value|
  if key == 'Access-Control-Allow-Origin'
    cors_found = true
    if value != "*"
      puts "Wrong CORS received: #{value}"
      system("taskkill /f /im main.exe > NUL 2>&1")
      exit 1
    end
  end
end

unless cors_found
  puts "CORS header not found"
  system("taskkill /f /im main.exe > NUL 2>&1")
  exit 1
end

# Convert response to JSON -> Check json response with CORS

json_response = JSON.parse(response.body)

if json_response["message"] != "Hello"
    puts "Wrong json response received in Test 1: #{json_response}"
    system("taskkill /f /im main.exe > NUL 2>&1")
    exit 1
end


# Test 2 -> Check JSON response
uri = URI('http://localhost:8080/status')
response = Net::HTTP.get_response(uri)
# Convert response to JSON
json_response = JSON.parse(response.body)

if json_response["status"] != "OK"
    puts "Wrong json response received in Test 2: #{json_response}"
    system("taskkill /f /im main.exe > NUL 2>&1")
    exit 1
end

# Test 3 -> Check int response
uri = URI('http://localhost:8080/int')
response = Net::HTTP.get_response(uri)
# Convert response to JSON
json_response = JSON.parse(response.body)




if json_response != 69
    puts "Wrong json response received in Test 3: #{json_response}"
   system("taskkill /f /im main.exe > NUL 2>&1")
    exit 1
end

cache_found = false

response.each_capitalized do |key, value|
  if key == 'Cache-Control'
    cache_found = true
    if value != "max-age=3600, public"
      puts "Wrong Cache received: #{value}"
      system("taskkill /f /im main.exe > NUL 2>&1")
      exit 1
    end
  end
end

unless cache_found
  puts "Cache header not found"
  system("taskkill /f /im main.exe > NUL 2>&1")
  exit 1
end


# Kill the process
system("taskkill /f /im main.exe > NUL 2>&1")
exit 0