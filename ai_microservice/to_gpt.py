from openai import OpenAI

client = OpenAI(
    api_key = "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"
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