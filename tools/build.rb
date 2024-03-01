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

# Create directories for build artifacts
linux_dir = "Scaffold_v#{version}_linux"
windows_dir = "Scaffold_v#{version}_windows"

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
    File.write("#{linux_dir}/main.yml", "")
    puts("Created main.yml for Linux")
rescue => e
    puts("Error during Linux compilation: #{e.message}")
    exit 1
end

# For Windows
begin
    system("set GOOS=windows&& set GOARCH=amd64&& go build -o #{windows_dir}\\Scaffold.exe")
    puts("Windows build completed successfully!")

    # Create main.yml for Windows
    File.write("#{windows_dir}/main.yml", "")
    puts("Created main.yml for Windows")
rescue => e
    puts("Error during Windows compilation: #{e.message}")
    exit 1
end
