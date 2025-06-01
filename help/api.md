# API Reference
---------------

Returns a response in the following format:
```json
{
    "sucess": true,                 //boolean
    "message": "server message",    //a message from the server
    "page": 1,                      //curent page you are on
    "pages": 18,                    //max number of pages
    "per_page": 20,                 //how many results are per page
    "skins": [                      //array containing skin items
        {                           //skin item
            "id": Int,
            "name": "String",
            "author": "String",
            "license": "String",
            "type": "String",
            "img": "String"         //base64 encoded image content string
        },
        ...repeat skin items
    ]
}
```

And accepts two GET params - these are capped at 1000:
```
per_page
page
```