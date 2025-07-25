from openai import OpenAI
import openlit
import sys
import os

#
def printf(format, *args):
    sys.stdout.write(format % args)

#
def get_chat_response(system_message: str, user_request: str, seed: int = None):
  try:
      messages = [
          {"role": "system", "content": system_message},
          {"role": "user", "content": user_request},
      ]

      response = client.chat.completions.create(
          model="gpt-4.1-nano",
          messages=messages,
          seed=seed,
          max_tokens=200,
          n=1,
          temperature=1,
      )
      response_content = response.choices[0].message.content
      system_fingerprint = response.system_fingerprint
      prompt_tokens = response.usage.prompt_tokens
      completion_tokens = response.usage.total_tokens - response.usage.prompt_tokens

      printf("chat system message: %s\n", system_message)
      printf("chat user request: %s\n", user_request)

      printf("response message: %s\n", response_content)
      printf("response system_fingerprint: %s\n", system_fingerprint)
      printf("response prompt_tokens: %s\n", prompt_tokens)
      printf("response completion_tokens: %s\n", completion_tokens)

      return response_content

  except Exception as e:
      print(f"An error occurred: {e}")
      return None



# main
openlit.init(otlp_endpoint=os.environ.get('OTEL_EXPORTER_OTLP_ENDPOINT'))

#api_key_file = os.environ.get('OPENAI_TOKEN_FILE')
#file = open(api_key_file)
#api_key_value = file.read

api_key_value = os.environ.get('OPENAI_TOKEN')

# Initialise openai client
client = OpenAI(
    api_key=api_key_value
)

print ("-----------------------------------------------------------------------")
# inline basic test
chat_message = "Return a one liner from any movie for me to guess"

response = client.chat.completions.create(
    model="gpt-4.1-nano",
    messages=[
        {
            "role": "user",
            "content": chat_message,
        }
    ],
)

response_content = response.choices[0].message.content
system_fingerprint = response.system_fingerprint
prompt_tokens = response.usage.prompt_tokens
completion_tokens = response.usage.total_tokens - response.usage.prompt_tokens

printf("chat request: %s\n", chat_message)
printf("chat response: %s\n", response_content)

# Print the HTML table content
#display(HTML(table))


#print ("-----------------------------------------------------------------------")
# with system message input
#system_message = "You are a creative writing companion, aiding me in crafting captivating stories."
#user_request = "Compose a thrilling short story about a time-traveling detective who solves perplexing mysteries in the streets of Victorian London, blending historical elements with a touch of the supernatural."

#response = get_chat_response(system_message=system_message, user_request=user_request)

print ("-----------------------------------------------------------------------")
