import OpenAI from 'openai';

const client = new OpenAI({
  baseURL: 'http://localhost:8080/v1',
  apiKey: 'dummy-key', // GoLLuM doesn't require real API keys
});

async function testGollum() {
  try {
    // Test models endpoint
    const models = await client.models.list();
    console.log('Available models:', models.data);

    // Test chat completion
    const completion = await client.chat.completions.create({
      model: 'toy-1',
      messages: [
        { role: 'user', content: 'Write a 7-word poem about llamas.' }
      ],
      max_tokens: 50,
      temperature: 0.7,
    });

    console.log('Response:', completion.choices[0].message.content);
  } catch (error) {
    console.error('Error:', error.message);
  }
}

testGollum();
