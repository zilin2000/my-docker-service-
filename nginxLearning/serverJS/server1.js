const http = require('http');

const server = http.createServer((req, res) => {
    res.writeHead(200, { 'Content-Type': 'text/plain' });
    res.end('Response from Server 1\n');
});

server.listen(3001, () => {
    console.log('Server 1 is running on http://localhost:3001');
});