"use client"

import { useState, useRef, useEffect } from "react"
import { Card } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Send, Sparkle, XCircle } from "lucide-react"

export default function ChatStreamUI() {
	const [question, setQuestion] = useState("")
	const [isLoading, setIsLoading] = useState(false)
	const [messages, setMessages] = useState([])
	const controllerRef = useRef(null)
	const messagesEndRef = useRef(null)
	const inputRef = useRef(null)

	// Auto-scroll to bottom when messages change
	useEffect(() => {
		messagesEndRef.current?.scrollIntoView({ behavior: "smooth" })
	}, [messages])

	// Focus input on load
	useEffect(() => {
		inputRef.current?.focus()
	}, [])

	const handleSubmit = async (e) => {
		e.preventDefault()
		if (!question.trim()) return

		const userMessage = { role: "user", content: question, id: Date.now().toString() }
		setMessages((prev) => [...prev, userMessage])
		setIsLoading(true)
		setQuestion("")

		const responseMessage = { role: "assistant", content: "", id: (Date.now() + 1).toString() }
		setMessages((prev) => [...prev, responseMessage])

		controllerRef.current = new AbortController()

		try {
			const res = await fetch("/completions", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({
					question,
					top_k: 3,
					threshold: 0.40,
					llm: { type: "ollama", model: "llama3.2:latest" },
				}),
				signal: controllerRef.current.signal,
			})

			if (!res.body) {
				setIsLoading(false)
				return
			}

			const reader = res.body.getReader()
			const decoder = new TextDecoder("utf-8")
			let buffer = ""

			while (true) {
				const { value, done } = await reader.read()
				if (done) break
				buffer += decoder.decode(value, { stream: true })

				const lines = buffer.split("\n")
				buffer = lines.pop() // carry over incomplete line

				for (const line of lines) {
					if (!line.trim()) continue
					try {
						const data = JSON.parse(line)
						setMessages((prev) => {
							const updated = [...prev]
							const last = { ...updated[updated.length - 1] } // clone the message
							last.content += data.response
							updated[updated.length - 1] = last
							return updated
						})
					} catch (err) {
						console.error("Failed to parse line:", line, err)
					}
				}
			}
		} catch (err) {
			if (err.name !== "AbortError") {
				console.error("Fetch error:", err)
				// Add error message
				setMessages((prev) => [
					...prev.slice(0, -1),
					{ role: "assistant", content: "Sorry, I encountered an error. Please try again.", id: Date.now().toString() },
				])
			}
		} finally {
			setIsLoading(false)
			inputRef.current?.focus()
		}
	}

	const cancelRequest = () => {
		if (controllerRef.current) {
			controllerRef.current.abort()
			controllerRef.current = null
			setIsLoading(false)
		}
	}

	return (
		<div className="max-w-4xl mx-auto py-8 px-4">
			<Card className="border border-gray-200 rounded-xl shadow-sm overflow-hidden">
				<div className="border-b border-gray-100 p-4 flex justify-between items-center">
					<div className="flex items-center gap-2">
						<Sparkle className="h-5 w-5 text-amber-400" />
						<h2 className="font-medium">Chat with your LiteTable Assistant</h2>
					</div>
					<div className="px-3 py-1 rounded-full bg-green-100 text-green-800 text-xs font-medium">llama3.2</div>
				</div>

				<ScrollArea className="h-[500px]">
					<div className="p-4">
						{messages.length === 0 && (
							<div className="flex flex-col items-center justify-center h-[400px] text-gray-400">
								<p className="text-center">Ask me anything to get started...</p>
							</div>
						)}

						{messages.map((msg) => (
							<div key={msg.id} className={`mb-4 ${msg.role === "user" ? "flex justify-end" : "flex justify-start"}`}>
								<div
									className={`
                    max-w-[80%] rounded-2xl px-4 py-3 shadow-sm
                    ${msg.role === "user" ? "bg-gray-500 text-white" : "bg-gray-100 text-gray-800"}
                  `}
								>
									<div className="whitespace-pre-wrap">{msg.content}</div>
								</div>
							</div>
						))}

						{isLoading && (
							<div className="flex items-center gap-2 text-gray-500 mb-4">
								<div className="flex space-x-1">
									<div
										className="w-2 h-2 rounded-full bg-gray-300"
										style={{
											animationName: "bounce",
											animationDuration: "1s",
											animationIterationCount: "infinite",
											animationDelay: "0ms",
										}}
									></div>
									<div
										className="w-2 h-2 rounded-full bg-gray-300"
										style={{
											animationName: "bounce",
											animationDuration: "1s",
											animationIterationCount: "infinite",
											animationDelay: "150ms",
										}}
									></div>
									<div
										className="w-2 h-2 rounded-full bg-gray-300"
										style={{
											animationName: "bounce",
											animationDuration: "1s",
											animationIterationCount: "infinite",
											animationDelay: "300ms",
										}}
									></div>
								</div>
							</div>
						)}
						<div ref={messagesEndRef} />
					</div>
				</ScrollArea>

				<div className="border-t border-gray-100 p-4">
					<form onSubmit={handleSubmit} className="flex items-center gap-2">
						<Input
							ref={inputRef}
							className="flex-1 border-gray-200 rounded-full focus-visible:ring-gray-300 focus-visible:ring-offset-0 focus-visible:border-gray-300"
							type="text"
							placeholder="Ask anything..."
							value={question}
							onChange={(e) => setQuestion(e.target.value)}
							disabled={isLoading}
						/>
						{isLoading ? (
							<Button
								type="button"
								variant="outline"
								size="icon"
								onClick={cancelRequest}
								className="rounded-full border-gray-200 hover:bg-gray-50"
							>
								<XCircle className="h-5 w-5 text-gray-500" />
							</Button>
						) : (
							<Button
								type="submit"
								disabled={!question.trim()}
								className="cursor-pointer bg-gray-500 hover:bg-gray-600 text-white rounded-full px-5"
							>
								<Send className="h-4 w-4 mr-2" />
								Send
							</Button>
						)}
					</form>
				</div>
			</Card>
		</div>
	)
}
