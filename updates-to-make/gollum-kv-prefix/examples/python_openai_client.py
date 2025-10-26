"""
Minimal OpenAI-compatible client targeting GoLLuM.

pip install openai
python examples/python_openai_client.py
"""
from openai import OpenAI
import os

client = OpenAI(
    base_url=os.environ.get("GOLLUM_BASE_URL", "http://localhost:8080/v1"),
    api_key=os.environ.get("OPENAI_API_KEY", "not-used"),
)

with client.chat.completions.stream(
    model="toy-1",
    messages=[{"role":"user","content":"Write a 7-word poem about llamas."}],
) as stream:
    for event in stream:
        if event.type == "chat.completion.chunk":
            delta = event.data.choices[0].delta.get("content") or ""
            print(delta, end="")
    print()
