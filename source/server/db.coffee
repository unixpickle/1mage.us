{MongoClient} = require 'mongodb'

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
      sorter = sequence: -1
      cursor = @collection.find({}).sort(sorter).limit 1
      cursor.nextObject (err, doc) =>
        cb err if err?
        @sequence = (doc?.sequence ? -1) + 1
        cb null, @sequence

  latest: (cb) -> cb null, @sequence

  findNextLast: (sequence, cb) ->
    gtQuery = sequence: $gt: sequence
    ltQuery = sequence: $lt: sequence
    cursor = @collection.find(gtQuery).sort(sequence: 1).limit 1
    cursor.nextObject (err, gtDoc) =>
      return cb err if err?
      cursor = @collection.find(ltQuery).sort(sequence: -1).limit 1
      cursor.nextObject (err, ltDoc) ->
        return cb err if err?
        cb null, gtDoc?.sequence, ltDoc?.sequence

  grabSequence: -> @sequence++

  insert: (mime, sequence, cb) ->
    doc = mime: mime, sequence: sequence
    @collection.insert doc, cb
  
  fetch: (sequence, cb) ->
    query = sequence: sequence
    cursor = @collection.find(query).limit 1
    cursor.nextObject cb

module.exports = new ImageDb()
