http = require 'http'
util = require 'util'
path = require 'path'

express = require 'express'

args = require './args'
db = require './db'
upload = require './upload'
nextlast = require './nextlast'
newest = require './newest'
etc = require './etc'

main = ->
  app = express()
  # app.use express.logger()
  app.use express.urlencoded()
  app.use express.json()
  
  # use both the static assets and the compiled ones
  app.use express.static path.resolve __dirname + '/../assets'
  app.use express.static path.resolve __dirname + '/../../static'
  
  indexFile = path.resolve __dirname + '/../../static/index.html'
  app.use (req, res, next) ->
    res.sendHome = -> res.sendfile indexFile
    next()
  
  app.post '/upload', upload
  app.get /^\/nextlast\/[0-9]+$/, nextlast
  app.get /^\/(nav\/*|error)$/, (req, res) -> res.sendHome()
  app.get '/last', newest
  app.get '/*', etc
  app.get '*', (req, res) -> res.redirect '/error'

  server = http.createServer app
  server.listen args.port
  server.on 'listening', ->
    console.log 'listening on port', args.port

db.connect (err) ->
  if err?
    util.error err
    process.exit 1
  main()
