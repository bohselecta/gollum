import openai

client = openai.OpenAI(
    base_url="http://localhost:8080/v1",
    api_key="dummy-key"  # GoLLuM doesn't require real API keys
)

def test_gollum():
    try:
        # Test models endpoint
        models = client.models.list()
        print("Available models:", [m.id for m in models.data])
        
        # Test chat completion
        completion = client.chat.completions.create(
            model="toy-1",
            messages=[
                {"role": "user", "content": "Write a 7-word poem about llamas."}
            ],
            max_tokens=50,
            temperature=0.7
        )
        
        print("Response:", completion.choices[0].message.content)
        
    except Exception as e:
        print("Error:", str(e))

if __name__ == "__main__":
    test_gollum()
