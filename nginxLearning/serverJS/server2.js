const http = require('http');

const server = http.createServer((req, res) => {
    res.writeHead(200, { 'Content-Type': 'text/plain' });
    res.end('Response from Server 2\n');
});

server.listen(3002, () => {
    console.log('Server 2 is running on http://localhost:3002');
});