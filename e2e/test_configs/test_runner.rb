require 'time'

def run_test(directory)
  system("cd #{directory} && ruby test.rb")
  exit_code = $?.exitstatus
  if exit_code != 0
    puts "Error: test.rb in #{directory} directory failed with exit code #{exit_code}"
    puts "FAIL in #{directory}"
    exit(exit_code)
  end
   puts "ok"
end

run_test('database')
run_test('no_database')

