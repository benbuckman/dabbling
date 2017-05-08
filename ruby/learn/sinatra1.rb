require 'sinatra'
require 'haml'

get '/' do
  'Hello world!'
end

get '/a' do haml :index end

__END__

@@index %html %body %h1 Welcome to Our Site

