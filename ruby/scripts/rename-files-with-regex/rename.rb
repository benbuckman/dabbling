#! /usr/bin/env ruby

require 'pry'
require 'optparse'

module FileRenamingScript

  def self.parse_args
    @options = {}
    required_params = []

    OptionParser.new do |params|
      params.banner = "Usage: #{__FILE__} [@options]"

      params.on("-h", "--help") do
        puts params.help
        exit
      end

      required_params << :dir
      params.on("--dir=DIR", "(Required) Base directory") do |dir|
        @options[:dir] = dir
      end

      required_params << :from
      params.on("--from=PATTERN", "(Required) Regex pattern to replace, without wrapping slashes (e.g. '\\d+')") do |from|
        @options[:from] = Regexp.new(from)
      end

      @options[:to] = ''
      params.on("--to=STRING", "(Required) New string to substitute for regex pattern (default to empty string)") do |to|
        @options[:to] = to
      end

      @options[:real] = false
      params.on("--real", "Actually rename. Otherwise dryrun.") do
        @options[:real] = true
      end

      params.on("--extensions=[EXTS]", "File extensions (comma-separated)") do |extensions|
        @options[:extensions] = extensions.split(',') unless extensions.empty?
      end

      # using `on_tail` so it parses `on`'s first, but still has access to `help` :-/
      params.on_tail do
        # OptionParser is silly, required args aren't actually required,
        # they just require a value if the key is passed.
        required_params.each do |key|
          if @options[key].nil? || (@options[key].is_a?(String) && @options[key].empty?)
            puts "Missing argument: #{key}"
            puts params.help
            exit
          end
        end
      end
    end.parse!

    @options
  end


  def self.handle_dir(dirpath)
    puts dirpath

    # `Dir.children` doesn't work in ruby 2.4.1,
    # see https://stackoverflow.com/questions/45719260/ruby-2-4-1-dir-children-dirname-returns-undefined-method-children-for-dir
    children = Dir.entries(dirpath) - [".", ".."]

    children.each do |filename|
      filepath = File.expand_path(filename, dirpath)
      if File.directory?(filepath)
        handle_dir(filepath)
      else
        handle_file(filepath)
      end
    end
  end

  def self.want_extension?(filepath)
    ext = File.extname(filepath)
    if @options[:extensions].nil?
      true
    elsif @options[:extensions].any? {|wanted_ext| ext == ".#{wanted_ext}"}
      true
    else
      false
    end
  end

  def self.handle_file(filepath)
    puts "#{' '*2} #{File.basename(filepath)}"

    unless want_extension?(filepath)
      puts "#{' '*4} ignoring, wrong extension"
      return
    end

    rename_file(filepath)
  end

  def self.rename_file(orig_filepath)
    path_parts = orig_filepath.split(File::SEPARATOR)
    filename = path_parts.pop
    unless filename.match?(@options[:from])
      puts "#{' '*4} does not match pattern"
      return
    end
    filename = filename.gsub(@options[:from], @options[:to])
    path_parts << filename
    filepath = path_parts.join(File::SEPARATOR)


    if @options[:real]
      puts "#{' '*4} renaming to #{filename}"
      File.rename(orig_filepath, filepath)
    else
      puts "#{' '*4} would rename to #{filename}"
    end
  end

  def self.run
    parse_args
    puts "Options: #{@options}"

    # recursively traverse files
    handle_dir(@options[:dir])
  end
end

FileRenamingScript.run if $0 == __FILE__