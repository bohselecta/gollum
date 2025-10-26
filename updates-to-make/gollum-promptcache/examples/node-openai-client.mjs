/**
 * Minimal OpenAI-compatible client targeting GoLLuM.
 * Requires: npm i openai
 * Run: node examples/node-openai-client.mjs
 */
import OpenAI from "openai";

const client = new OpenAI({
  baseURL: process.env.GOLLUM_BASE_URL || "http://localhost:8080/v1",
  apiKey: process.env.OPENAI_API_KEY || "not-used",
});

const res = await client.chat.completions.create({
  model: "toy-1",
  messages: [{ role: "user", content: "Write a 7-word poem about llamas." }],
  stream: true,
});

for await (const chunk of res) {
  const delta = chunk.choices?.[0]?.delta?.content || "";
  process.stdout.write(delta);
}
process.stdout.write("\n");
