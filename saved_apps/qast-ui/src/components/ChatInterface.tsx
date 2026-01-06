import React, { useState } from 'react';

interface ChatResponse {
    answer: string;
    context: string[];
}

interface IngestResponse {
    result: any;
    // Add other fields if needed for display
}

type Mode = 'chat' | 'ingest';

export function ChatInterface() {
    const [mode, setMode] = useState<Mode>('chat');

    const [expertiseLevel, setExpertiseLevel] = useState('novice');

    // Chat State
    const [query, setQuery] = useState('');
    const [chatResponse, setChatResponse] = useState<ChatResponse | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    // Ingest State
    const [ingestText, setIngestText] = useState('');
    const [ingestResponse, setIngestResponse] = useState<IngestResponse | null>(null);

    const handleChatSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!query.trim()) return;

        setLoading(true);
        setError(null);
        setChatResponse(null);

        try {
            const res = await fetch('http://localhost:8081/api/v1/rag/ask', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    query,
                    expertise_level: expertiseLevel
                }),
            });

            if (!res.ok) throw new Error(`Error: ${res.statusText}`);

            const data = await res.json();
            setChatResponse(data);
        } catch (err: any) {
            setError(err.message || 'Failed to fetch response');
        } finally {
            setLoading(false);
        }
    };

    const handleIngestSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!ingestText.trim()) return;

        setLoading(true);
        setError(null);
        setIngestResponse(null);

        try {
            const res = await fetch('http://localhost:8081/api/v1/privacy/extract', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    text: ingestText,
                    user_id: 'frontend_user',
                    template_name: 'extract_knowledge_graph'
                }),
            });

            if (!res.ok) throw new Error(`Error: ${res.statusText}`);

            const data = await res.json();
            setIngestResponse(data);
            setIngestText('');
        } catch (err: any) {
            setError(err.message || 'Failed to ingest data');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="flex flex-col h-screen max-w-4xl mx-auto p-6 bg-white dark:bg-zinc-900 text-zinc-900 dark:text-zinc-100">
            <header className="mb-8 flex justify-between items-center">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Qast RAG Interface</h1>
                    <p className="text-zinc-500 mt-2">Privacy-preserving Knowledge Retrieval</p>
                </div>

                <div className="flex items-center gap-4">
                    {/* Expertise Selector (Chat Mode Only) */}
                    {mode === 'chat' && (
                        <div className="flex items-center gap-2">
                            <label htmlFor="expertise" className="text-sm font-medium text-zinc-600 dark:text-zinc-400">Level:</label>
                            <select
                                id="expertise"
                                value={expertiseLevel}
                                onChange={(e) => setExpertiseLevel(e.target.value)}
                                className="px-3 py-2 rounded-md text-sm border border-zinc-300 dark:border-zinc-700 bg-white dark:bg-zinc-800 focus:outline-none focus:ring-2 focus:ring-blue-500"
                            >
                                <option value="novice">Novice</option>
                                <option value="expert">Expert</option>
                            </select>
                        </div>
                    )}

                    {/* Mode Toggle */}
                    <div className="flex bg-zinc-100 dark:bg-zinc-800 p-1 rounded-lg">
                        <button
                            onClick={() => setMode('chat')}
                            className={`px-4 py-2 rounded-md text-sm font-medium transition-all ${mode === 'chat' ? 'bg-white dark:bg-zinc-700 shadow-sm text-blue-600 dark:text-blue-400' : 'text-zinc-500 hover:text-zinc-700 dark:hover:text-zinc-300'}`}
                        >
                            Chat
                        </button>
                        <button
                            onClick={() => setMode('ingest')}
                            className={`px-4 py-2 rounded-md text-sm font-medium transition-all ${mode === 'ingest' ? 'bg-white dark:bg-zinc-700 shadow-sm text-emerald-600 dark:text-emerald-400' : 'text-zinc-500 hover:text-zinc-700 dark:hover:text-zinc-300'}`}
                        >
                            Ingest
                        </button>
                    </div>
                </div>
            </header>

            {/* Main Content Area */}
            <div className="flex-1 overflow-y-auto mb-6 p-4 rounded-xl bg-zinc-50 dark:bg-zinc-800/50 border border-zinc-200 dark:border-zinc-800 relative min-h-[300px]">

                {loading && (
                    <div className="absolute inset-0 z-10 bg-white/50 dark:bg-zinc-900/50 backdrop-blur-sm flex items-center justify-center text-blue-600">
                        <div className="flex flex-col items-center gap-2">
                            <span className="animate-spin text-2xl">⚡</span>
                            <span className="font-medium animate-pulse">{mode === 'chat' ? 'Generating Response...' : 'Ingesting Data...'}</span>
                        </div>
                    </div>
                )}

                {error && (
                    <div className="p-4 rounded-lg bg-red-50 text-red-600 border border-red-200 mb-4">
                        {error}
                    </div>
                )}

                {/* CHAT MODE OUTPUT */}
                {mode === 'chat' && (
                    <>
                        {!chatResponse && !loading && !error && (
                            <div className="absolute inset-0 flex items-center justify-center text-zinc-400">
                                <p>Ask a question about the knowledge base.</p>
                            </div>
                        )}
                        {chatResponse && (
                            <div className="space-y-6">
                                <div className="prose dark:prose-invert">
                                    <h3 className="text-sm font-semibold text-zinc-400 uppercase tracking-wider mb-2">Answer</h3>
                                    <div className="text-lg leading-relaxed whitespace-pre-wrap">
                                        {chatResponse.answer}
                                    </div>
                                </div>

                                {chatResponse.context && chatResponse.context.length > 0 && (
                                    <div className="border-t border-zinc-200 dark:border-zinc-700 pt-4 mt-6">
                                        <h3 className="text-sm font-semibold text-zinc-400 uppercase tracking-wider mb-3">Context Sources</h3>
                                        <div className="space-y-2">
                                            {chatResponse.context.map((ctx, idx) => (
                                                <div key={idx} className="text-xs bg-zinc-100 dark:bg-zinc-800 p-2 rounded border border-zinc-200 dark:border-zinc-700 text-zinc-600 dark:text-zinc-400 font-mono">
                                                    {ctx}
                                                </div>
                                            ))}
                                        </div>
                                    </div>
                                )}
                            </div>
                        )}
                    </>
                )}

                {/* INGEST MODE OUTPUT */}
                {mode === 'ingest' && (
                    <>
                        {!ingestResponse && !loading && !error && (
                            <div className="absolute inset-0 flex items-center justify-center text-zinc-400">
                                <p>Enter text to add to the knowledge graph.</p>
                            </div>
                        )}
                        {ingestResponse && (
                            <div className="space-y-4">
                                <div className="p-4 bg-emerald-50 dark:bg-emerald-900/20 text-emerald-700 dark:text-emerald-300 rounded border border-emerald-100 dark:border-emerald-800 mb-4">
                                    ✅ Ingestion Successful
                                </div>

                                <div className="text-xs font-mono bg-zinc-100 dark:bg-zinc-900 p-4 rounded overflow-auto max-h-96 border border-zinc-200 dark:border-zinc-700">
                                    <pre>{JSON.stringify(ingestResponse, null, 2)}</pre>
                                </div>
                            </div>
                        )}
                    </>
                )}
            </div>

            {/* Input Area */}
            {mode === 'chat' ? (
                <form onSubmit={handleChatSubmit} className="flex flex-col gap-3">
                    <div className="relative">
                        <textarea
                            value={query}
                            onChange={(e) => setQuery(e.target.value)}
                            placeholder="Ask a question..."
                            className="w-full p-4 rounded-xl border border-zinc-300 dark:border-zinc-700 bg-white dark:bg-zinc-800 focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none resize-none shadow-sm transition-all"
                            rows={3}
                            onKeyDown={(e) => {
                                if (e.key === 'Enter' && !e.shiftKey) {
                                    e.preventDefault();
                                    handleChatSubmit(e);
                                }
                            }}
                        />
                        <div className="absolute bottom-3 right-3 text-xs text-zinc-400">Press Enter</div>
                    </div>
                    <button
                        type="submit"
                        disabled={loading || !query.trim()}
                        className="self-end px-6 py-2 rounded-lg bg-blue-600 hover:bg-blue-700 text-white font-medium disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                    >
                        {loading ? 'Thinking...' : 'Generate Response'}
                    </button>
                </form>
            ) : (
                <form onSubmit={handleIngestSubmit} className="flex flex-col gap-3">
                    <div className="relative">
                        <textarea
                            value={ingestText}
                            onChange={(e) => setIngestText(e.target.value)}
                            placeholder="Enter facts or information to start ingestion (e.g. 'Emalie went to Dubai')..."
                            className="w-full p-4 rounded-xl border border-emerald-300 dark:border-emerald-700 bg-white dark:bg-zinc-800 focus:ring-2 focus:ring-emerald-500 focus:border-transparent outline-none resize-none shadow-sm transition-all"
                            rows={5}
                        />
                    </div>
                    <button
                        type="submit"
                        disabled={loading || !ingestText.trim()}
                        className="self-end px-6 py-2 rounded-lg bg-emerald-600 hover:bg-emerald-700 text-white font-medium disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                    >
                        {loading ? 'Processing...' : 'Ingest Data'}
                    </button>
                </form>
            )}

        </div>
    );
}
