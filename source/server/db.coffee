{MongoClient} = require 'mongodb'
args = require './args'
pathJoin = require('path').join
fs = require 'fs'

class ImageDb
  constructor: ->
    @collection = null
    @sequence = null
  
  connect: (cb) ->
    info = 'mongodb://127.0.0.1:27017/1mage'
    MongoClient.connect info, (err, db) =>
      cb err if err?
      @collection = db.collection 'images'
      # TODO: ensure index here for id
    
      # query for the highest sequence number
      @latest (err, seq) =>
        cb err if err?
        @sequence = seq
        cb null, @sequence

  latest: (cb) ->
    sorter = sequence: -1
    cursor = @collection.find({}).sort(sorter).limit 1
    cursor.nextObject (err, doc) =>
      cursor.close()
      cb null, (doc?.sequence ? -1)

  findNextLast: (sequence, cb) ->
    gtQuery = sequence: $gt: sequence
    ltQuery = sequence: $lt: sequence
    cursor = @collection.find(gtQuery).sort(sequence: 1).limit 1
    cursor.nextObject (err, gtDoc) =>
      cursor.close()
      return cb err if err?
      cursor = @collection.find(ltQuery).sort(sequence: -1).limit 1
      cursor.nextObject (err, ltDoc) ->
        cursor.close()
        return cb err if err?
        cb null, gtDoc?.sequence, ltDoc?.sequence

  grabSequence: -> @sequence++

  insert: (mime, sequence, cb) ->
    doc = mime: mime, sequence: sequence
    @collection.insert doc, cb
  
  fetch: (sequence, cb) ->
    query = sequence: sequence
    cursor = @collection.find(query).limit 1
    cursor.nextObject (args...) ->
      cursor.close()
      cb args...

  delete: (sequence, cb) ->
    query = sequence: sequence
    @collection.remove query, (err) ->
      return cb err if err?
      deletePath = pathJoin args.directory, sequence + ''
      fs.unlink deletePath, cb

module.exports = new ImageDb()
