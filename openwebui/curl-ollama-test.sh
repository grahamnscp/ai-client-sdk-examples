
curl --insecure https://$OLLAMA_HOSTNAME:443/api/generate -d '{
  "model": "gemma:7B",
  "prompt": "Why is the sky blue?"
}'

exit

curl https://$OLLAMA_HOSTNAME:443/vi/api/chat -sk -d '{
   "model": "gemma:7B",
   "messages": [
     {
       "role": "user",
       "content": "How hot is the Sun?"
     }
   ],
   "stream": false,
 }'

exit

curl -sk https://$OLLAMA_HOSTNAME/api/chat/completions \
    -H "Content-Type: application/json" \
    -X POST \
    -d '{
          "model": "gemma:7B",
          "messages": [
            {
              "role": "system",
              "content": "You are a helpful assistant."
            },
            {
              "role": "user",
              "content": "How hot is the Sun?"
            }
          ]
        }' 


