def increment_version(version)
  major, minor, patch = version.split('.').map(&:to_i)

  if patch < 9
    patch += 1
  elsif minor < 9
    minor += 1
    patch = 0
  else
    major += 1
    minor = 0
    patch = 0
  end

  "v.#{major}.#{minor}.#{patch}"
end

# Start building
start_time = Time.now
puts("Running Scaffold Build Tool!\n==================================")

# Run tests
test_output = `./Testing.bat`
if test_output.length > 1
    puts(test_output)
    puts("Testing Failed -> Build Failed")
    exit 1
end

puts("Testing passed")

# Read version from VERSION file
version = File.read("VERSION").strip
version = increment_version(version)

# Create directories for build artifacts
linux_dir = "Scaffold_#{version}_linux"
windows_dir = "Scaffold_#{version}_windows"

begin
    Dir.mkdir(linux_dir)
    Dir.mkdir(windows_dir)
rescue => e
    puts("Error creating directories: #{e.message}")
    exit 1
end

puts("Created directories: #{linux_dir}, #{windows_dir}")

# For linux
begin
    system("set GOOS=linux&& set GOARCH=amd64&& go build -o #{linux_dir}\\Scaffold")
    puts("Linux build completed successfully!")

    # Create main.yml for Linux
    File.write("#{linux_dir}/VERSION", version)
rescue => e
    puts("Error during Linux compilation: #{e.message}")
    exit 1
end

# For Windows
begin
    system("set GOOS=windows&& set GOARCH=amd64&& go build -o #{windows_dir}\\Scaffold.exe")
    puts("Windows build completed successfully!")

    # Create main.yml for Windows
    File.write("#{windows_dir}/VERSION", version)
rescue => e
    puts("Error during Windows compilation: #{e.message}")
    exit 1
end

end_time = Time.now
elapsed_time = end_time - start_time

puts "Build took: #{elapsed_time} seconds"
