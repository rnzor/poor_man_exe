import { useState, useRef, useEffect } from 'react'
import './App.css'

interface Message {
  id: string
  role: 'user' | 'assistant'
  content: string
  timestamp: Date
}

interface WebSocketMessage {
  type: 'chunk' | 'done' | 'complete' | 'error'
  content?: string
  response?: string
  conversation_id?: string
  full_response?: string
  error?: string
}

function App() {
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState('')
  const [isConnected, setIsConnected] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [ws, setWs] = useState<WebSocket | null>(null)
  const messagesEndRef = useRef<HTMLDivElement>(null)
  const conversationId = useRef<string>('demo-' + Math.random().toString(36).substring(7))

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }

  useEffect(() => {
    scrollToBottom()
  }, [messages])

  useEffect(() => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const wsUrl = `${protocol}//${window.location.host}/ws/${conversationId.current}`

    let socket: WebSocket
    let reconnectTimeout: NodeJS.Timeout

    const connect = () => {
      try {
        socket = new WebSocket(wsUrl)
        setWs(socket)

        socket.onopen = () => {
          setIsConnected(true)
          console.log('WebSocket connected')
        }

        socket.onclose = () => {
          setIsConnected(false)
          console.log('WebSocket disconnected, reconnecting in 3s...')
          reconnectTimeout = setTimeout(connect, 3000)
        }

        socket.onerror = (error) => {
          console.error('WebSocket error:', error)
        }

        socket.onmessage = (event) => {
          try {
            const data: WebSocketMessage = JSON.parse(event.data)

            if (data.type === 'chunk' && data.content) {
              setMessages(prev => {
                const lastMessage = prev[prev.length - 1]
                if (lastMessage && lastMessage.role === 'assistant') {
                  return [...prev.slice(0, -1), {
                    ...lastMessage,
                    content: lastMessage.content + data.content!
                  }]
                }
                return prev
              })
            } else if (data.type === 'done' || data.type === 'complete') {
              setIsLoading(false)
            } else if (data.type === 'error') {
              setIsLoading(false)
              setMessages(prev => [...prev, {
                id: Date.now().toString(),
                role: 'assistant',
                content: `Error: ${data.error}`,
                timestamp: new Date()
              }])
            }
          } catch (e) {
            console.error('Failed to parse WebSocket message:', e)
          }
        }
      } catch (error) {
        console.error('Failed to connect WebSocket:', error)
      }
    }

    connect()

    return () => {
      clearTimeout(reconnectTimeout)
      socket?.close()
    }
  }, [])

  const sendMessage = async () => {
    if (!input.trim() || isLoading) return

    const userMessage: Message = {
      id: Date.now().toString(),
      role: 'user',
      content: input.trim(),
      timestamp: new Date()
    }

    setMessages(prev => [...prev, userMessage])
    setInput('')
    setIsLoading(true)

    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({
        message: userMessage.content,
        stream: true
      }))
    } else {
      try {
        const response = await fetch('/chat', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ message: userMessage.content, stream: false })
        })
        const data = await response.json()
        setMessages(prev => [...prev, {
          id: Date.now().toString(),
          role: 'assistant',
          content: data.response,
          timestamp: new Date()
        }])
      } catch (error) {
        setMessages(prev => [...prev, {
          id: Date.now().toString(),
          role: 'assistant',
          content: 'Failed to send message. Please try again.',
          timestamp: new Date()
        }])
      } finally {
        setIsLoading(false)
      }
    }
  }

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      sendMessage()
    }
  }

  const clearChat = () => {
    setMessages([])
    conversationId.current = 'demo-' + Math.random().toString(36).substring(7)
  }

  return (
    <div className="app">
      <header className="header">
        <div className="header-content">
          <h1>AI Chat</h1>
          <span className={`status ${isConnected ? 'connected' : 'disconnected'}`}>
            {isConnected ? 'Connected' : 'Connecting...'}
          </span>
        </div>
        <p className="subtitle">Poor Man's exe.dev Demo - Real-time WebSocket Chat</p>
      </header>

      <main className="chat-container">
        <div className="messages">
          {messages.length === 0 && (
            <div className="welcome">
              <div className="welcome-icon">🤖</div>
              <h2>Welcome to AI Chat!</h2>
              <p>Start a conversation and I'll respond in real-time.</p>
              <p className="hint">Try saying "Hello" or "What can you do?"</p>
            </div>
          )}

          {messages.map(message => (
            <div key={message.id} className={`message ${message.role}`}>
              <div className="message-avatar">
                {message.role === 'user' ? '👤' : '🤖'}
              </div>
              <div className="message-content">
                <div className="message-text">{message.content}</div>
                <div className="message-time">
                  {message.timestamp.toLocaleTimeString()}
                </div>
              </div>
            </div>
          ))}

          {isLoading && (
            <div className="message assistant">
              <div className="message-avatar">🤖</div>
              <div className="message-content">
                <div className="typing-indicator">
                  <span></span>
                  <span></span>
                  <span></span>
                </div>
              </div>
            </div>
          )}

          <div ref={messagesEndRef} />
        </div>

        <div className="input-area">
          <textarea
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyPress={handleKeyPress}
            placeholder="Type your message..."
            rows={1}
            disabled={isLoading}
          />
          <button
            onClick={sendMessage}
            disabled={!input.trim() || isLoading}
            className="send-button"
          >
            {isLoading ? 'Sending...' : 'Send'}
          </button>
          <button
            onClick={clearChat}
            className="clear-button"
            title="Clear chat"
          >
            🗑️
          </button>
        </div>
      </main>

      <footer className="footer">
        <p>Running on Poor Man's exe.dev • Docker + FastAPI + React + WebSocket</p>
      </footer>
    </div>
  )
}

export default App
