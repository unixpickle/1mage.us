util = require 'util'

if process.argv.length isnt 4
  util.error 'Usage: node server.js <port> <directory>'
  process.exit 1

if isNaN port = parseInt process.argv[2]
  util.error 'invalid port: ' + process.argv[2]
  process.exit 1
imageDirectory = process.argv[3]

module.exports =
  port: port
  directory: imageDirectory
