class Things
  #attr_accessor :names

  def initialize
    @names = []
  end

  def add(name)
    unless @names.include? name
      @names.push name.to_s
    end
  end

  def list
    if @names.nil? or @names.empty?
      puts "--- Missing ---"
    else
      @names.each do |name|
        puts "Hello #{name}"
      end
    end
  end

end