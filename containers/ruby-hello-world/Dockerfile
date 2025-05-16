# Use the official Ruby image from the Docker Hub
FROM ruby:3.4.4

# Set the working directory in the container
WORKDIR /app

# Copy the Gemfile and Gemfile.lock into the container
COPY Gemfile ./

# Install the dependencies
RUN bundle install && bundle add rackup puma

# Copy the rest of the application code into the container
COPY app.rb .

# Run the application
CMD ["ruby", "app.rb"]
