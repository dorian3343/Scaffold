require 'time'

def run_test(directory)
  puts "Test #{directory} start #{Time.now.strftime("%T")}"
  system("cd #{directory} && ruby test.rb")
  exit_code = $?.exitstatus
  if exit_code != 0
    puts "Error: test.rb in #{directory} directory failed with exit code #{exit_code}"
    exit(exit_code)
  end
  puts "Test #{directory} end #{Time.now.strftime("%T")}"
end

run_test('database')
run_test('no_database')

puts "E2E Tests Passed! Yippee"
