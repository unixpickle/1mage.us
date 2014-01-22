if process.argv.length isnt 4
  console.log 'Usage: node server.js <port> <directory>'
  process.exit 1

if isNaN port = parseInt process.argv[2]
  console.log 'invalid port: ' + process.argv[2]
  process.exit 1
imageDirectory = process.argv[3]

express = require 'express'
{MongoClient} = require 'mongodb'
http = require 'http'

app = express()
app.use express.static __dirname + '/assets'

server = http.createServer app
server.listen port
