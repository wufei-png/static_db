#!/usr/bin/env ruby

defined = {}
IO.foreach(ARGV[0]) do |line|
  next unless line =~ /typedef.+?tcu([a-zA-Z0-9_]+)\((.*)\)/
  n = $1
  next if n == 'Init'
  next if defined[n]
  args = $2
  new_args = args.split(',').map{|e| e.split(' ').last.gsub('*', '')}.join(" ,")
  if new_args == 'void'
    new_args = ""
  end
  puts """CUresult wrapCu#{n}(#{args}) {
    return __cu#{n}(#{new_args});
}
  """
  defined[n] = true
end
