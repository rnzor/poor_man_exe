"""
Demo AI Chat Backend - Real-time WebSocket Chat
Poor Man's exe.dev - Demo Application
"""

import asyncio
import json
import random
import uuid
from datetime import datetime
from fastapi import FastAPI, WebSocket, WebSocketDisconnect, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import HTMLResponse
from pydantic import BaseModel
from typing import List, Optional
import httpx

app = FastAPI(title="AI Chat API", version="1.0.0")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


class Message(BaseModel):
    role: str
    content: str
    timestamp: Optional[str] = None


class ChatRequest(BaseModel):
    message: str
    stream: bool = True


class ChatResponse(BaseModel):
    response: str
    conversation_id: str
    timestamp: str


conversations = {}

AI_RESPONSES = [
    "That's a fascinating perspective! Tell me more.",
    "I see what you mean. How does that make you feel?",
    "Interesting! What led you to that conclusion?",
    "That's a great point. Have you considered the alternative?",
    "I'm learning so much from this conversation. What else?",
    "That's really insightful! How might we explore this further?",
    "I appreciate you sharing that. What do you think comes next?",
    "That's thought-provoking. Let's dig deeper into that idea.",
]


def generate_ai_response(user_message: str) -> str:
    """Generate a thoughtful AI-like response"""
    user_lower = user_message.lower()

    if any(word in user_lower for word in ["hello", "hi", "hey"]):
        responses = [
            "Hello! Great to connect with you today. What's on your mind?",
            "Hi there! I'm excited to chat. What would you like to discuss?",
            "Hey! Welcome to the conversation. What brings you here?",
        ]
        return random.choice(responses)

    if any(word in user_lower for word in ["how are you", "how do you work"]):
        return "I'm doing well, thanks for asking! I'm a simple demo chatbot running on FastAPI with WebSocket support. I use async Python to handle multiple connections efficiently. How about you?"

    if any(word in user_lower for word in ["what is", "explain", "tell me about"]):
        return f"You asked about '{user_message}'. That's a great topic! While I'm a demo bot, I can help explore ideas. What specific aspect interests you most?"

    if any(word in user_lower for word in ["help", "capabilities", "what can you do"]):
        return "I'm a demo chat application showcasing real-time WebSocket communication! I can:\n\n• Hold conversations in real-time\n• Stream responses character by character (like ChatGPT)\n• Remember our conversation context\n• Run entirely in Docker containers\n• Scale across multiple instances\n\nPretty cool for a simple demo, right?"

    return (
        random.choice(AI_RESPONSES)
        + f'\n\n(You said: "{user_message[:50]}{"..." if len(user_message) > 50 else ""}")'
    )


@app.get("/")
async def root():
    return {
        "name": "AI Chat Demo",
        "version": "1.0.0",
        "status": "running",
        "endpoints": {
            "websocket": "/ws/{conversation_id}",
            "chat": "/chat",
            "health": "/health",
        },
    }


@app.get("/health")
async def health():
    return {"status": "healthy", "timestamp": datetime.utcnow().isoformat()}


@app.post("/chat", response_model=ChatResponse)
async def chat(request: ChatRequest):
    conversation_id = (
        request.message[:8] if len(request.message) > 8 else request.message
    )

    response_text = generate_ai_response(request.message)

    return ChatResponse(
        response=response_text,
        conversation_id=conversation_id,
        timestamp=datetime.utcnow().isoformat(),
    )


@app.websocket("/ws/{conversation_id}")
async def websocket_endpoint(websocket: WebSocket, conversation_id: str):
    await websocket.accept()

    if conversation_id not in conversations:
        conversations[conversation_id] = []

    conversation = conversations[conversation_id]

    try:
        while True:
            data = await websocket.receive_text()
            try:
                message_data = json.loads(data)
                user_message = message_data.get("message", "")

                if not user_message:
                    await websocket.send_json({"error": "Empty message"})
                    continue

                conversation.append(
                    {
                        "role": "user",
                        "content": user_message,
                        "timestamp": datetime.utcnow().isoformat(),
                    }
                )

                ai_response = generate_ai_response(user_message)

                conversation.append(
                    {
                        "role": "assistant",
                        "content": ai_response,
                        "timestamp": datetime.utcnow().isoformat(),
                    }
                )

                if message_data.get("stream", True):
                    for char in ai_response:
                        await websocket.send_json(
                            {
                                "type": "chunk",
                                "content": char,
                                "conversation_id": conversation_id,
                            }
                        )
                        await asyncio.sleep(0.02)

                    await websocket.send_json(
                        {
                            "type": "done",
                            "conversation_id": conversation_id,
                            "full_response": ai_response,
                        }
                    )
                else:
                    await websocket.send_json(
                        {
                            "type": "complete",
                            "response": ai_response,
                            "conversation_id": conversation_id,
                        }
                    )

            except json.JSONDecodeError:
                await websocket.send_json({"error": "Invalid JSON"})
            except Exception as e:
                await websocket.send_json({"error": str(e)})

    except WebSocketDisconnect:
        pass


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host="0.0.0.0", port=8000)
