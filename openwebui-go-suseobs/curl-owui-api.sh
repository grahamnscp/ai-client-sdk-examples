#!/bin/bash

curl --insecure -X POST https://$OPEN_WEBUI_HOSTNAME/api/chat/completions \
-H "Authorization: Bearer $OPEN_WEBUI_API_KEY" \
-H "Content-Type: application/json" \
-d '{
      "model": "gemma:2b",
      "messages": [
        {
          "role": "user",
          "content": "how tall is the shard in london?"
        }
      ]
    }'
