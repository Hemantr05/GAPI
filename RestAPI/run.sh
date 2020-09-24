go run server.go

go build .

gnome-terminal --tab --title="test" --command="bash -c 'curl http://localhost:3000/admin'"
