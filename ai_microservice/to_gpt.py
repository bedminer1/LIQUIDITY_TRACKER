from openai import OpenAI
from dotenv import load_dotenv
import os

load_dotenv("secrets.env")

client = OpenAI(
    api_key = os.getenv("OPENAI_API_KEY")
)

prompt = "How to create a python code that links to chat gpt?"

chat_completion = client.chat.completions.create(
    messages = [
        {
            "role":"user",
            "content":prompt
        }
    ],
    model = "gpt-4o-mini"
)

print(chat_completion.choices[0].message.content)