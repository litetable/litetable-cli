"use client"

import { useState, useRef, useEffect } from "react"
import { Card } from "@/components/ui/card.jsx"
import { Input } from "@/components/ui/input.jsx"
import { Button } from "@/components/ui/button.jsx"
import { ScrollArea } from "@/components/ui/scroll-area.jsx"
import { Send, Sparkle, XCircle, Search, BookOpen, Clock, ArrowUpRight } from "lucide-react"

export default function QuestionAnswerUI() {
	const [question, setQuestion] = useState("")
	const [isLoading, setIsLoading] = useState(false)
	const [currentQuery, setCurrentQuery] = useState(null)
	const [answers, setAnswers] = useState([])
	const controllerRef = useRef(null)
	const answerEndRef = useRef(null)
	const inputRef = useRef(null)

	// Auto-scroll to bottom when answer changes
	useEffect(() => {
		answerEndRef.current?.scrollIntoView({ behavior: "smooth" })
	}, [answers])

	// Focus input on load
	useEffect(() => {
		inputRef.current?.focus()
	}, [])

	const handleSubmit = async (e) => {
		e.preventDefault()
		if (!question.trim()) return

		// Set the current query
		setCurrentQuery(question)

		// Create a new answer entry
		const newAnswerId = Date.now().toString()
		setAnswers((prev) => [
			...prev,
			{
				id: newAnswerId,
				query: question,
				response: "",
				status: "loading",
				timestamp: new Date().toISOString(),
			},
		])

		setIsLoading(true)
		setQuestion("")

		controllerRef.current = new AbortController()

		try {
			const res = await fetch("/completions", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({
					question,
					top_k: 3,
					threshold: 0.4,
					llm: { type: "ollama", model: "llama3.2:latest" },
				}),
				signal: controllerRef.current.signal,
			})

			if (!res.body) {
				setIsLoading(false)
				setAnswers((prev) =>
					prev.map((a) =>
						a.id === newAnswerId
							? { ...a, status: "error", response: "Failed to get a response. Please try again." }
							: a,
					),
				)
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
						setAnswers((prev) =>
							prev.map((a) =>
								a.id === newAnswerId ? { ...a, response: a.response + data.response, status: "streaming" } : a,
							),
						)
					} catch (err) {
						console.error("Failed to parse line:", line, err)
					}
				}
			}

			// Mark as completed when done
			setAnswers((prev) => prev.map((a) => (a.id === newAnswerId ? { ...a, status: "completed" } : a)))
		} catch (err) {
			if (err.name !== "AbortError") {
				console.error("Fetch error:", err)
				setAnswers((prev) =>
					prev.map((a) =>
						a.id === newAnswerId
							? { ...a, status: "error", response: "Sorry, I encountered an error. Please try again." }
							: a,
					),
				)
			} else {
				// Handle abort - mark as cancelled
				setAnswers((prev) =>
					prev.map((a) =>
						a.id === newAnswerId ? { ...a, status: "cancelled", response: a.response || "Request cancelled." } : a,
					),
				)
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

	const formatTimestamp = (timestamp) => {
		const date = new Date(timestamp)
		return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })
	}

	return (
		<div className="max-w-5xl mx-auto py-8 px-4">
			<Card className="border border-gray-200 rounded-xl shadow-sm overflow-hidden">
				<div className="border-b border-gray-100 p-4 flex justify-between items-center">
					<div className="flex items-center gap-2">
						<Sparkle className="h-5 w-5 text-amber-400" />
						<h2 className="font-medium">LiteTable Research Assistant</h2>
					</div>
					<div className="px-3 py-1 rounded-full bg-green-100 text-green-800 text-xs font-medium">llama3.2</div>
				</div>

				<ScrollArea className="h-[500px]">
					<div className="p-4">
						{answers.length === 0 && (
							<div className="flex flex-col items-center justify-center h-[400px] text-gray-400">
								<Search className="h-12 w-12 mb-4 text-gray-300" />
								<h3 className="text-lg font-medium text-gray-500 mb-2">Ask me a research question</h3>
								<p className="text-center max-w-md text-gray-400">
									I'll search through available context to find the most relevant information for you.
								</p>
							</div>
						)}

						{answers.map((item) => (
							<div key={item.id} className="mb-8 last:mb-2">
								<div className="flex items-start gap-2 mb-2">
									<div className="bg-gray-500 text-white p-2 rounded-full">
										<Search className="h-4 w-4" />
									</div>
									<div className="flex-1">
										<div className="flex justify-between items-center mb-1">
											<h3 className="font-medium text-gray-800">{item.query}</h3>
											<span className="text-xs text-gray-400 flex items-center">
                        <Clock className="h-3 w-3 mr-1" />
												{formatTimestamp(item.timestamp)}
                      </span>
										</div>
									</div>
								</div>

								<div className="ml-10 pl-4 border-l-2 border-gray-200">
									<div className="bg-gray-100 rounded-xl p-4 shadow-sm">
										{item.status === "loading" && (
											<div className="flex items-center gap-2 text-gray-500">
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
												<span className="text-sm">Researching answer...</span>
											</div>
										)}

										{item.response && <div className="whitespace-pre-wrap text-gray-800">{item.response}</div>}

										{item.status === "completed" && item.response && (
											<div className="mt-4 pt-3 border-t border-gray-200 flex items-center justify-between">
												<div className="flex items-center text-xs text-gray-500">
													<BookOpen className="h-3 w-3 mr-1" />
													<span>Answer based on available context</span>
												</div>
												<Button
													variant="ghost"
													size="sm"
													className="text-xs text-gray-500 hover:text-gray-700 p-1 h-auto"
												>
													<ArrowUpRight className="h-3 w-3 mr-1" />
													View sources
												</Button>
											</div>
										)}
									</div>
								</div>
							</div>
						))}
						<div ref={answerEndRef} />
					</div>
				</ScrollArea>

				<div className="border-t border-gray-100 p-4">
					<form onSubmit={handleSubmit} className="flex items-center gap-2">
						<Input
							ref={inputRef}
							className="flex-1 border-gray-200 rounded-full focus-visible:ring-gray-300 focus-visible:ring-offset-0 focus-visible:border-gray-300"
							type="text"
							placeholder="Ask a research question..."
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
