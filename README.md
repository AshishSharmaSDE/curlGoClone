# curlGoClone

Implementing curl functionality in Go Lang

# Build the App

go build -o curl-clone.exe

# Commands to test the application.

```bash
- curl-clone -X GET "https://jsonplaceholder.typicode.com/posts/1"
- echo '{"title": "foo", "body": "bar", "userId": 1}' > data.json
- curl-clone -X POST "https://jsonplaceholder.typicode.com/posts" -d @data.json
- curl-clone -X POST "https://jsonplaceholder.typicode.com/posts" -d '{"title": "foo", "body": "bar", "userId": 1}' -H "Content-Type: application/json"

```
