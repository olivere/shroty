task :deploy do
  `git push web`
  `ssh as2 rake compile restart`
end

task :compile do
  `go build -o shroty`
end

task :restart do
end
