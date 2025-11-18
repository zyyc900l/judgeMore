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
			"name": {
				"type": "text",
				"analyzer": "my_analyzer"
			},
			"category_id": { "type": "long" },
			"price": { "type": "double" },
			"shipping": { "type": "boolean" }
		}
	}
}`
