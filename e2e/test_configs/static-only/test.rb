require 'net/http'
require 'json'

# Start the server and redirect logging
system("start /b go run main.go > NUL 2>&1")
sleep(5)

#Test 1 -> Check static
uri = URI('http://localhost:8080/')
response = Net::HTTP.get_response(uri)

# Read the contents of the file
file_contents = File.read('./static/index.html')

# Check if response body is not equal to file contents
if response.code != "200"
    puts "Wrong response code."
    system("taskkill /f /im main.exe > NUL 2>&1")
    exit 1
end
# Test 2 -> Check API response
uri = URI('http://localhost:8080/api')
response = Net::HTTP.get_response(uri)

json_response = JSON.parse(response.body)

if json_response["message"] != "API"
    puts "Wrong response received in Test 2: #{json_response}"
    system("taskkill /f /im main.exe > NUL 2>&1")
    exit 1
end

system("taskkill /f /im main.exe > NUL 2>&1")
exit 0
