package es

const mapping = `{
  "settings": {
    "analysis": {
      "analyzer": {
        "my_analyzer": {
          "type": "custom", 
          "tokenizer": "ik_max_word"
        }
      }
    }
  },    
  "mappings": {
    "properties": {
      "id": { "type": "long" },
      "recognized_event_name": {
        "type": "text",
        "analyzer": "my_analyzer",
        "fields": {
          "keyword": {
            "type": "keyword"
          }
        }
      },
      "organizer": {
        "type": "text", 
        "analyzer": "my_analyzer",
        "fields": {
          "keyword": {
            "type": "keyword"
          }
        }
      }
    }
  }
}`
