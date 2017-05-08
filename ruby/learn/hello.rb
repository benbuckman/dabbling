class Greeter1

  def initialize(name = "Stranger")
    @name = name
  end

  def say
    puts "Hello #{@name}"  #" or #{self.name}"
  end

  def to_s
    "Hello (name=#{@name})"
  end

end


p = Greeter1.new("Ben")

p.say
puts "string: #{p.to_s}"

###############################

class Greeter2
  attr_accessor :name
end